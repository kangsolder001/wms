import { useState } from 'react';
import { Table, Button, Modal, Form, Input, InputNumber, Select, Space, Typography, message, Tag } from 'antd';
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { poApi, type PurchaseOrder, type CreatePORequest } from '../../api/purchase-orders';
import { itemApi, type Item } from '../../api/items';
import type { ColumnsType } from 'antd/es/table';

const statusColors: Record<string, string> = {
  pending: 'orange',
  partial: 'blue',
  received: 'green',
  cancelled: 'red',
};

export default function PurchaseOrdersPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['purchase-orders'],
    queryFn: () => poApi.list().then((res) => res.data),
  });

  const { data: itemsData } = useQuery({
    queryKey: ['items'],
    queryFn: () => itemApi.list(1, 100).then((res) => res.data),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreatePORequest) => poApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['purchase-orders'] });
      message.success('Purchase Order created');
      setIsModalOpen(false);
      form.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to create PO'),
  });

  const columns: ColumnsType<PurchaseOrder> = [
    { title: 'PO Number', dataIndex: 'po_number', key: 'po_number' },
    { title: 'Supplier', dataIndex: 'supplier_name', key: 'supplier_name', responsive: ['md'] },
    {
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => <Tag color={statusColors[status]}>{status.toUpperCase()}</Tag>,
    },
    { title: 'Created By', dataIndex: 'created_by_name', key: 'created_by_name', responsive: ['lg'] },
    {
      title: 'Expected Date',
      dataIndex: 'expected_date',
      key: 'expected_date',
      responsive: ['lg'],
      render: (date: string) => date ? new Date(date).toLocaleDateString() : '-',
    },
    {
      title: 'Created At',
      dataIndex: 'created_at',
      key: 'created_at',
      responsive: ['xl'],
      render: (date: string) => new Date(date).toLocaleDateString(),
    },
  ];

  const handleSubmit = (values: any) => {
    const items = values.items?.map((item: any) => ({
      item_id: item.item_id,
      expected_quantity: item.expected_quantity,
      unit_price: item.unit_price || 0,
    })) || [];

    createMutation.mutate({
      supplier_name: values.supplier_name,
      expected_date: values.expected_date,
      notes: values.notes,
      items,
    });
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16, flexWrap: 'wrap', gap: 8 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>Purchase Orders</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setIsModalOpen(true)}>
          Add PO
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
        title="Create Purchase Order"
        open={isModalOpen}
        onCancel={() => { setIsModalOpen(false); form.resetFields(); }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 700 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit}>
          <Form.Item name="supplier_name" label="Supplier Name" rules={[{ required: true }]}>
            <Input placeholder="e.g. PT Supplier Maju" />
          </Form.Item>
          <Form.Item name="expected_date" label="Expected Date">
            <Input type="date" />
          </Form.Item>
          <Form.Item name="notes" label="Notes">
            <Input.TextArea rows={2} />
          </Form.Item>

          <Typography.Title level={5}>Items</Typography.Title>
          <Form.List name="items" rules={[{ validator: async (_, value) => {
            if (!value || value.length === 0) return Promise.reject(new Error('At least one item required'));
          }}]}>
            {(fields, { add, remove }, { errors }) => (
              <>
                {fields.map(({ key, name, ...restField }) => (
                  <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="start">
                    <Form.Item
                      {...restField}
                      name={[name, 'item_id']}
                      rules={[{ required: true, message: 'Select item' }]}
                      style={{ marginBottom: 0, minWidth: 200 }}
                    >
                      <Select placeholder="Select Item" showSearch optionFilterProp="label">
                        {itemsData?.data?.map((item: Item) => (
                          <Select.Option key={item.id} value={item.id} label={`${item.sku} - ${item.name}`}>
                            {item.sku} - {item.name}
                          </Select.Option>
                        ))}
                      </Select>
                    </Form.Item>
                    <Form.Item
                      {...restField}
                      name={[name, 'expected_quantity']}
                      rules={[{ required: true, message: 'Qty' }]}
                      style={{ marginBottom: 0, minWidth: 100 }}
                    >
                      <InputNumber min={1} placeholder="Qty" />
                    </Form.Item>
                    <Form.Item
                      {...restField}
                      name={[name, 'unit_price']}
                      style={{ marginBottom: 0, minWidth: 120 }}
                    >
                      <InputNumber min={0} placeholder="Unit Price" />
                    </Form.Item>
                    <Button icon={<DeleteOutlined />} danger onClick={() => remove(name)} />
                  </Space>
                ))}
                <Form.Item>
                  <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                    Add Item
                  </Button>
                  <Form.ErrorList errors={errors} />
                </Form.Item>
              </>
            )}
          </Form.List>
        </Form>
      </Modal>
    </div>
  );
}
