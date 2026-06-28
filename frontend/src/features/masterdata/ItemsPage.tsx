import { useState } from 'react';
import { Table, Button, Modal, Form, Input, InputNumber, Select, Space, Typography, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { itemApi } from '../../api/items';
import type { Item, CreateItemRequest } from '../../api/items';
import type { ColumnsType } from 'antd/es/table';

export default function ItemsPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<Item | null>(null);
  const [form] = Form.useForm();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['items'],
    queryFn: () => itemApi.list().then((res) => res.data),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateItemRequest) => itemApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] });
      message.success('Item created');
      setIsModalOpen(false);
      form.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to create item'),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => itemApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] });
      message.success('Item deleted');
    },
  });

  const columns: ColumnsType<Item> = [
    { title: 'SKU', dataIndex: 'sku', key: 'sku' },
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Category', dataIndex: 'category', key: 'category', responsive: ['md'] },
    { title: 'UoM', dataIndex: 'unit_of_measure', key: 'unit_of_measure', responsive: ['md'] },
    { title: 'Weight', dataIndex: 'weight', key: 'weight', responsive: ['lg'] },
    {
      title: 'Actions',
      key: 'actions',
      fixed: 'right',
      width: 100,
      render: (_: any, record: Item) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => { setEditingItem(record); setIsModalOpen(true); }} />
          <Button icon={<DeleteOutlined />} danger onClick={() => deleteMutation.mutate(record.id)} />
        </Space>
      ),
    },
  ];

  const handleSubmit = (values: CreateItemRequest) => {
    createMutation.mutate(values);
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16, flexWrap: 'wrap', gap: 8 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>Items</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingItem(null); setIsModalOpen(true); }}>
          Add Item
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
        title={editingItem ? 'Edit Item' : 'Add Item'}
        open={isModalOpen}
        onCancel={() => { setIsModalOpen(false); setEditingItem(null); form.resetFields(); }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 520 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit} initialValues={editingItem || {}}>
          <Form.Item name="sku" label="SKU" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="Description">
            <Input.TextArea />
          </Form.Item>
          <Form.Item name="category" label="Category">
            <Input />
          </Form.Item>
          <Form.Item name="unit_of_measure" label="Unit of Measure" rules={[{ required: true }]}>
            <Select options={[{ value: 'pcs', label: 'Pieces' }, { value: 'kg', label: 'Kilogram' }, { value: 'box', label: 'Box' }]} />
          </Form.Item>
          <Form.Item name="weight" label="Weight">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
