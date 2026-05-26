<template>
  <el-card shadow="never">
    <template #header><span>版本历史</span></template>
    <el-table :data="backups" v-loading="loading" stripe @row-click="handleRowClick">
      <el-table-column label="设备" min-width="150">
        <template #default="{ row }">{{ row.device?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="IP" width="140">
        <template #default="{ row }">{{ row.device?.management_ip || '-' }}</template>
      </el-table-column>
      <el-table-column prop="created_at" label="备份时间" width="180" />
      <el-table-column prop="triggered_by" label="触发方式" width="100" />
      <el-table-column label="Git提交" width="110">
        <template #default="{ row }">
          <code style="font-size:12px;">{{ row.git_commit_id?.substring(0, 8) || '-' }}</code>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="handleView(row)">查看</el-button>
          <el-button type="warning" link size="small" @click="handleCompare(row)">对比</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div v-if="total > pageSize" class="pagination-wrap">
      <el-pagination
        v-model:current-page="page"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadData"
      />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'
import { ElMessage } from 'element-plus'

const backups = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

let compareTarget: any = null

async function loadData() {
  loading.value = true
  try {
    const res = await request.get('/configs/history/all', { params: { page: page.value, page_size: pageSize.value } })
    const data = res.data || { list: [], total: 0 }
    backups.value = data.list || []
    total.value = data.total || 0
  } catch { backups.value = [] }
  finally { loading.value = false }
}

function handleView(row: any) {
  ElMessage.info(`查看备份: #${row.id} - ${row.device?.name || '设备'}`)
  // TODO: 打开配置内容查看对话框
}

function handleCompare(row: any) {
  if (!compareTarget) {
    compareTarget = row
    ElMessage.info(`已选择 #${row.id} 作为对比基准，请点击另一个版本的"对比"`)
  } else if (compareTarget.id === row.id) {
    ElMessage.warning('请选择不同的版本进行对比')
  } else {
    const id1 = compareTarget.id
    const id2 = row.id
    ElMessage.info(`对比版本 #${id1} ↔ #${id2}`)
    compareTarget = null
    // TODO: 打开配置对比对话框
  }
}

function handleRowClick(row: any) {
  handleView(row)
}

onMounted(loadData)
</script>

<style scoped>
.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
