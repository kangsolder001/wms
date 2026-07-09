import api from './client';

export interface Item {
  id: string;
  sku: string;
  name: string;
  description: string;
  category: string;
  category_id: string;
  barcode: string;
  unit_of_measure: string;
  weight: number;
  length: number;
  width: number;
  height: number;
  is_active: boolean;
}

export interface CreateItemRequest {
  category_id: string;
  name: string;
  description?: string;
  unit_of_measure: string;
  weight?: number;
  length?: number;
  width?: number;
  height?: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    page: number;
    limit: number;
    total: number;
  };
}

export const itemApi = {
  list: (page = 1, limit = 10) =>
    api.get<PaginatedResponse<Item>>('/items', { params: { page, limit } }),
  get: (id: string) => api.get<Item>(`/items/${id}`),
  generateSKU: (categoryId: string) =>
    api.post<{ sku: string }>('/items/generate-sku', { category_id: categoryId }),
  create: (data: CreateItemRequest) => api.post<Item>('/items', data),
  update: (id: string, data: Partial<Item>) => api.put<Item>(`/items/${id}`, data),
  delete: (id: string) => api.delete(`/items/${id}`),
};
