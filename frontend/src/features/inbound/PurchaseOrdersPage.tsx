import { Table, Typography } from 'antd';
import { useQuery } from '@tanstack/react-query';
import api from '../../api/client';

export default function PurchaseOrdersPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['purchase-orders'],
    queryFn: () => api.get('/purchase-orders').then((res) => res.data),
  });

  const columns = [
    { title: 'PO Number', dataIndex: 'po_number', key: 'po_number' },
    { title: 'Supplier', dataIndex: 'supplier_name', key: 'supplier_name', responsive: ['md' as const] },
    { title: 'Status', dataIndex: 'status', key: 'status' },
    { title: 'Created At', dataIndex: 'created_at', key: 'created_at', responsive: ['lg' as const] },
  ];

  return (
    <div>
      <Typography.Title level={3}>Purchase Orders</Typography.Title>
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
