import { Table, Typography } from 'antd';
import { useQuery } from '@tanstack/react-query';
import api from '../../api/client';

export default function TransfersPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['transfers'],
    queryFn: () => api.get('/transfers').then((res) => res.data),
  });

  const columns = [
    { title: 'Transfer Number', dataIndex: 'transfer_number', key: 'transfer_number' },
    { title: 'From', dataIndex: 'from_location_id', key: 'from_location_id', responsive: ['md' as const] },
    { title: 'To', dataIndex: 'to_location_id', key: 'to_location_id', responsive: ['md' as const] },
    { title: 'Item ID', dataIndex: 'item_id', key: 'item_id', responsive: ['lg' as const] },
    { title: 'Quantity', dataIndex: 'quantity', key: 'quantity' },
    { title: 'Status', dataIndex: 'status', key: 'status' },
  ];

  return (
    <div>
      <Typography.Title level={3}>Stock Transfers</Typography.Title>
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
