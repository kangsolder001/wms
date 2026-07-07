import { useState } from 'react';
import { Table, Typography, Input, Select, Space, Button, Modal, Form, InputNumber, message } from 'antd';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { stockApi, type Stock } from '../../api/stock';
import { itemApi, type Item } from '../../api/items';
import api from '../../api/client';
import type { ColumnsType } from 'antd/es/table';

interface Location {
  id: string;
  code: string;
  name: string;
  type: string;
}

export default function StockPage() {
  const [search, setSearch] = useState('');
  const [itemFilter, setItemFilter] = useState<string | undefined>(undefined);
  const [locationFilter, setLocationFilter] = useState<string | undefined>(undefined);
  const [isAdjustModalOpen, setIsAdjustModalOpen] = useState(false);
  const [isOpnameModalOpen, setIsOpnameModalOpen] = useState(false);
  const [selectedStock, setSelectedStock] = useState<Stock | null>(null);
  const [adjustForm] = Form.useForm();
  const [opnameForm] = Form.useForm();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['stock', search, itemFilter, locationFilter],
    queryFn: () => stockApi.list(1, 100, itemFilter, locationFilter, search).then((res) => res.data),
  });

  const { data: itemsData } = useQuery({
    queryKey: ['items'],
    queryFn: () => itemApi.list(1, 100).then((res) => res.data),
  });

  const { data: locationsData } = useQuery({
    queryKey: ['locations'],
    queryFn: () => api.get('/locations', { params: { page: 1, limit: 100 } }).then((res) => res.data),
  });

  const adjustMutation = useMutation({
    mutationFn: stockApi.adjust,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stock'] });
      message.success('Stock adjusted');
      setIsAdjustModalOpen(false);
      setSelectedStock(null);
      adjustForm.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to adjust'),
  });

  const opnameMutation = useMutation({
    mutationFn: stockApi.opname,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['stock'] });
      message.success('Stock opname completed');
      setIsOpnameModalOpen(false);
      opnameForm.resetFields();
    },
    onError: (error: any) => message.error(error.response?.data?.error || 'Failed to submit opname'),
  });

  const columns: ColumnsType<Stock> = [
    { title: 'SKU', dataIndex: 'item_sku', key: 'item_sku' },
    { title: 'Item', dataIndex: 'item_name', key: 'item_name' },
    { title: 'Location', dataIndex: 'location_code', key: 'location_code' },
    { title: 'Quantity', dataIndex: 'quantity', key: 'quantity' },
    { title: 'Reserved', dataIndex: 'reserved_quantity', key: 'reserved_quantity', responsive: ['md'] },
    { title: 'Batch', dataIndex: 'batch_number', key: 'batch_number', responsive: ['lg'] },
    {
      title: 'Actions',
      key: 'actions',
      fixed: 'right',
      width: 100,
      render: (_: any, record: Stock) => (
        <Button size="small" onClick={() => {
          setSelectedStock(record);
          adjustForm.setFieldsValue({
            item_id: record.item_id,
            location_id: record.location_id,
            quantity: record.quantity,
          });
          setIsAdjustModalOpen(true);
        }}>
          Adjust
        </Button>
      ),
    },
  ];

  const handleAdjustSubmit = (values: any) => {
    adjustMutation.mutate(values);
  };

  const handleOpnameSubmit = (values: any) => {
    const locationId = values.location_id;
    const stockItems = data?.data?.filter((s: Stock) => s.location_id === locationId) || [];

    const items = stockItems.map((s: Stock) => ({
      item_id: s.item_id,
      system_quantity: s.quantity,
      actual_quantity: values[`qty_${s.item_id}`] ?? s.quantity,
    }));

    opnameMutation.mutate({
      location_id: locationId,
      notes: values.notes,
      items,
    });
  };

  const stockAtLocation = data?.data?.filter((s: Stock) => {
    const opnameLocation = opnameForm.getFieldValue('location_id');
    return !opnameLocation || s.location_id === opnameLocation;
  }) || [];

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16, flexWrap: 'wrap', gap: 8 }}>
        <Typography.Title level={3} style={{ margin: 0 }}>Stock Levels</Typography.Title>
        <Space>
          <Button onClick={() => { setIsOpnameModalOpen(true); opnameForm.resetFields(); }}>
            Stock Opname
          </Button>
        </Space>
      </div>

      <Space style={{ marginBottom: 16 }} wrap>
        <Input.Search
          placeholder="Search SKU / Item Name / Location"
          allowClear
          onSearch={(value) => setSearch(value)}
          style={{ width: 300 }}
        />
        <Select
          placeholder="Filter by Item"
          allowClear
          showSearch
          optionFilterProp="label"
          style={{ width: 250 }}
          onChange={(value) => setItemFilter(value)}
        >
          {itemsData?.data?.map((item: Item) => (
            <Select.Option key={item.id} value={item.id} label={`${item.sku} - ${item.name}`}>
              {item.sku} - {item.name}
            </Select.Option>
          ))}
        </Select>
        <Select
          placeholder="Filter by Location"
          allowClear
          showSearch
          optionFilterProp="label"
          style={{ width: 250 }}
          onChange={(value) => setLocationFilter(value)}
        >
          {locationsData?.data?.map((loc: Location) => (
            <Select.Option key={loc.id} value={loc.id} label={`${loc.code} - ${loc.name}`}>
              {loc.code} - {loc.name}
            </Select.Option>
          ))}
        </Select>
      </Space>

      <Table
        columns={columns}
        dataSource={data?.data || []}
        rowKey="id"
        loading={isLoading}
        pagination={{ pageSize: 10 }}
        scroll={{ x: 'max-content' }}
      />

      {/* Adjust Modal */}
      <Modal
        title="Adjust Stock"
        open={isAdjustModalOpen}
        onCancel={() => { setIsAdjustModalOpen(false); setSelectedStock(null); adjustForm.resetFields(); }}
        onOk={() => adjustForm.submit()}
        confirmLoading={adjustMutation.isPending}
      >
        <Form form={adjustForm} layout="vertical" onFinish={handleAdjustSubmit}>
          <Form.Item name="item_id" hidden><Input /></Form.Item>
          <Form.Item name="location_id" hidden><Input /></Form.Item>
          <Form.Item label="Current Quantity">
            <InputNumber value={selectedStock?.quantity} disabled style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="quantity" label="New Quantity" rules={[{ required: true }]}>
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="notes" label="Notes">
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>

      {/* Stock Opname Modal */}
      <Modal
        title="Stock Opname"
        open={isOpnameModalOpen}
        onCancel={() => { setIsOpnameModalOpen(false); opnameForm.resetFields(); }}
        onOk={() => opnameForm.submit()}
        width="90%"
        style={{ maxWidth: 800 }}
        confirmLoading={opnameMutation.isPending}
      >
        <Form form={opnameForm} layout="vertical" onFinish={handleOpnameSubmit}>
          <Form.Item name="location_id" label="Location" rules={[{ required: true }]}>
            <Select placeholder="Select location to count" showSearch optionFilterProp="label"
              onChange={() => opnameForm.resetFields(['items'])}
            >
              {locationsData?.data?.map((loc: Location) => (
                <Select.Option key={loc.id} value={loc.id} label={`${loc.code} - ${loc.name}`}>
                  {loc.code} - {loc.name}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>
          <Form.Item name="notes" label="Notes">
            <Input.TextArea rows={2} placeholder="Opname notes..." />
          </Form.Item>

          <Typography.Title level={5}>Count Items</Typography.Title>
          {opnameForm.getFieldValue('location_id') && stockAtLocation.length > 0 ? (
            stockAtLocation.map((stock: Stock) => (
              <Space key={stock.id} style={{ display: 'flex', marginBottom: 8 }} align="start">
                <div style={{ minWidth: 250 }}>
                  <strong>{stock.item_sku}</strong> - {stock.item_name}
                </div>
                <div style={{ minWidth: 100, color: '#999' }}>
                  System: {stock.quantity}
                </div>
                <Form.Item
                  name={`qty_${stock.item_id}`}
                  style={{ marginBottom: 0, minWidth: 120 }}
                  initialValue={stock.quantity}
                >
                  <InputNumber min={0} placeholder="Actual qty" />
                </Form.Item>
              </Space>
            ))
          ) : opnameForm.getFieldValue('location_id') ? (
            <Typography.Text type="secondary">No stock at this location</Typography.Text>
          ) : (
            <Typography.Text type="secondary">Select a location first</Typography.Text>
          )}
        </Form>
      </Modal>
    </div>
  );
}
