<template>
  <div>
    <el-card shadow="never">
      <template #header>
        <div style="display:flex;justify-content:space-between;align-items:center;">
          <span>巡检报告</span>
          <div style="display:flex;gap:8px;">
            <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" size="small" style="width:240px" @change="loadData" />
            <el-button size="small" @click="handleExport" :loading="exporting">导出 CSV</el-button>
            <el-button size="small" @click="handleRefresh">刷新</el-button>
          </div>
        </div>
      </template>
      <el-row :gutter="20" style="margin-bottom: 20px;">
        <el-col :span="8"><el-statistic title="总巡检次数" :value="summary.total" /></el-col>
        <el-col :span="8"><el-statistic title="正常" :value="summary.normal" /></el-col>
        <el-col :span="8"><el-statistic title="异常" :value="summary.anomaly" /></el-col>
      </el-row>
      <el-table :data="inspections" v-loading="loading" stripe>
        <el-table-column label="设备" min-width="150">
          <template #default="{ row }">{{ row.device?.name || '-' }}</template>
        </el-table-column>
        <el-table-column label="IP" width="140">
          <template #default="{ row }">{{ row.device?.management_ip || '-' }}</template>
        </el-table-column>
        <el-table-column label="CPU" width="80"><template #default="{ row }">{{ row.cpu_usage ?? '-' }}%</template></el-table-column>
        <el-table-column label="内存" width="80"><template #default="{ row }">{{ row.memory_usage ?? '-' }}%</template></el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_anomaly" type="danger" size="small">异常</el-tag>
            <el-tag v-else type="success" size="small">正常</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="inspected_at" label="时间" width="180" />
      </el-table>
      <div v-if="total > pageSize" style="display:flex;justify-content:flex-end;margin-top:16px;">
        <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadData" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getInspectionReport } from '@/api/device'
import request from '@/api/request'
import { ElMessage } from 'element-plus'

const inspections = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const summary = ref({ total: 0, normal: 0, anomaly: 0 })
const dateRange = ref<[Date, Date] | null>(null)
const exporting = ref(false)

async function loadData() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: page.value, page_size: pageSize.value }
    if (dateRange.value) {
      params.start_date = dateRange.value[0].toISOString().substring(0, 10)
      params.end_date = dateRange.value[1].toISOString().substring(0, 10)
    }
    const res = await getInspectionReport(params)
    inspections.value = res.data?.list || []
    total.value = res.data?.total || 0
    // 从 latest 接口获取汇总统计
    const latestRes = await request.get('/inspections/latest', { params: { limit: 1 } })
    summary.value = latestRes.data?.summary || { total: 0, normal: 0, anomaly: 0 }
  } catch { inspections.value = [] }
  finally { loading.value = false }
}

function handleRefresh() {
  page.value = 1
  loadData()
}

async function handleExport() {
  exporting.value = true
  try {
    const params: Record<string, any> = {}
    if (dateRange.value) {
      params.start_date = dateRange.value[0].toISOString().substring(0, 10)
      params.end_date = dateRange.value[1].toISOString().substring(0, 10)
    }
    const res = await request.get('/inspections/export', {
      params,
      responseType: 'blob',
    })
    const url = window.URL.createObjectURL(new Blob([res as any], { type: 'text/csv;charset=utf-8;' }))
    const a = document.createElement('a')
    a.href = url
    a.download = `inspection_report_${new Date().toISOString().substring(0, 10)}.csv`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch {
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

onMounted(loadData)
</script>
