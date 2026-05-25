<template>
  <div class="detail-page" v-if="device">
    <div class="detail-header">
      <button class="back-btn" @click="goBack">
        <el-icon><ArrowLeft /></el-icon> 返回
      </button>
      <div class="header-info">
        <h2 class="device-name">{{ device.name }}</h2>
        <span class="status-badge" :class="device.status ? 'online' : 'offline'">
          <span class="dot"></span>
          {{ device.status ? '在线' : '离线' }}
        </span>
      </div>
      <div class="header-actions">
        <el-button type="primary" size="small" @click="handleEdit">编辑设备</el-button>
        <el-button size="small" @click="handleInspect" :loading="inspecting">执行巡检</el-button>
        <el-button size="small" @click="handleBackup">配置备份</el-button>
      </div>
    </div>

    <div class="detail-grid">
      <div class="info-card">
        <div class="card-label">设备信息</div>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-key">设备名称</span>
            <span class="info-val">{{ device.name }}</span>
          </div>
          <div class="info-item">
            <span class="info-key">管理IP</span>
            <span class="info-val mono">{{ device.management_ip }}</span>
          </div>
          <div class="info-item">
            <span class="info-key">SSH端口</span>
            <span class="info-val mono">{{ device.ssh_port }}</span>
          </div>
          <div class="info-item">
            <span class="info-key">设备类型</span>
            <span class="info-val">{{ device.device_type }}</span>
          </div>
          <div class="info-item">
            <span class="info-key">厂商</span>
            <span class="info-val">{{ (device.vendor || '').toUpperCase() || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-key">型号</span>
            <span class="info-val mono">{{ device.model || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-key">角色</span>
            <span class="info-val"><span class="role-tag" :class="'role-'+device.role">{{ roleLabel(device.role) }}</span></span>
          </div>
          <div class="info-item">
            <span class="info-key">楼栋 / 楼层</span>
            <span class="info-val">{{ device.building || '-' }} {{ device.floor != null ? device.floor + 'F' : '' }}</span>
          </div>
          <div class="info-item" v-if="device.ap_name">
            <span class="info-key">AP名称</span>
            <span class="info-val mono">{{ device.ap_name }}</span>
          </div>
          <div class="info-item">
            <span class="info-key">创建时间</span>
            <span class="info-val mono-sm">{{ device.created_at }}</span>
          </div>
        </div>
      </div>

      <div class="info-card">
        <div class="card-label">最近巡检</div>
        <div v-if="inspectionHistory.length" class="mini-list">
          <div v-for="r in inspectionHistory.slice(0, 5)" :key="r.id" class="mini-item">
            <div class="mini-top">
              <span class="mini-time">{{ r.inspected_at?.substring(0, 16) || '-' }}</span>
              <span v-if="r.status === 'success'" class="mini-status ok">正常</span>
              <span v-else class="mini-status fail">失败</span>
            </div>
            <div class="mini-bars">
              <span class="mini-bar-label">CPU</span>
              <div class="mini-bar-track"><div class="mini-bar-fill" :style="{ width: (r.cpu_usage || 0) + '%', background: (r.cpu_usage || 0) > 80 ? '#ef4444' : '#00d4ff' }"></div></div>
              <span class="mini-bar-label">MEM</span>
              <div class="mini-bar-track"><div class="mini-bar-fill" :style="{ width: (r.memory_usage || 0) + '%', background: (r.memory_usage || 0) > 80 ? '#ef4444' : '#7c3aed' }"></div></div>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">暂无巡检数据</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getDevice, inspectDevice } from '@/api/device'
import { backupConfig } from '@/api/config'
import { ElMessage } from 'element-plus'

const route = useRoute(); const router = useRouter()
const device = ref<any>(null)
const inspectionHistory = ref<any[]>([])
const inspecting = ref(false)

function roleLabel(role: string) {
  const map: Record<string, string> = { core: '核心', ac: 'AC', aggregation: '本体', access: '接入', ap: 'AP' }
  return map[role] || role
}

async function loadDevice() {
  const id = Number(route.params.id)
  if (!id) return
  try { const res = await getDevice(id); device.value = res.data
  } catch { ElMessage.error('设备不存在'); router.push('/devices') }
}

async function loadHistory() {
  const id = Number(route.params.id)
  if (!id) return
  try {
    const { getInspectionReport } = await import('@/api/device')
    const res = await getInspectionReport({ device_id: id, page_size: 30 })
    inspectionHistory.value = res.data?.list || []
  } catch { }
}

function goBack() { router.push('/devices') }
function handleEdit() { router.push(`/devices/${device.value!.id}/edit`) }
async function handleInspect() {
  inspecting.value = true
  try { await inspectDevice(device.value!.id); ElMessage.success('巡检任务已创建') } catch { }
  finally { inspecting.value = false }
}
async function handleBackup() {
  try { await backupConfig([device.value!.id]); ElMessage.success('备份任务已创建') } catch { }
}

onMounted(() => { loadDevice(); loadHistory() })
</script>

<style scoped>
.detail-page { display: flex; flex-direction: column; gap: 20px; }

.detail-header {
  display: flex; align-items: center; gap: 20px;
  padding: 16px 20px;
  background: rgba(26, 35, 59, 0.7); backdrop-filter: blur(16px);
  border: 1px solid rgba(0, 212, 255, 0.06); border-radius: 12px;
}
.back-btn {
  display: flex; align-items: center; gap: 4px;
  background: none; border: 1px solid rgba(0,212,255,0.1); color: #94a3b8;
  padding: 6px 12px; border-radius: 6px; cursor: pointer; font-size: 13px;
  &:hover { border-color: rgba(0,212,255,0.3); color: #e2e8f0; }
}
.header-info { flex: 1; display: flex; align-items: center; gap: 16px; }
.device-name { font-size: 20px; font-weight: 700; color: #e2e8f0; margin: 0; }
.header-actions { display: flex; gap: 8px; }

.status-badge {
  display: flex; align-items: center; gap: 6px; font-size: 12px; padding: 4px 12px;
  border-radius: 6px; font-weight: 500;
}
.status-badge.online { background: rgba(16,185,129,0.1); color: #10b981; }
.status-badge.offline { background: rgba(239,68,68,0.1); color: #ef4444; }
.status-badge .dot { width: 6px; height: 6px; border-radius: 50%; }
.status-badge.online .dot { background: #10b981; box-shadow: 0 0 6px rgba(16,185,129,0.6); }
.status-badge.offline .dot { background: #ef4444; box-shadow: 0 0 6px rgba(239,68,68,0.4); }

.detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }

.info-card {
  padding: 20px;
  background: rgba(26, 35, 59, 0.7); backdrop-filter: blur(16px);
  border: 1px solid rgba(0, 212, 255, 0.06); border-radius: 12px;
}
.card-label { font-size: 13px; font-weight: 600; color: #e2e8f0; margin-bottom: 16px; letter-spacing: 0.5px; }

.info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.info-item { display: flex; flex-direction: column; gap: 2px; }
.info-key { font-size: 11px; color: #64748b; text-transform: uppercase; letter-spacing: 0.5px; }
.info-val { font-size: 14px; color: #e2e8f0; }
.info-val.mono { font-family: 'JetBrains Mono', monospace; font-size: 13px; }
.info-val.mono-sm { font-family: 'JetBrains Mono', monospace; font-size: 12px; color: #64748b; }

.role-tag { font-size: 11px; padding: 2px 10px; border-radius: 4px; font-weight: 500; }
.role-core { background: rgba(239,68,68,0.1); color: #ef4444; }
.role-ac { background: rgba(245,158,11,0.1); color: #f59e0b; }
.role-aggregation { background: rgba(0,212,255,0.08); color: #00d4ff; }
.role-access { background: rgba(100,116,139,0.1); color: #94a3b8; }
.role-ap { background: rgba(16,185,129,0.1); color: #10b981; }

/* Mini list */
.mini-list { display: flex; flex-direction: column; gap: 12px; }
.mini-item {
  padding: 12px; border-radius: 8px; background: rgba(0,0,0,0.15);
}
.mini-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.mini-time { font-size: 12px; color: #64748b; font-family: 'JetBrains Mono', monospace; }
.mini-status { font-size: 11px; padding: 1px 8px; border-radius: 3px; font-weight: 500; }
.mini-status.ok { background: rgba(16,185,129,0.1); color: #10b981; }
.mini-status.fail { background: rgba(239,68,68,0.1); color: #ef4444; }
.mini-bars { display: flex; align-items: center; gap: 8px; }
.mini-bar-label { font-size: 10px; color: #64748b; width: 30px; }
.mini-bar-track { flex: 1; height: 4px; background: rgba(255,255,255,0.04); border-radius: 2px; overflow: hidden; }
.mini-bar-fill { height: 100%; border-radius: 2px; transition: width 0.6s; }

.empty-state { padding: 32px; text-align: center; color: #475569; font-size: 13px; }
</style>
