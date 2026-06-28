import api from './client';

export interface Stock {
  id: string;
  item_id: string;
  location_id: string;
  quantity: number;
  reserved_quantity: number;
  batch_number: string;
  item_sku?: string;
  item_name?: string;
  location_code?: string;
}

export interface AdjustStockRequest {
  item_id: string;
  location_id: string;
  quantity: number;
  notes?: string;
}

export interface StockMovement {
  id: string;
  item_id: string;
  from_location_id?: string;
  to_location_id?: string;
  quantity: number;
  movement_type: string;
  reference_type: string;
  reference_id: string;
  notes: string;
  created_by: string;
}

export const stockApi = {
  list: (page = 1, limit = 10) =>
    api.get('/stock', { params: { page, limit } }),
  getByItem: (itemId: string) =>
    api.get<Stock[]>('/stock', { params: { item_id: itemId } }),
  adjust: (data: AdjustStockRequest) => api.post('/stock/adjust', data),
  listMovements: (page = 1, limit = 10) =>
    api.get('/stock/movements', { params: { page, limit } }),
};
