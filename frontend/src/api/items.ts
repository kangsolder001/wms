import api from './client';

export interface Item {
  id: string;
  sku: string;
  name: string;
  description: string;
  category: string;
  unit_of_measure: string;
  weight: number;
  is_active: boolean;
}

export interface CreateItemRequest {
  sku: string;
  name: string;
  description?: string;
  category?: string;
  unit_of_measure: string;
  weight?: number;
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
  create: (data: CreateItemRequest) => api.post<Item>('/items', data),
  update: (id: string, data: Partial<Item>) => api.put<Item>(`/items/${id}`, data),
  delete: (id: string) => api.delete(`/items/${id}`),
};
