<template>
  <div class="page-container" v-if="device">
    <!-- Header Card -->
    <el-card shadow="never">
      <div class="detail-header">
        <div class="header-left">
          <el-button text @click="goBack">
            <el-icon><ArrowLeft /></el-icon> 返回
          </el-button>
          <div class="header-divider"></div>
          <h2 class="device-name">{{ device.name }}</h2>
          <span class="status-badge" :class="device.status ? 'online' : 'offline'">
            <span class="dot"></span>
            {{ device.status ? '在线' : '离线' }}
          </span>
        </div>
        <div class="header-actions">
          <el-button size="small" @click="handleEdit">编辑设备</el-button>
          <el-button type="primary" size="small" @click="handleInspect" :loading="inspecting">执行巡检</el-button>
          <el-button size="small" @click="handleBackup">配置备份</el-button>
        </div>
      </div>
    </el-card>

    <!-- Info Grid -->
    <el-row :gutter="16">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>设备信息</template>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="设备名称" :span="2">{{ device.name }}</el-descriptions-item>
            <el-descriptions-item label="管理IP">
              <code style="font-size:12px">{{ device.management_ip }}</code>
            </el-descriptions-item>
            <el-descriptions-item label="SSH端口">
              <code style="font-size:12px">{{ device.ssh_port }}</code>
            </el-descriptions-item>
            <el-descriptions-item label="设备类型">{{ device.device_type }}</el-descriptions-item>
            <el-descriptions-item label="厂商">{{ (device.vendor || '').toUpperCase() || '-' }}</el-descriptions-item>
            <el-descriptions-item label="型号">{{ device.model || '-' }}</el-descriptions-item>
            <el-descriptions-item label="角色">
              <el-tag :type="roleTagType(device.role)" size="small">{{ roleLabel(device.role) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="楼栋/楼层">{{ device.building || '-' }} {{ device.floor != null ? device.floor + 'F' : '' }}</el-descriptions-item>
            <el-descriptions-item label="AP名称" v-if="device.ap_name">{{ device.ap_name }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ device.created_at }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>最近巡检</template>
          <div v-if="inspectionHistory.length" class="inspection-list">
            <div v-for="r in inspectionHistory.slice(0, 5)" :key="r.id" class="inspection-item">
              <div class="inspection-top">
                <span class="inspection-time">{{ r.inspected_at?.substring(0, 16) || '-' }}</span>
                <el-tag v-if="r.status === 'success'" type="success" size="small">正常</el-tag>
                <el-tag v-else type="danger" size="small">失败</el-tag>
              </div>
              <div class="inspection-bars">
                <div class="bar-item">
                  <span class="bar-label">CPU</span>
                  <el-progress :percentage="r.cpu_usage || 0" :stroke-width="6" :color="barColor(r.cpu_usage)" />
                </div>
                <div class="bar-item">
                  <span class="bar-label">MEM</span>
                  <el-progress :percentage="r.memory_usage || 0" :stroke-width="6" :color="barColor(r.memory_usage)" />
                </div>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无巡检数据" :image-size="60" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getDevice, inspectDevice, getInspectionReport } from '@/api/device'
import { backupConfig } from '@/api/config'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const device = ref<any>(null)
const inspectionHistory = ref<any[]>([])
const inspecting = ref(false)

function roleLabel(role: string) {
  const map: Record<string, string> = { core: '核心', ac: 'AC', aggregation: '本体', access: '接入', ap: 'AP' }
  return map[role] || role
}

function roleTagType(role?: string) {
  const map: Record<string, 'danger' | 'warning' | 'info' | 'success'> = { core: 'danger', ac: 'warning', aggregation: 'info', access: 'info', ap: 'success' }
  return map[role || ''] || 'info'
}

function barColor(val: number) {
  if (val > 80) return '#ff4d4f'
  if (val > 50) return '#fa8c16'
  return '#1890ff'
}

async function loadDevice() {
  const id = Number(route.params.id)
  if (!id) return
  try {
    const res = await getDevice(id)
    device.value = res.data
  } catch {
    ElMessage.error('设备不存在')
    router.push('/devices')
  }
}

async function loadHistory() {
  const id = Number(route.params.id)
  if (!id) return
  try {
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
.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-divider {
  width: 1px;
  height: 24px;
  background: #f0f0f0;
}

.device-name {
  font-size: 18px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
}

/* Inspection list */
.inspection-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.inspection-item {
  padding: 12px;
  background: #fafafa;
  border-radius: 6px;
  border: 1px solid #f0f0f0;
}

.inspection-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.inspection-time {
  font-size: 12px;
  color: #8c8c8c;
  font-family: 'JetBrains Mono', monospace;
}

.inspection-bars {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.bar-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.bar-label {
  font-size: 11px;
  color: #8c8c8c;
  width: 36px;
  flex-shrink: 0;
}

.bar-item :deep(.el-progress) {
  flex: 1;
}

.bar-item :deep(.el-progress-bar__outer) {
  background: #f0f0f0;
}
</style>
