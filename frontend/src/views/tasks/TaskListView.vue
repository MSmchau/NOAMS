<template>
  <el-card shadow="hover">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <span>定时任务</span>
      </div>
    </template>
    <el-table :data="tasks" v-loading="loading" stripe>
      <el-table-column prop="name" label="任务名称" min-width="160" />
      <el-table-column prop="task_type" label="类型" width="120" />
      <el-table-column prop="cron_expr" label="Cron表达式" width="140" />
      <el-table-column label="状态" width="80">
        <template #default="{ row }">
          <el-tag v-if="row.status" type="success" size="small">启用</el-tag>
          <el-tag v-else type="info" size="small">停用</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="last_run_at" label="上次执行" width="180" />
      <el-table-column label="结果" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.last_result === 'success'" type="success" size="small">成功</el-tag>
          <el-tag v-else-if="row.last_result === 'failed'" type="danger" size="small">失败</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-switch :model-value="!!row.status" @change="handleToggle(row)" size="small" />
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'
import { ElMessage } from 'element-plus'

const tasks = ref<any[]>([])
const loading = ref(false)

async function loadData() {
  loading.value = true
  try {
    const res = await request.get('/tasks')
    tasks.value = res.data?.list || []
  } catch { tasks.value = [] }
  finally { loading.value = false }
}

async function handleToggle(row: any) {
  try {
    await request.put(`/tasks/${row.id}/toggle`)
    ElMessage.success('状态已切换')
    loadData()
  } catch { }
}

onMounted(loadData)
</script>
