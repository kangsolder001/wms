import { useState } from 'react';
import { Table, Button, Modal, Form, Input, InputNumber, Select, Space, Typography, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { locationApi } from '../../api/locations';
import type { Location, CreateLocationRequest } from '../../api/locations';
import type { ColumnsType } from 'antd/es/table';

export default function LocationsPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingLocation, setEditingLocation] = useState<Location | null>(null);
  const [form] = Form.useForm();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['locations'],
    queryFn: () => locationApi.list().then((res) => res.data),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateLocationRequest) => locationApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['locations'] });
      message.success('Location created');
      setIsModalOpen(false);
      form.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to create location'),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => locationApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['locations'] });
      message.success('Location deleted');
    },
  });

  const columns: ColumnsType<Location> = [
    { title: 'Code', dataIndex: 'code', key: 'code' },
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Zone', dataIndex: 'zone', key: 'zone', responsive: ['md'] },
    { title: 'Type', dataIndex: 'type', key: 'type', responsive: ['md'] },
    { title: 'Capacity', dataIndex: 'capacity', key: 'capacity', responsive: ['lg'] },
    {
      title: 'Actions',
      key: 'actions',
      fixed: 'right',
      width: 100,
      render: (_: any, record: Location) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => { setEditingLocation(record); setIsModalOpen(true); }} />
          <Button icon={<DeleteOutlined />} danger onClick={() => deleteMutation.mutate(record.id)} />
        </Space>
      ),
    },
  ];

  const handleSubmit = (values: CreateLocationRequest) => {
    createMutation.mutate(values);
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16, flexWrap: 'wrap', gap: 8 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>Locations</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingLocation(null); setIsModalOpen(true); }}>
          Add Location
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
        title={editingLocation ? 'Edit Location' : 'Add Location'}
        open={isModalOpen}
        onCancel={() => { setIsModalOpen(false); setEditingLocation(null); form.resetFields(); }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 520 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit} initialValues={editingLocation || {}}>
          <Form.Item name="code" label="Code" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="zone" label="Zone">
            <Input />
          </Form.Item>
          <Form.Item name="type" label="Type" rules={[{ required: true }]}>
            <Select options={[{ value: 'storage', label: 'Storage' }, { value: 'receiving', label: 'Receiving' }, { value: 'shipping', label: 'Shipping' }, { value: 'staging', label: 'Staging' }]} />
          </Form.Item>
          <Form.Item name="capacity" label="Capacity">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
