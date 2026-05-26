import request from './request'

export interface Credential {
  id: number
  name: string
  username: string
  auth_method: string
  description: string
  created_at: string
  updated_at: string
}

export function getCredentials() {
  return request.get<Credential[]>('/credentials')
}

export function getCredential(id: number) {
  return request.get<Credential>(`/credentials/${id}`)
}

export function createCredential(data: {
  name: string
  username: string
  password: string
  enable_pw?: string
  auth_method?: string
  description?: string
}) {
  return request.post<Credential>('/credentials', data)
}

export function updateCredential(id: number, data: Partial<Credential>) {
  return request.put(`/credentials/${id}`, data)
}

export function deleteCredential(id: number) {
  return request.delete(`/credentials/${id}`)
}
