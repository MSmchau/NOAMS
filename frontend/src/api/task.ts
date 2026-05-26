import request from './request'

export function getTasks(params?: Record<string, any>) {
  return request.get('/tasks', { params })
}

export function createTask(data: Record<string, any>) {
  return request.post('/tasks', data)
}

export function updateTask(id: number, data: Record<string, any>) {
  return request.put(`/tasks/${id}`, data)
}

export function deleteTask(id: number) {
  return request.delete(`/tasks/${id}`)
}

export function toggleTask(id: number) {
  return request.put(`/tasks/${id}/toggle`)
}

export function getTaskLogs(params?: Record<string, any>) {
  return request.get('/tasks/logs', { params })
}
