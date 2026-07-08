import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ConfigProvider, theme } from 'antd';
import { useAuthStore } from './stores/authStore';
import { useThemeStore } from './stores/themeStore';
import LoginPage from './features/auth/LoginPage';
import UsersPage from './features/auth/UsersPage';
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
  colorBgBase: '#0d1117',
  colorBgContainer: '#161b22',
  colorBgElevated: '#1c2128',
  colorBgLayout: '#010409',
  colorBgSpotlight: '#1c2128',
  colorText: '#c9d1d9',
  colorTextSecondary: '#8b949e',
  colorTextTertiary: '#6e7681',
  colorTextQuaternary: '#484f58',
  colorBorder: '#30363d',
  colorBorderSecondary: '#21262d',
  colorFill: '#1c2128',
  colorFillSecondary: '#161b22',
  colorFillTertiary: '#0d1117',
  colorFillQuaternary: '#010409',
  controlItemBgActive: '#1c2128',
  controlItemBgHover: '#161b22',
  colorPrimary: '#58a6ff',
  colorPrimaryBg: '#0d1117',
  colorPrimaryBgHover: '#161b22',
  colorPrimaryBorder: '#58a6ff',
  colorPrimaryBorderHover: '#79c0ff',
  colorPrimaryHover: '#79c0ff',
  colorPrimaryActive: '#a5d6ff',
  colorPrimaryTextHover: '#79c0ff',
  colorPrimaryText: '#58a6ff',
  colorPrimaryTextActive: '#a5d6ff',
  colorErrorBg: '#3d1f23',
  colorErrorBorder: '#f85149',
  colorSuccessBg: '#1a3a2a',
  colorSuccessBorder: '#3fb950',
  colorWarningBg: '#3d2e1a',
  colorWarningBorder: '#d29922',
};

const lightThemeTokens = {
  colorBgBase: '#ffffff',
  colorBgContainer: '#f6f8fa',
  colorBgElevated: '#ffffff',
  colorBgLayout: '#f6f8fa',
  colorBgSpotlight: '#ffffff',
  colorText: '#24292f',
  colorTextSecondary: '#57606a',
  colorTextTertiary: '#6e7781',
  colorTextQuaternary: '#8c959f',
  colorBorder: '#d0d7de',
  colorBorderSecondary: '#d8dee4',
  colorFill: '#f6f8fa',
  colorFillSecondary: '#ffffff',
  colorFillTertiary: '#f6f8fa',
  colorFillQuaternary: '#ffffff',
  controlItemBgActive: '#f3f4f6',
  controlItemBgHover: '#f6f8fa',
  colorPrimary: '#0969da',
  colorPrimaryBg: '#f0f7ff',
  colorPrimaryBgHover: '#dbe9ff',
  colorPrimaryBorder: '#0969da',
  colorPrimaryBorderHover: '#54aeff',
  colorPrimaryHover: '#54aeff',
  colorPrimaryActive: '#0550ae',
  colorPrimaryTextHover: '#54aeff',
  colorPrimaryText: '#0969da',
  colorPrimaryTextActive: '#0550ae',
  colorErrorBg: '#fff1f0',
  colorErrorBorder: '#cf222e',
  colorSuccessBg: '#dafbe1',
  colorSuccessBorder: '#1a7f37',
  colorWarningBg: '#fff8c5',
  colorWarningBorder: '#9a6700',
};

function App() {
  const isDark = useThemeStore((s) => s.isDark);

  const themeConfig = isDark ? {
    token: darkThemeTokens,
    algorithm: theme.darkAlgorithm,
    components: {
      Layout: {
        headerBg: '#010409',
        siderBg: '#010409',
        bodyBg: '#0d1117',
        headerColor: '#c9d1d9',
        triggerBg: '#010409',
        triggerColor: '#c9d1d9',
      },
      Menu: {
        darkItemBg: '#010409',
        darkSubMenuItemBg: '#010409',
        darkItemColor: '#8b949e',
        darkItemSelectedBg: '#1c2128',
        darkItemSelectedColor: '#c9d1d9',
        darkItemHoverColor: '#c9d1d9',
        darkItemHoverBg: '#161b22',
      },
      Card: {
        colorBgContainer: '#161b22',
        colorBorderSecondary: '#30363d',
      },
      Table: {
        colorBgContainer: '#161b22',
        headerBg: '#1c2128',
        headerColor: '#c9d1d9',
        rowHoverBg: '#1c2128',
        borderColor: '#30363d',
        colorText: '#c9d1d9',
      },
      Modal: {
        contentBg: '#161b22',
        headerBg: '#161b22',
        titleColor: '#c9d1d9',
      },
      Button: {
        colorBgContainer: '#1c2128',
        defaultBg: '#1c2128',
        defaultBorderColor: '#30363d',
        defaultColor: '#c9d1d9',
      },
      Input: {
        colorBgContainer: '#1c2128',
        activeBorderColor: '#58a6ff',
        hoverBorderColor: '#484f58',
        addonBg: '#1c2128',
      },
      Select: {
        colorBgContainer: '#1c2128',
        optionSelectedBg: '#1c2128',
        optionActiveBg: '#161b22',
      },
      Drawer: {
        colorBgElevated: '#010409',
      },
      Pagination: {
        colorBgContainer: '#161b22',
      },
      Statistic: {
        colorTextDescription: '#8b949e',
      },
      Form: {
        labelColor: '#c9d1d9',
      },
    },
  } : {
    token: lightThemeTokens,
    algorithm: theme.defaultAlgorithm,
    components: {
      Layout: {
        headerBg: '#f6f8fa',
        siderBg: '#f6f8fa',
        bodyBg: '#ffffff',
        headerColor: '#24292f',
        triggerBg: '#f6f8fa',
        triggerColor: '#24292f',
      },
      Menu: {
        darkItemBg: '#f6f8fa',
        darkSubMenuItemBg: '#f6f8fa',
        darkItemColor: '#57606a',
        darkItemSelectedBg: '#dbe9ff',
        darkItemSelectedColor: '#0969da',
        darkItemHoverColor: '#24292f',
        darkItemHoverBg: '#f3f4f6',
      },
    },
  };

  return (
    <ConfigProvider theme={themeConfig}>
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
              <Route path="users" element={<UsersPage />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    </ConfigProvider>
  );
}

export default App;
