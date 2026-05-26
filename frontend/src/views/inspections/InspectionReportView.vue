<template>
  <div>
    <el-card shadow="never">
      <template #header>
        <span>巡检报告</span>
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
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getLatestInspections } from '@/api/device'

const inspections = ref<any[]>([])
const loading = ref(false)
const summary = ref({ total: 0, normal: 0, anomaly: 0 })

async function loadData() {
  loading.value = true
  try {
    const res = await getLatestInspections(50)
    inspections.value = res.data?.results || []
    summary.value = res.data?.summary || { total: 0, normal: 0, anomaly: 0 }
  } catch { inspections.value = [] }
  finally { loading.value = false }
}

onMounted(loadData)
</script>
