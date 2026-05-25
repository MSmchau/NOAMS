<template>
  <div class="device-detail">
    <el-card shadow="hover" v-if="device">
      <template #header>
        <div style="display: flex; align-items: center; gap: 12px;">
          <el-button @click="goBack">返回</el-button>
          <span style="font-size: 18px; font-weight: bold;">{{ device.name }}</span>
          <el-tag v-if="device.status === 1" type="success" size="small">在线</el-tag>
          <el-tag v-else type="danger" size="small">离线</el-tag>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="设备名称" :span="2">{{ device.name }}</el-descriptions-item>
        <el-descriptions-item label="管理IP">{{ device.management_ip }}</el-descriptions-item>
        <el-descriptions-item label="SSH端口">{{ device.ssh_port }}</el-descriptions-item>
        <el-descriptions-item label="设备类型">{{ device.device_type }}</el-descriptions-item>
        <el-descriptions-item label="厂商">{{ device.vendor?.toUpperCase() }}</el-descriptions-item>
        <el-descriptions-item label="型号">{{ device.model }}</el-descriptions-item>
        <el-descriptions-item label="角色">{{ roleLabel(device.role) }}</el-descriptions-item>
        <el-descriptions-item label="楼栋">{{ device.building || '-' }}</el-descriptions-item>
        <el-descriptions-item label="楼层">{{ device.floor != null ? device.floor + '层' : '-' }}</el-descriptions-item>
        <el-descriptions-item label="AP名称">{{ device.ap_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="最后在线">{{ device.last_seen || '-' }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ device.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ device.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ device.updated_at }}</el-descriptions-item>
      </el-descriptions>

      <div style="margin-top: 20px; display: flex; gap: 12px;">
        <el-button type="primary" @click="handleEdit">编辑设备</el-button>
        <el-button type="success" @click="handleInspect">执行巡检</el-button>
        <el-button type="warning" @click="handleBackup">配置备份</el-button>
      </div>
    </el-card>

    <el-card shadow="hover" style="margin-top: 20px;" v-if="inspectionHistory.length">
      <template #header>
        <span>巡检历史（最近30次）</span>
      </template>
      <el-table :data="inspectionHistory" stripe>
        <el-table-column prop="inspected_at" label="巡检时间" width="180" />
        <el-table-column prop="cpu_usage" label="CPU%" width="100" />
        <el-table-column prop="memory_usage" label="内存%" width="100" />
        <el-table-column prop="temperature" label="温度(℃)" width="100" />
        <el-table-column prop="uptime" label="运行时长" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status === 'success' ? '正常' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getDevice, inspectDevice } from '@/api/device'
import { ElMessage } from 'element-plus'
import type { Device } from '@/api/device'

const route = useRoute()
const router = useRouter()
const device = ref<Device | null>(null)
const inspectionHistory = ref<any[]>([])

function roleLabel(role: string) {
  const map: Record<string, string> = {
    core: '核心交换机',
    ac: '无线控制器',
    aggregation: '本体交换机',
    access: '接入交换机',
    ap: '无线AP',
  }
  return map[role] || role
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

async function loadInspectionHistory() {
  const id = Number(route.params.id)
  if (!id) return
  try {
    const { getInspectionReport } = await import('@/api/device')
    const res = await getInspectionReport({ device_id: id, page_size: 30 })
    inspectionHistory.value = res.data?.list || []
  } catch {
    // silently fail
  }
}

function goBack() {
  router.push('/devices')
}

function handleEdit() {
  router.push(`/devices/${device.value!.id}/edit`)
}

async function handleInspect() {
  try {
    await inspectDevice(device.value!.id)
    ElMessage.success('巡检任务已创建')
  } catch {
    // handled by interceptor
  }
}

async function handleBackup() {
  try {
    const { backupConfig } = await import('@/api/config')
    await backupConfig([device.value!.id])
    ElMessage.success('备份任务已创建')
  } catch {
    // handled by interceptor
  }
}

onMounted(() => {
  loadDevice()
  loadInspectionHistory()
})
</script>
