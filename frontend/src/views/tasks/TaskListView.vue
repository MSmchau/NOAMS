<template>
  <el-card shadow="hover">
    <template #header>
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <span>定时任务</span>
        <el-button type="primary" size="small" @click="showCreateDialog">
          <el-icon><Plus /></el-icon> 新建任务
        </el-button>
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
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-switch :model-value="!!row.status" @change="handleToggle(row)" size="small" style="margin-right:8px" />
          <el-button type="primary" link size="small" @click="showEditDialog(row)">编辑</el-button>
          <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div v-if="total > pageSize" style="display:flex;justify-content:flex-end;margin-top:16px;">
      <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadData" />
    </div>

    <!-- 新建/编辑任务对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEditing ? '编辑任务' : '新建任务'" width="500px" :close-on-click-modal="false">
      <el-form :model="form" label-width="100px" size="small">
        <el-form-item label="任务名称" required>
          <el-input v-model="form.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="任务类型" required>
          <el-select v-model="form.task_type" style="width:100%">
            <el-option label="设备巡检" value="inspection" />
            <el-option label="配置备份" value="backup" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="Cron表达式" required>
          <el-input v-model="form.cron_expr" placeholder="如: 0 6 * * * (每天6点)" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="任务描述（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getTasks, createTask, updateTask, deleteTask, toggleTask } from '@/api/task'
import { ElMessage, ElMessageBox } from 'element-plus'

const tasks = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

// Dialog
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref<number | null>(null)
const saving = ref(false)
const form = ref({
  name: '',
  task_type: 'inspection',
  cron_expr: '',
  description: '',
})

async function loadData() {
  loading.value = true
  try {
    const res = await getTasks({ page: page.value, page_size: pageSize.value })
    tasks.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch { tasks.value = [] }
  finally { loading.value = false }
}

function resetForm() {
  form.value = { name: '', task_type: 'inspection', cron_expr: '', description: '' }
}

function showCreateDialog() {
  isEditing.value = false
  editingId.value = null
  resetForm()
  dialogVisible.value = true
}

function showEditDialog(row: any) {
  isEditing.value = true
  editingId.value = row.id
  form.value = {
    name: row.name || '',
    task_type: row.task_type || 'inspection',
    cron_expr: row.cron_expr || '',
    description: row.description || '',
  }
  dialogVisible.value = true
}

async function handleSave() {
  if (!form.value.name || !form.value.task_type || !form.value.cron_expr) {
    ElMessage.warning('请填写必填项')
    return
  }
  saving.value = true
  try {
    if (isEditing.value && editingId.value) {
      await updateTask(editingId.value, form.value)
      ElMessage.success('任务已更新')
    } else {
      await createTask(form.value)
      ElMessage.success('任务已创建')
    }
    dialogVisible.value = false
    loadData()
  } catch {
    ElMessage.error(isEditing.value ? '更新失败' : '创建失败')
  } finally {
    saving.value = false
  }
}

async function handleToggle(row: any) {
  try {
    await toggleTask(row.id)
    ElMessage.success('状态已切换')
    loadData()
  } catch { }
}

async function handleDelete(row: any) {
  try {
    await ElMessageBox.confirm(`确认删除任务「${row.name}」？`, '提示', {
      type: 'warning',
      confirmButtonText: '确认',
      cancelButtonText: '取消',
    })
    await deleteTask(row.id)
    ElMessage.success('任务已删除')
    loadData()
  } catch { }
}

onMounted(loadData)
</script>
