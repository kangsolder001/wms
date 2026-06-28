import api from './client';

export interface Location {
  id: string;
  code: string;
  name: string;
  zone: string;
  aisle: string;
  rack: string;
  level: string;
  bin: string;
  type: string;
  capacity: number;
  is_active: boolean;
}

export interface CreateLocationRequest {
  code: string;
  name: string;
  zone?: string;
  aisle?: string;
  rack?: string;
  level?: string;
  bin?: string;
  type: string;
  capacity?: number;
}

export const locationApi = {
  list: (page = 1, limit = 10) =>
    api.get('/locations', { params: { page, limit } }),
  get: (id: string) => api.get<Location>(`/locations/${id}`),
  create: (data: CreateLocationRequest) => api.post<Location>('/locations', data),
  update: (id: string, data: Partial<Location>) => api.put<Location>(`/locations/${id}`, data),
  delete: (id: string) => api.delete(`/locations/${id}`),
};
