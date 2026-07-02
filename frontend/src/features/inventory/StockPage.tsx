import { useState } from 'react';
import { Table, Typography, Input, Select, Space } from 'antd';
import { useQuery } from '@tanstack/react-query';
import { stockApi, type Stock } from '../../api/stock';
import { itemApi, type Item } from '../../api/items';
import api from '../../api/client';
import type { ColumnsType } from 'antd/es/table';

interface Location {
  id: string;
  code: string;
  name: string;
}

export default function StockPage() {
  const [search, setSearch] = useState('');
  const [itemFilter, setItemFilter] = useState<string | undefined>(undefined);
  const [locationFilter, setLocationFilter] = useState<string | undefined>(undefined);

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

  const columns: ColumnsType<Stock> = [
    { title: 'SKU', dataIndex: 'item_sku', key: 'item_sku' },
    { title: 'Item', dataIndex: 'item_name', key: 'item_name' },
    { title: 'Location', dataIndex: 'location_code', key: 'location_code' },
    { title: 'Quantity', dataIndex: 'quantity', key: 'quantity' },
    { title: 'Reserved', dataIndex: 'reserved_quantity', key: 'reserved_quantity', responsive: ['md'] },
    { title: 'Batch', dataIndex: 'batch_number', key: 'batch_number', responsive: ['lg'] },
  ];

  return (
    <div>
      <Typography.Title level={3}>Stock Levels</Typography.Title>

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
    </div>
  );
}
