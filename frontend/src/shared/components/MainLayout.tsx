import { useState, useEffect } from 'react';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { Layout, Menu, Button, Typography, Drawer, Modal, Space, Switch } from 'antd';
import {
  DashboardOutlined,
  ShoppingCartOutlined,
  EnvironmentOutlined,
  DatabaseOutlined,
  InboxOutlined,
  ShopOutlined,
  SwapOutlined,
  LogoutOutlined,
  MenuOutlined,
  SunOutlined,
  UserAddOutlined,
  AppstoreOutlined,
} from '@ant-design/icons';
import { useAuthStore } from '../../stores/authStore';
import { useThemeStore } from '../../stores/themeStore';

const { Header, Sider, Content } = Layout;

const baseMenuItems = [
  { key: '/', icon: <DashboardOutlined />, label: 'Dashboard' },
  { key: '/stock', icon: <DatabaseOutlined />, label: 'Stock' },
  { key: '/purchase-orders', icon: <InboxOutlined />, label: 'Purchase Orders' },
  { key: '/sales-orders', icon: <ShopOutlined />, label: 'Sales Orders' },
  { key: '/transfers', icon: <SwapOutlined />, label: 'Transfers' },
  {
    key: 'master-data',
    icon: <AppstoreOutlined />,
    label: 'Master Data',
    children: [
      { key: '/items', icon: <ShoppingCartOutlined />, label: 'Items' },
      { key: '/categories', icon: <AppstoreOutlined />, label: 'Categories' },
      { key: '/locations', icon: <EnvironmentOutlined />, label: 'Locations' },
      { key: '/zones', icon: <EnvironmentOutlined />, label: 'Zones' },
    ],
  },
];

const MOBILE_BREAKPOINT = 768;

export default function MainLayout() {
  const [collapsed, setCollapsed] = useState(false);
  const [drawerOpen, setDrawerOpen] = useState(false);
  const [isMobile, setIsMobile] = useState(window.innerWidth < MOBILE_BREAKPOINT);
  const navigate = useNavigate();
  const location = useLocation();
  const logout = useAuthStore((s) => s.logout);
  const user = useAuthStore((s) => s.user);
  const { isDark, toggle: toggleTheme } = useThemeStore();

  const menuItems = [
    ...baseMenuItems,
    ...((user?.role === 'superadmin' || user?.role === 'manager') ? [{ key: '/users', icon: <UserAddOutlined />, label: 'User Management' }] : []),
  ];

  const handleLogout = () => {
    Modal.confirm({
      title: 'Logout',
      content: 'Are you sure you want to logout?',
      okText: 'Yes',
      cancelText: 'No',
      onOk: () => logout(),
    });
  };

  useEffect(() => {
    const handleResize = () => setIsMobile(window.innerWidth < MOBILE_BREAKPOINT);
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  useEffect(() => {
    setDrawerOpen(false);
  }, [location.pathname]);

  const sidebarContent = (
    <>
      <div style={{ margin: 16, textAlign: 'center' }}>
        <img src="/logo.png" alt="WMS Logo" style={{ width: collapsed ? 40 : 80, height: 'auto' }} />
        {!collapsed && <div style={{ color: isDark ? '#c9d1d9' : '#24292f', fontSize: 14, marginTop: 8 }}>WMS System</div>}
      </div>
      <Menu
        theme="dark"
        mode="inline"
        selectedKeys={[location.pathname]}
        items={menuItems}
        onClick={({ key }) => navigate(key)}
      />
    </>
  );

  return (
    <Layout style={{ minHeight: '100vh' }}>
      {isMobile ? (
        <Drawer
          placement="left"
          open={drawerOpen}
          onClose={() => setDrawerOpen(false)}
          width={200}
          styles={{ body: { padding: 0 }, header: { display: 'none' } }}
        >
          {sidebarContent}
        </Drawer>
      ) : (
        <Sider collapsible collapsed={collapsed} onCollapse={setCollapsed}>
          {sidebarContent}
        </Sider>
      )}
      <Layout>
        <Header style={{ padding: '0 16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center', gap: 8 }}>
          {isMobile && (
            <Button
              type="text"
              icon={<MenuOutlined />}
              onClick={() => setDrawerOpen(true)}
              style={{ fontSize: 18 }}
            />
          )}
          <Typography.Title level={5} style={{ margin: 0, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', flex: 1 }}>
            Warehouse Management System
          </Typography.Title>
          <Space size={12}>
            <SunOutlined style={{ fontSize: 14 }} />
            <Switch
              checked={isDark}
              onChange={toggleTheme}
              checkedChildren="Dark"
              unCheckedChildren="Light"
            />
            <Button icon={<LogoutOutlined />} onClick={handleLogout} size={isMobile ? 'small' : 'middle'}>
              {!isMobile && 'Logout'}
            </Button>
          </Space>
        </Header>
        <Content style={{ margin: isMobile ? 8 : 24, padding: isMobile ? 12 : 24, borderRadius: 8, overflow: 'auto' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
