<template>
  <div>
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="6"><el-card shadow="hover"><el-statistic title="严重" :value="stats.critical" /></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover"><el-statistic title="警告" :value="stats.warning" /></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover"><el-statistic title="提示" :value="stats.info" /></el-card></el-col>
      <el-col :span="6"><el-card shadow="hover"><el-statistic title="已解决" :value="stats.resolved" /></el-card></el-col>
    </el-row>

    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; gap: 12px;">
          <el-select v-model="filters.severity" placeholder="级别" clearable size="small" style="width: 120px;" @change="loadData">
            <el-option label="严重" value="critical" /><el-option label="警告" value="warning" /><el-option label="提示" value="info" />
          </el-select>
          <el-select v-model="filters.status" placeholder="状态" clearable size="small" style="width: 120px;" @change="loadData">
            <el-option label="已触发" value="triggered" /><el-option label="已确认" value="confirmed" /><el-option label="已解决" value="resolved" />
          </el-select>
        </div>
      </template>
      <el-table :data="alerts" v-loading="loading" stripe>
        <el-table-column label="级别" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.severity === 'critical'" type="danger" size="small">严重</el-tag>
            <el-tag v-else-if="row.severity === 'warning'" type="warning" size="small">警告</el-tag>
            <el-tag v-else type="info" size="small">提示</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="device?.name" label="设备" min-width="150" />
        <el-table-column prop="alert_type" label="类型" width="120" />
        <el-table-column prop="message" label="消息" min-width="200" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'triggered'" type="danger" size="small">待处理</el-tag>
            <el-tag v-else-if="row.status === 'confirmed'" type="warning" size="small">已确认</el-tag>
            <el-tag v-else type="success" size="small">已解决</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="triggered_at" label="触发时间" width="180" />
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status === 'triggered'" type="primary" link size="small" @click="handleConfirm(row)">确认</el-button>
            <el-button v-if="row.status !== 'resolved'" type="success" link size="small" @click="handleResolve(row)">解决</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total"
        layout="total, prev, pager, next" style="margin-top: 16px; justify-content: flex-end; display: flex;"
        @current-change="loadData" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getAlerts, getAlertStats, confirmAlert, resolveAlert } from '@/api/device'
import { ElMessage } from 'element-plus'

const alerts = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const stats = ref({ critical: 0, warning: 0, info: 0, triggered: 0, resolved: 0 })

const filters = reactive({ severity: '', status: '' })

async function loadData() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: page.value, page_size: pageSize.value }
    if (filters.severity) params.severity = filters.severity
    if (filters.status) params.status = filters.status
    const res = await getAlerts(params)
    alerts.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch { alerts.value = [] }
  finally { loading.value = false }
}

async function loadStats() {
  try { const res = await getAlertStats(); stats.value = res.data } catch { }
}

async function handleConfirm(row: any) {
  try { await confirmAlert(row.id); ElMessage.success('已确认'); loadData() } catch { }
}
async function handleResolve(row: any) {
  try { await resolveAlert(row.id); ElMessage.success('已解决'); loadData() } catch { }
}

onMounted(() => { loadData(); loadStats() })
</script>
