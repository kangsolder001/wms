import { useState } from 'react';
import { Table, Button, Modal, Form, Input, Select, Tag, Space, Typography, message, Popconfirm, Switch } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { authApi, type UserResponse } from '../../api/auth';
import type { ColumnsType } from 'antd/es/table';

const roleColors: Record<string, string> = {
  superadmin: 'red',
  manager: 'blue',
  operator: 'green',
  viewer: 'default',
};

export default function UsersPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<UserResponse | null>(null);
  const [form] = Form.useForm();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['users'],
    queryFn: () => authApi.listUsers().then((res) => res.data),
  });

  const createMutation = useMutation({
    mutationFn: (data: any) => authApi.register(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      message.success('User created');
      setIsModalOpen(false);
      form.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to create user'),
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => authApi.updateUser(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      message.success('User updated');
      setIsModalOpen(false);
      setEditingUser(null);
      form.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to update user'),
  });

  const deactivateMutation = useMutation({
    mutationFn: (id: string) => authApi.deleteUser(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      message.success('User deactivated');
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to deactivate user'),
  });

  const toggleActiveMutation = useMutation({
    mutationFn: ({ id, is_active }: { id: string; is_active: boolean }) => authApi.updateUser(id, { is_active }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['users'] });
      message.success('User status updated');
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to update status'),
  });

  const columns: ColumnsType<UserResponse> = [
    { title: 'Username', dataIndex: 'username', key: 'username' },
    { title: 'Full Name', dataIndex: 'full_name', key: 'full_name', responsive: ['md'] },
    { title: 'Email', dataIndex: 'email', key: 'email', responsive: ['lg'] },
    {
      title: 'Role',
      dataIndex: 'role',
      key: 'role',
      render: (role: string) => <Tag color={roleColors[role]}>{role.toUpperCase()}</Tag>,
    },
    {
      title: 'Status',
      dataIndex: 'is_active',
      key: 'is_active',
      render: (isActive: boolean, record) => (
        <Popconfirm
          title={isActive ? 'Deactivate this user?' : 'Activate this user?'}
          onConfirm={() => toggleActiveMutation.mutate({ id: record.id, is_active: !isActive })}
          okText="Yes"
          cancelText="No"
        >
          <Switch checked={isActive} size="small" />
        </Popconfirm>
      ),
    },
    {
      title: 'Actions',
      key: 'actions',
      fixed: 'right',
      width: 100,
      render: (_: any, record: UserResponse) => (
        <Space>
          <Button icon={<EditOutlined />} onClick={() => { setEditingUser(record); setIsModalOpen(true); }} />
          <Popconfirm title="Deactivate this user?" onConfirm={() => deactivateMutation.mutate(record.id)}>
            <Button icon={<DeleteOutlined />} danger />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const handleSubmit = (values: any) => {
    if (editingUser) {
      const updateData: any = {
        email: values.email,
        full_name: values.full_name,
        role: values.role,
      };
      if (values.password) updateData.password = values.password;
      updateMutation.mutate({ id: editingUser.id, data: updateData });
    } else {
      createMutation.mutate(values);
    }
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16, flexWrap: 'wrap', gap: 8 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>User Management</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingUser(null); setIsModalOpen(true); }}>
          Add User
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
        title={editingUser ? 'Edit User' : 'Add User'}
        open={isModalOpen}
        onCancel={() => { setIsModalOpen(false); setEditingUser(null); form.resetFields(); }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 520 }}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
          initialValues={editingUser ? { ...editingUser } : { role: 'operator' }}
        >
          {!editingUser && (
            <Form.Item name="username" label="Username" rules={[{ required: true, min: 3 }]}>
              <Input placeholder="e.g. operator1" />
            </Form.Item>
          )}
          <Form.Item name="full_name" label="Full Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="email" label="Email" rules={[{ required: !editingUser, type: 'email' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="role" label="Role" rules={[{ required: true }]}>
            <Select>
              <Select.Option value="operator">Operator</Select.Option>
              <Select.Option value="manager">Manager</Select.Option>
              <Select.Option value="superadmin">Super Admin</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item
            name="password"
            label={editingUser ? 'New Password (leave empty to keep)' : 'Password'}
            rules={editingUser ? [] : [{ required: true, min: 6 }]}
          >
            <Input.Password />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
}
