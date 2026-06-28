import { Table, Typography } from 'antd';
import { useQuery } from '@tanstack/react-query';
import api from '../../api/client';

export default function SalesOrdersPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['sales-orders'],
    queryFn: () => api.get('/sales-orders').then((res) => res.data),
  });

  const columns = [
    { title: 'SO Number', dataIndex: 'so_number', key: 'so_number' },
    { title: 'Customer', dataIndex: 'customer_name', key: 'customer_name', responsive: ['md' as const] },
    { title: 'Status', dataIndex: 'status', key: 'status' },
    { title: 'Created At', dataIndex: 'created_at', key: 'created_at', responsive: ['lg' as const] },
  ];

  return (
    <div>
      <Typography.Title level={3}>Sales Orders</Typography.Title>
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
