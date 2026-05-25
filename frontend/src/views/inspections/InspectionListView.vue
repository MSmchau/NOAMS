<template>
  <div>
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span>巡检记录</span>
          <el-button type="primary" size="small" @click="handleRefresh">
            <el-icon><Refresh /></el-icon> 刷新
          </el-button>
        </div>
      </template>
      <el-table :data="inspections" v-loading="loading" stripe>
        <el-table-column prop="device?.name" label="设备" min-width="150" />
        <el-table-column label="CPU" width="100">
          <template #default="{ row }"> <el-tag v-if="row.cpu_usage" size="small">{{ row.cpu_usage }}%</el-tag><span v-else>-</span> </template>
        </el-table-column>
        <el-table-column label="内存" width="100">
          <template #default="{ row }"> <el-tag v-if="row.memory_usage" size="small">{{ row.memory_usage }}%</el-tag><span v-else>-</span> </template>
        </el-table-column>
        <el-table-column prop="uptime" label="运行时长" min-width="120" />
        <el-table-column label="异常" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_anomaly" type="danger" size="small">是</el-tag>
            <el-tag v-else type="success" size="small">否</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="inspected_at" label="巡检时间" width="180" />
      </el-table>

      <el-pagination
        v-model:current-page="page" :page-size="pageSize" :total="total"
        layout="total, prev, pager, next" style="margin-top: 16px; justify-content: flex-end; display: flex;"
        @current-change="loadData"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getInspectionReport } from '@/api/device'

const inspections = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

async function loadData() {
  loading.value = true
  try {
    const res = await getInspectionReport({ page: page.value, page_size: pageSize.value })
    inspections.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch { inspections.value = [] }
  finally { loading.value = false }
}

function handleRefresh() { page.value = 1; loadData() }
onMounted(loadData)
</script>
