import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ConfigProvider, theme } from 'antd';
import { useAuthStore } from './stores/authStore';
import { useThemeStore } from './stores/themeStore';
import LoginPage from './features/auth/LoginPage';
import RegisterPage from './features/auth/RegisterPage';
import DashboardPage from './features/dashboard/DashboardPage';
import ItemsPage from './features/masterdata/ItemsPage';
import LocationsPage from './features/location/LocationsPage';
import StockPage from './features/inventory/StockPage';
import PurchaseOrdersPage from './features/inbound/PurchaseOrdersPage';
import SalesOrdersPage from './features/outbound/SalesOrdersPage';
import TransfersPage from './features/transfer/TransfersPage';
import MainLayout from './shared/components/MainLayout';

const queryClient = new QueryClient();

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const token = localStorage.getItem('token');
  if (!isAuthenticated && !token) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}

const darkThemeTokens = {
  colorBgBase: '#1E1E2E',
  colorBgContainer: '#2B2B3B',
  colorBgElevated: '#313145',
  colorBgLayout: '#181825',
  colorBgSpotlight: '#313145',
  colorText: '#CDD6F4',
  colorTextSecondary: '#A6ADC8',
  colorTextTertiary: '#7F849C',
  colorTextQuaternary: '#585B70',
  colorBorder: '#45475A',
  colorBorderSecondary: '#313145',
  colorFill: '#313145',
  colorFillSecondary: '#2B2B3B',
  colorFillTertiary: '#242434',
  colorFillQuaternary: '#1E1E2E',
  controlItemBgActive: '#313145',
  controlItemBgHover: '#2B2B3B',
  colorPrimary: '#89B4FA',
  colorPrimaryBg: '#1E1E2E',
  colorPrimaryBgHover: '#2B2B3B',
  colorPrimaryBorder: '#89B4FA',
  colorPrimaryBorderHover: '#B4BEFE',
  colorPrimaryHover: '#B4BEFE',
  colorPrimaryActive: '#CBA6F7',
  colorPrimaryTextHover: '#B4BEFE',
  colorPrimaryText: '#89B4FA',
  colorPrimaryTextActive: '#CBA6F7',
  colorErrorBg: '#2D1B1E',
  colorErrorBorder: '#F38BA8',
  colorSuccessBg: '#1B2E1D',
  colorSuccessBorder: '#A6E3A1',
  colorWarningBg: '#2E2A1D',
  colorWarningBorder: '#F9E2AF',
};

function App() {
  const isDark = useThemeStore((s) => s.isDark);

  return (
    <ConfigProvider theme={{
      token: { colorPrimary: '#1677ff' },
      algorithm: isDark ? theme.darkAlgorithm : theme.defaultAlgorithm,
      ...(isDark && {
        token: darkThemeTokens,
        components: {
          Layout: {
            headerBg: '#181825',
            siderBg: '#181825',
            bodyBg: '#1E1E2E',
            headerColor: '#CDD6F4',
            triggerBg: '#181825',
            triggerColor: '#CDD6F4',
          },
          Menu: {
            darkItemBg: '#181825',
            darkSubMenuItemBg: '#181825',
            darkItemColor: '#A6ADC8',
            darkItemSelectedBg: '#313145',
            darkItemSelectedColor: '#CDD6F4',
            darkItemHoverColor: '#CDD6F4',
            darkItemHoverBg: '#2B2B3B',
          },
          Card: {
            colorBgContainer: '#2B2B3B',
            colorBorderSecondary: '#45475A',
          },
          Table: {
            colorBgContainer: '#2B2B3B',
            headerBg: '#313145',
            headerColor: '#CDD6F4',
            rowHoverBg: '#313145',
            borderColor: '#45475A',
            colorText: '#CDD6F4',
          },
          Modal: {
            contentBg: '#2B2B3B',
            headerBg: '#2B2B3B',
            titleColor: '#CDD6F4',
          },
          Button: {
            colorBgContainer: '#313145',
            defaultBg: '#313145',
            defaultBorderColor: '#45475A',
            defaultColor: '#CDD6F4',
          },
          Input: {
            colorBgContainer: '#313145',
            activeBorderColor: '#89B4FA',
            hoverBorderColor: '#585B70',
            addonBg: '#313145',
          },
          Select: {
            colorBgContainer: '#313145',
            optionSelectedBg: '#313145',
            optionActiveBg: '#2B2B3B',
          },
          Drawer: {
            colorBgElevated: '#181825',
          },
          Pagination: {
            colorBgContainer: '#2B2B3B',
          },
          Statistic: {
            colorTextDescription: '#A6ADC8',
          },
          Form: {
            labelColor: '#CDD6F4',
          },
        },
      }),
    }}>
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route
              path="/"
              element={
                <ProtectedRoute>
                  <MainLayout />
                </ProtectedRoute>
              }
            >
              <Route index element={<DashboardPage />} />
              <Route path="items" element={<ItemsPage />} />
              <Route path="locations" element={<LocationsPage />} />
              <Route path="stock" element={<StockPage />} />
              <Route path="purchase-orders" element={<PurchaseOrdersPage />} />
              <Route path="sales-orders" element={<SalesOrdersPage />} />
              <Route path="transfers" element={<TransfersPage />} />
              <Route path="register" element={<RegisterPage />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    </ConfigProvider>
  );
}

export default App;
