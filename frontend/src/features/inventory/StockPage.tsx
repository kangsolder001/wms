import { Table, Typography } from 'antd';
import { useQuery } from '@tanstack/react-query';
import { stockApi, type Stock } from '../../api/stock';
import type { ColumnsType } from 'antd/es/table';

export default function StockPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['stock'],
    queryFn: () => stockApi.list().then((res) => res.data),
  });

  const columns: ColumnsType<Stock> = [
    { title: 'SKU', dataIndex: 'item_sku', key: 'item_sku' },
    { title: 'Item', dataIndex: 'item_name', key: 'item_name' },
    { title: 'Location', dataIndex: 'location_code', key: 'location_code', responsive: ['md'] },
    { title: 'Quantity', dataIndex: 'quantity', key: 'quantity' },
    { title: 'Reserved', dataIndex: 'reserved_quantity', key: 'reserved_quantity', responsive: ['md'] },
    { title: 'Batch', dataIndex: 'batch_number', key: 'batch_number', responsive: ['lg'] },
  ];

  return (
    <div>
      <Typography.Title level={3}>Stock Levels</Typography.Title>
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
