import { useState } from 'react';
import { Table, Button, Modal, Form, Input, Space, Typography, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { categoryApi } from '../../api/categories';
import type { Category, CreateCategoryRequest } from '../../api/categories';
import type { ColumnsType } from 'antd/es/table';

export default function CategoriesPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<Category | null>(null);
  const [form] = Form.useForm();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['categories'],
    queryFn: () => categoryApi.list().then((res) => res.data),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateCategoryRequest) => categoryApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      message.success('Category created');
      setIsModalOpen(false);
      form.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to create category'),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => categoryApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['categories'] });
      message.success('Category deleted');
    },
  });

  const columns: ColumnsType<Category> = [
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Abbreviation', dataIndex: 'abbreviation', key: 'abbreviation' },
    { title: 'Description', dataIndex: 'description', key: 'description', responsive: ['md'] },
    {
      title: 'Actions',
      key: 'actions',
      fixed: 'right',
      width: 100,
      render: (_: unknown, record: Category) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => { setEditingCategory(record); setIsModalOpen(true); }} />
          <Button icon={<DeleteOutlined />} danger onClick={() => deleteMutation.mutate(record.id)} />
        </Space>
      ),
    },
  ];

  const handleSubmit = (values: CreateCategoryRequest) => {
    createMutation.mutate(values);
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16, flexWrap: 'wrap', gap: 8 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>Categories</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingCategory(null); setIsModalOpen(true); }}>
          Add Category
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
        title={editingCategory ? 'Edit Category' : 'Add Category'}
        open={isModalOpen}
        onCancel={() => { setIsModalOpen(false); setEditingCategory(null); form.resetFields(); }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 520 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit} initialValues={editingCategory || {}}>
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="abbreviation" label="Abbreviation" rules={[{ required: true }, { len: 3, message: 'Must be exactly 3 characters' }]}>
            <Input maxLength={3} />
          </Form.Item>
          <Form.Item name="description" label="Description">
            <Input.TextArea />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
