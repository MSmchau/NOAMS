import request from './request'
import type { PageData } from './request'

export interface Device {
  id: number
  name: string
  device_type: string
  vendor: string
  model: string
  role: string
  management_ip: string
  ssh_port: number
  credential_id: number | null
  group_id: number | null
  building: string
  floor: number
  ap_name: string
  status: number
  last_seen: string | null
  description: string
  created_at: string
  updated_at: string
  group?: { id: number; name: string }
  credential?: { id: number; name: string }
}

export interface DeviceStats {
  total: number
  online: number
  offline: number
  role_counts: Record<string, number>
}

export function getDevices(params?: Record<string, any>) {
  return request.get<PageData<Device>>('/devices', { params })
}

export function getAllDevices() {
  return request.get<Device[]>('/devices/all')
}

export function getDevice(id: number) {
  return request.get<Device>(`/devices/${id}`)
}

export function createDevice(data: Partial<Device>) {
  return request.post<Device>('/devices', data)
}

export function updateDevice(id: number, data: Partial<Device>) {
  return request.put(`/devices/${id}`, data)
}

export function deleteDevice(id: number) {
  return request.delete(`/devices/${id}`)
}

export function getDeviceStats() {
  return request.get<DeviceStats>('/devices/stats')
}

export function inspectDevice(id: number) {
  return request.post(`/devices/${id}/inspect`)
}

export function batchInspect(deviceIds: number[]) {
  return request.post('/inspections/batch', { device_ids: deviceIds })
}

export function getInspectionReport(params?: Record<string, any>) {
  return request.get('/inspections/report', { params })
}

export function getLatestInspections(limit = 10) {
  return request.get('/inspections/latest', { params: { limit } })
}

export function getDashboard() {
  return request.get('/monitor/dashboard')
}

export function getAlerts(params?: Record<string, any>) {
  return request.get('/alerts', { params })
}

export function getAlertStats() {
  return request.get('/alerts/stats')
}

export function confirmAlert(id: number) {
  return request.put(`/alerts/${id}/confirm`)
}

export function resolveAlert(id: number) {
  return request.put(`/alerts/${id}/resolve`)
}
