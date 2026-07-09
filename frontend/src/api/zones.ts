import api from './client';

export interface Zone {
  id: string;
  code: string;
  name: string;
  description: string;
  is_active: boolean;
}

export interface CreateZoneRequest {
  code: string;
  name: string;
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

export const zoneApi = {
  list: (page = 1, limit = 10) =>
    api.get<PaginatedResponse<Zone>>('/zones', { params: { page, limit } }),
  listAll: () => api.get<Zone[]>('/zones/all'),
  get: (id: string) => api.get<Zone>(`/zones/${id}`),
  create: (data: CreateZoneRequest) => api.post<Zone>('/zones', data),
  update: (id: string, data: Partial<Zone>) => api.put<Zone>(`/zones/${id}`, data),
  delete: (id: string) => api.delete(`/zones/${id}`),
};
