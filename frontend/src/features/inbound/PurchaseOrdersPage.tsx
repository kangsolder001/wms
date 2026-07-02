import { useState } from 'react';
import { Table, Button, Modal, Form, Input, InputNumber, Select, Space, Typography, message, Tag, Descriptions } from 'antd';
import { PlusOutlined, DeleteOutlined, InboxOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { poApi, type PurchaseOrder, type CreatePORequest, type ReceiveGoodsRequest } from '../../api/purchase-orders';
import { itemApi, type Item } from '../../api/items';
import api from '../../api/client';
import type { ColumnsType } from 'antd/es/table';

const statusColors: Record<string, string> = {
  pending: 'orange',
  partial: 'blue',
  received: 'green',
  cancelled: 'red',
};

interface Location {
  id: string;
  code: string;
  name: string;
  type: string;
}

export default function PurchaseOrdersPage() {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isReceiveModalOpen, setIsReceiveModalOpen] = useState(false);
  const [selectedPO, setSelectedPO] = useState<PurchaseOrder | null>(null);
  const [createForm] = Form.useForm();
  const [receiveForm] = Form.useForm();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['purchase-orders'],
    queryFn: () => poApi.list().then((res) => res.data),
  });

  const { data: itemsData } = useQuery({
    queryKey: ['items'],
    queryFn: () => itemApi.list(1, 100).then((res) => res.data),
  });

  const { data: locationsData } = useQuery({
    queryKey: ['locations'],
    queryFn: () => api.get('/locations', { params: { page: 1, limit: 100 } }).then((res) => res.data),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreatePORequest) => poApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['purchase-orders'] });
      message.success('Purchase Order created');
      setIsCreateModalOpen(false);
      createForm.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to create PO'),
  });

  const receiveMutation = useMutation({
    mutationFn: ({ id, data }: { id: string; data: ReceiveGoodsRequest }) => poApi.receive(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['purchase-orders'] });
      message.success('Goods received successfully');
      setIsReceiveModalOpen(false);
      setSelectedPO(null);
      receiveForm.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to receive goods'),
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
      title: 'Actions',
      key: 'actions',
      fixed: 'right',
      width: 120,
      render: (_: any, record: PurchaseOrder) => (
        <Space>
          {record.status === 'pending' && (
            <Button
              type="primary"
              size="small"
              icon={<InboxOutlined />}
              onClick={() => {
                setSelectedPO(record);
                setIsReceiveModalOpen(true);
                receiveForm.resetFields();
                // Pre-fill items from PO
                if (record.items && record.items.length > 0) {
                  const receiveItems = record.items.map((item) => ({
                    item_id: item.item_id,
                    quantity: item.expected_quantity,
                    batch_number: '',
                    location_id: undefined,
                  }));
                  receiveForm.setFieldsValue({ items: receiveItems });
                }
              }}
            >
              Receive
            </Button>
          )}
        </Space>
      ),
    },
  ];

  const handleCreateSubmit = (values: any) => {
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

  const handleReceiveSubmit = (values: any) => {
    if (!selectedPO) return;

    const items = values.items?.map((item: any) => ({
      item_id: item.item_id,
      quantity: item.quantity,
      location_id: item.location_id,
    })) || [];

    receiveMutation.mutate({
      id: selectedPO.id,
      data: {
        notes: values.notes,
        items,
      },
    });
  };

  const getItemName = (itemId: string) => {
    const item = itemsData?.data?.find((i: Item) => i.id === itemId);
    return item ? `${item.sku} - ${item.name}` : itemId;
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16, flexWrap: 'wrap', gap: 8 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>Purchase Orders</Typography.Title>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => setIsCreateModalOpen(true)}>
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

      {/* Create PO Modal */}
      <Modal
        title="Create Purchase Order"
        open={isCreateModalOpen}
        onCancel={() => { setIsCreateModalOpen(false); createForm.resetFields(); }}
        onOk={() => createForm.submit()}
        width="90%"
        style={{ maxWidth: 700 }}
      >
        <Form form={createForm} layout="vertical" onFinish={handleCreateSubmit}>
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

      {/* Receive Goods Modal */}
      <Modal
        title={`Receive Goods - ${selectedPO?.po_number || ''}`}
        open={isReceiveModalOpen}
        onCancel={() => { setIsReceiveModalOpen(false); setSelectedPO(null); receiveForm.resetFields(); }}
        onOk={() => receiveForm.submit()}
        width="90%"
        style={{ maxWidth: 800 }}
        confirmLoading={receiveMutation.isPending}
      >
        {selectedPO && (
          <>
            <Descriptions bordered size="small" style={{ marginBottom: 16 }}>
              <Descriptions.Item label="Supplier">{selectedPO.supplier_name}</Descriptions.Item>
              <Descriptions.Item label="Expected Date">
                {selectedPO.expected_date ? new Date(selectedPO.expected_date).toLocaleDateString() : '-'}
              </Descriptions.Item>
              <Descriptions.Item label="Status">
                <Tag color={statusColors[selectedPO.status]}>{selectedPO.status.toUpperCase()}</Tag>
              </Descriptions.Item>
            </Descriptions>

            <Form form={receiveForm} layout="vertical" onFinish={handleReceiveSubmit}>
              <Form.Item name="notes" label="Notes">
                <Input.TextArea rows={2} />
              </Form.Item>

              <Typography.Title level={5}>Receive Items</Typography.Title>
              <Form.List name="items">
                {(fields) => (
                  <>
                    {fields.map(({ key, name }) => (
                      <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="start" wrap>
                        <div style={{ minWidth: 200, fontWeight: 500 }}>
                          {getItemName(receiveForm.getFieldInstance(['items', name, 'item_id'])?.props?.value || selectedPO.items?.[name]?.item_id || '')}
                        </div>
                        <Form.Item
                          name={[name, 'item_id']}
                          hidden
                        >
                          <Input />
                        </Form.Item>
                        <Form.Item
                          name={[name, 'quantity']}
                          rules={[{ required: true, message: 'Qty' }]}
                          style={{ marginBottom: 0, minWidth: 100 }}
                        >
                          <InputNumber min={1} placeholder="Qty" />
                        </Form.Item>
                        <Form.Item
                          name={[name, 'location_id']}
                          rules={[{ required: true, message: 'Location' }]}
                          style={{ marginBottom: 0, minWidth: 200 }}
                        >
                          <Select placeholder="Select Location">
                            {locationsData?.data?.map((loc: Location) => (
                              <Select.Option key={loc.id} value={loc.id}>
                                {loc.code} - {loc.name}
                              </Select.Option>
                            ))}
                          </Select>
                        </Form.Item>
                      </Space>
                    ))}
                  </>
                )}
              </Form.List>
            </Form>
          </>
        )}
      </Modal>
    </div>
  );
}
