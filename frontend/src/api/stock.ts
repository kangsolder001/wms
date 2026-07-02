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

export interface StockOpnameRequest {
  location_id: string;
  notes?: string;
  items: StockOpnameItemRequest[];
}

export interface StockOpnameItemRequest {
  item_id: string;
  system_quantity: number;
  actual_quantity: number;
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
  list: (page = 1, limit = 10, itemId?: string, locationId?: string, search?: string) => {
    const params: any = { page, limit };
    if (itemId) params.item_id = itemId;
    if (locationId) params.location_id = locationId;
    if (search) params.search = search;
    return api.get('/stock', { params });
  },
  getByItem: (itemId: string) =>
    api.get<Stock[]>('/stock', { params: { item_id: itemId } }),
  adjust: (data: AdjustStockRequest) => api.post('/stock/adjust', data),
  opname: (data: StockOpnameRequest) => api.post('/stock/opname', data),
  listMovements: (page = 1, limit = 10) =>
    api.get('/stock/movements', { params: { page, limit } }),
};
