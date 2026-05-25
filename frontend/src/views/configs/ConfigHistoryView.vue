<template>
  <el-card shadow="hover">
    <template #header><span>版本历史</span></template>
    <el-table :data="backups" v-loading="loading" stripe>
      <el-table-column prop="device?.name" label="设备" min-width="150" />
      <el-table-column prop="device?.management_ip" label="IP" width="140" />
      <el-table-column prop="created_at" label="备份时间" width="180" />
      <el-table-column prop="triggered_by" label="触发方式" width="100" />
      <el-table-column prop="git_commit_id" label="Git提交" width="100">
        <template #default="{ row }">{{ row.git_commit_id?.substring(0, 8) || '-' }}</template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small">查看</el-button>
          <el-button type="warning" link size="small">对比</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'

const backups = ref<any[]>([])
const loading = ref(false)

async function loadData() {
  loading.value = true
  try {
    const res = await request.get('/configs/history/all')
    backups.value = res.data?.list || []
  } catch { backups.value = [] }
  finally { loading.value = false }
}

onMounted(loadData)
</script>
