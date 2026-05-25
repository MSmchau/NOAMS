import request from './request'

export function backupConfig(deviceIds: number[]) {
  return request.post('/configs/backup', { device_ids: deviceIds })
}

export function getConfigHistory(deviceId: number, params?: Record<string, any>) {
  return request.get(`/configs/history/${deviceId}`, { params })
}

export function rollbackConfig(backupId: number) {
  return request.post('/configs/rollback', { backup_id: backupId })
}

export function diffConfig(id1: number, id2: number) {
  return request.get('/configs/diff', { params: { id1, id2 } })
}
