import { useState } from 'react';
import { Table, Button, Modal, Form, Input, InputNumber, Select, Space, Row, Col, Typography, message } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, PrinterOutlined } from '@ant-design/icons';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { itemApi } from '../../api/items';
import { categoryApi } from '../../api/categories';
import type { Item, CreateItemRequest } from '../../api/items';
import type { ColumnsType } from 'antd/es/table';
import QRLabelModal from './components/QRLabelModal';

export default function ItemsPage() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<Item | null>(null);
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([]);
  const [printModalItems, setPrintModalItems] = useState<Item[]>([]);
  const [isPrintModalOpen, setIsPrintModalOpen] = useState(false);
  const [generatedSKU, setGeneratedSKU] = useState<string>('');
  const [form] = Form.useForm();
  const queryClient = useQueryClient();

  const { data, isLoading } = useQuery({
    queryKey: ['items'],
    queryFn: () => itemApi.list().then((res) => res.data),
  });

  const { data: categoriesData } = useQuery({
    queryKey: ['categories-all'],
    queryFn: () => categoryApi.listAll().then((res) => res.data.data),
  });

  const createMutation = useMutation({
    mutationFn: (data: CreateItemRequest) => itemApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] });
      message.success('Item created');
      setIsModalOpen(false);
      form.resetFields();
      setGeneratedSKU('');
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

  const handleCategoryChange = async (categoryId: string) => {
    if (categoryId) {
      try {
        const res = await itemApi.generateSKU(categoryId);
        setGeneratedSKU(res.data.sku);
        form.setFieldValue('sku', res.data.sku);
      } catch {
        message.error('Failed to generate SKU');
      }
    } else {
      setGeneratedSKU('');
      form.setFieldValue('sku', '');
    }
  };

  const handlePrintSingle = (item: Item) => {
    setPrintModalItems([item]);
    setIsPrintModalOpen(true);
  };

  const handlePrintSelected = () => {
    const selectedItems = data?.data?.filter((item: Item) =>
      selectedRowKeys.includes(item.id)
    ) || [];
    setPrintModalItems(selectedItems);
    setIsPrintModalOpen(true);
    setSelectedRowKeys([]);
  };

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
      width: 150,
      render: (_: unknown, record: Item) => (
        <Space>
          <Button icon={<PrinterOutlined />} onClick={() => handlePrintSingle(record)} />
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
        <Space>
          {selectedRowKeys.length > 0 && (
            <Button icon={<PrinterOutlined />} onClick={handlePrintSelected}>
              Print Selected ({selectedRowKeys.length})
            </Button>
          )}
          <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditingItem(null); setGeneratedSKU(''); setIsModalOpen(true); }}>
            Add Item
          </Button>
        </Space>
      </div>

      <Table
        columns={columns}
        dataSource={data?.data || []}
        rowKey="id"
        loading={isLoading}
        pagination={{ pageSize: 10 }}
        scroll={{ x: 'max-content' }}
        rowSelection={{
          selectedRowKeys,
          onChange: setSelectedRowKeys,
        }}
      />

      <Modal
        title={editingItem ? 'Edit Item' : 'Add Item'}
        open={isModalOpen}
        onCancel={() => { setIsModalOpen(false); setEditingItem(null); form.resetFields(); setGeneratedSKU(''); }}
        onOk={() => form.submit()}
        width="90%"
        style={{ maxWidth: 520 }}
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit} initialValues={editingItem || {}}>
          <Form.Item name="category_id" label="Category" rules={[{ required: true }]}>
            <Select
              placeholder="Select category"
              onChange={handleCategoryChange}
              options={(categoriesData || []).map((cat: any) => ({ value: cat.id, label: `${cat.name} (${cat.abbreviation})` }))}
            />
          </Form.Item>
          <Form.Item name="sku" label="SKU">
            <Input disabled placeholder="Auto-generated" />
          </Form.Item>
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="Description">
            <Input.TextArea />
          </Form.Item>
          <Form.Item name="unit_of_measure" label="Unit of Measure" rules={[{ required: true }]}>
            <Select options={[{ value: 'pcs', label: 'Pieces' }, { value: 'kg', label: 'Kilogram' }, { value: 'box', label: 'Box' }]} />
          </Form.Item>
          <Form.Item name="weight" label="Weight">
            <InputNumber min={0} style={{ width: '100%' }} />
          </Form.Item>
          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="length" label="Length (cm)">
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="width" label="Width (cm)">
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="height" label="Height (cm)">
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </Modal>

      <QRLabelModal
        items={printModalItems}
        isOpen={isPrintModalOpen}
        onClose={() => setIsPrintModalOpen(false)}
      />
    </div>
  );
}
