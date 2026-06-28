import { Card, Row, Col, Statistic, Typography } from 'antd';
import { useQuery } from '@tanstack/react-query';
import { DatabaseOutlined, ShoppingCartOutlined, EnvironmentOutlined } from '@ant-design/icons';
import api from '../../api/client';

export default function DashboardPage() {
  const { data: summary, isLoading } = useQuery({
    queryKey: ['dashboard-summary'],
    queryFn: () => api.get('/dashboard/summary').then((res) => res.data),
  });

  return (
    <div>
      <Typography.Title level={3}>Dashboard</Typography.Title>
      <Row gutter={[16, 16]}>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="Total Stock Items"
              value={summary?.data?.total_stock_items || 0}
              prefix={<DatabaseOutlined />}
              loading={isLoading}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="Total Quantity"
              value={summary?.data?.total_quantity || 0}
              prefix={<ShoppingCartOutlined />}
              loading={isLoading}
            />
          </Card>
        </Col>
        <Col xs={24} sm={8}>
          <Card>
            <Statistic
              title="Active Locations"
              value={0}
              prefix={<EnvironmentOutlined />}
              loading={isLoading}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
}
