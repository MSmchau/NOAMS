<template>
  <el-card shadow="never">
    <template #header><span>执行日志</span></template>
    <el-table :data="logs" v-loading="loading" stripe>
      <el-table-column prop="task_id" label="任务ID" width="160" />
      <el-table-column label="设备" min-width="150">
        <template #default="{ row }">{{ row.device?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="IP" width="140">
        <template #default="{ row }">{{ row.device?.management_ip || '-' }}</template>
      </el-table-column>
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.status === 'success'" type="success" size="small">成功</el-tag>
          <el-tag v-else-if="row.status === 'failed'" type="danger" size="small">失败</el-tag>
          <el-tag v-else type="warning" size="small">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="duration" label="耗时(s)" width="80">
        <template #default="{ row }">{{ row.duration || '-' }}</template>
      </el-table-column>
      <el-table-column label="异常" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.is_anomaly" type="danger" size="small">是</el-tag>
          <el-tag v-else type="success" size="small">否</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="180" />
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'

const logs = ref<any[]>([])
const loading = ref(false)

async function loadData() {
  loading.value = true
  try {
    const res = await request.get('/tasks/logs')
    logs.value = res.data?.list || []
  } catch { logs.value = [] }
  finally { loading.value = false }
}

onMounted(loadData)
</script>
