import { useState } from 'react';
import { Table, Button, Modal, Form, Input, Space, Typography, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { zoneApi } from '../../api/zones';
import type { Zone, CreateZoneRequest } from '../../api/zones';
import type { ColumnsType } from 'antd/es/table';

export default function ZonesPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingZone, setEditingZone] = useState<Zone | null>(null);
  const [form] = Form.useForm();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['zones'],
    queryFn: () => zoneApi.list().then((res) => res.data),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateZoneRequest) => zoneApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['zones'] });
      message.success('Zone created');
      setIsModalOpen(false);
      form.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to create zone'),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => zoneApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['zones'] });
      message.success('Zone deleted');
    },
  });

  const columns: ColumnsType<Zone> = [
    { title: 'Code', dataIndex: 'code', key: 'code' },
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Description', dataIndex: 'description', key: 'description', responsive: ['md'] },
    {
      title: 'Actions',
      key: 'actions',
      fixed: 'right',
      width: 100,
      render: (_: unknown, record: Zone) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => { setEditingZone(record); setIsModalOpen(true); }} />
          <Button icon={<DeleteOutlined />} danger onClick={() => deleteMutation.mutate(record.id)} />
        </Space>
      ),
    },
  ];

  const handleSubmit = (values: CreateZoneRequest) => {
    createMutation.mutate(values);
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16, flexWrap: 'wrap', gap: 8 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>Zones</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingZone(null); setIsModalOpen(true); }}>
          Add Zone
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={data?.data || []}
        rowKey="id"
        loading={isLoading}
        pagination={{ pageSize: 10 }}
        scroll={{ x: 'max-content' }}
      />

      <Modal
        title={editingZone ? 'Edit Zone' : 'Add Zone'}
        open={isModalOpen}
        onCancel={() => { setIsModalOpen(false); setEditingZone(null); form.resetFields(); }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 520 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit} initialValues={editingZone || {}}>
          <Form.Item name="code" label="Code" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="Description">
            <Input.TextArea />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
