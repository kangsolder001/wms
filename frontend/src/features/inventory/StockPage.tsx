import { Table, Typography } from 'antd';
import { useQuery } from '@tanstack/react-query';
import { stockApi } from '../../api/stock';

export default function StockPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['stock'],
    queryFn: () => stockApi.list().then((res) => res.data),
  });

  const columns = [
    { title: 'Item ID', dataIndex: 'item_id', key: 'item_id' },
    { title: 'Location ID', dataIndex: 'location_id', key: 'location_id', responsive: ['md' as const] },
    { title: 'Quantity', dataIndex: 'quantity', key: 'quantity' },
    { title: 'Reserved', dataIndex: 'reserved_quantity', key: 'reserved_quantity', responsive: ['md' as const] },
    { title: 'Batch', dataIndex: 'batch_number', key: 'batch_number', responsive: ['lg' as const] },
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
