<template>
  <el-card shadow="hover">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <span>配置备份</span>
        <div>
          <el-select v-model="selectedDevice" placeholder="选择设备" clearable style="width: 200px; margin-right: 12px;">
            <el-option v-for="d in devices" :key="d.id" :label="d.name" :value="d.id" />
          </el-select>
          <el-button type="primary" size="small" @click="handleBackup" :loading="backingUp">立即备份</el-button>
        </div>
      </div>
    </template>
    <el-empty v-if="!selectedDevice" description="请选择设备查看备份历史" />
    <el-table v-else :data="history" v-loading="loading" stripe>
      <el-table-column prop="id" label="版本" width="80" />
      <el-table-column prop="created_at" label="备份时间" width="180" />
      <el-table-column prop="triggered_by" label="触发方式" width="120" />
      <el-table-column prop="operator" label="操作人" width="120" />
      <el-table-column prop="config_hash" label="配置Hash" width="120">
        <template #default="{ row }">{{ row.config_hash?.substring(0, 12) || '-' }}</template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { getAllDevices } from '@/api/device'
import { backupConfig, getConfigHistory } from '@/api/config'
import { ElMessage } from 'element-plus'

const devices = ref<any[]>([])
const selectedDevice = ref<number>()
const history = ref<any[]>([])
const loading = ref(false)
const backingUp = ref(false)

async function loadDevices() {
  try { const res = await getAllDevices(); devices.value = res.data || [] } catch { }
}

async function loadHistory() {
  if (!selectedDevice.value) return
  loading.value = true
  try { const res = await getConfigHistory(selectedDevice.value); history.value = res.data?.list || [] }
  catch { history.value = [] }
  finally { loading.value = false }
}

async function handleBackup() {
  if (!selectedDevice.value) { ElMessage.warning('请先选择设备'); return }
  backingUp.value = true
  try { await backupConfig([selectedDevice.value]); ElMessage.success('备份任务已创建'); loadHistory() }
  catch { }
  finally { backingUp.value = false }
}

watch(selectedDevice, () => loadHistory())
onMounted(loadDevices)
</script>
