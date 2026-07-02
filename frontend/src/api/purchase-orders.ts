import api from './client';

export interface PurchaseOrder {
  id: string;
  po_number: string;
  supplier_name: string;
  status: string;
  expected_date?: string;
  notes?: string;
  created_by: string;
  created_by_name?: string;
  created_at: string;
  items?: POItem[];
}

export interface POItem {
  id: string;
  item_id: string;
  expected_quantity: number;
  received_quantity: number;
  unit_price: number;
}

export interface CreatePORequest {
  supplier_name: string;
  expected_date?: string;
  notes?: string;
  items: CreatePOItemRequest[];
}

export interface CreatePOItemRequest {
  item_id: string;
  expected_quantity: number;
  unit_price?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    page: number;
    limit: number;
    total: number;
  };
}

export interface ReceiveGoodsRequest {
  grn_number: string;
  notes?: string;
  items: ReceiveItemRequest[];
}

export interface ReceiveItemRequest {
  item_id: string;
  quantity: number;
  batch_number?: string;
  location_id: string;
}

export const poApi = {
  list: (page = 1, limit = 10) =>
    api.get<PaginatedResponse<PurchaseOrder>>('/purchase-orders', { params: { page, limit } }),
  get: (id: string) => api.get<PurchaseOrder>(`/purchase-orders/${id}`),
  create: (data: CreatePORequest) => api.post<PurchaseOrder>('/purchase-orders', data),
  receive: (id: string, data: ReceiveGoodsRequest) =>
    api.post(`/purchase-orders/${id}/receive`, data),
};
