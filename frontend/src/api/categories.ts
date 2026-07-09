import api from './client';

export interface Category {
  id: string;
  name: string;
  abbreviation: string;
  description: string;
  is_active: boolean;
}

export interface CreateCategoryRequest {
  name: string;
  abbreviation: string;
  description?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  meta: {
    page: number;
    limit: number;
    total: number;
  };
}

export const categoryApi = {
  list: (page = 1, limit = 10) =>
    api.get<PaginatedResponse<Category>>('/categories', { params: { page, limit } }),
  listAll: () => api.get<Category[]>('/categories/all'),
  get: (id: string) => api.get<Category>(`/categories/${id}`),
  create: (data: CreateCategoryRequest) => api.post<Category>('/categories', data),
  update: (id: string, data: Partial<Category>) => api.put<Category>(`/categories/${id}`, data),
  delete: (id: string) => api.delete(`/categories/${id}`),
};
