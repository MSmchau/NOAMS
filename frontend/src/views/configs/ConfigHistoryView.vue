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
      <el-table-column label="配置Hash" width="130">
        <template #default="{ row }">
          <code style="font-size:12px;">{{ row.config_hash?.substring(0, 12) || '-' }}</code>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="210" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="handleView(row)">查看</el-button>
          <el-button type="warning" link size="small" @click="handleCompare(row)">对比</el-button>
          <el-button type="success" link size="small" @click="handleDownload(row)">下载</el-button>
          <el-button type="danger" link size="small" @click="handleRollback(row)">回滚</el-button>
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

    <!-- 查看配置对话框 -->
    <el-dialog v-model="viewDialogVisible" :title="viewTitle" width="700px" :close-on-click-modal="false">
      <pre class="config-content">{{ viewContent || '暂无配置内容' }}</pre>
      <template #footer>
        <el-button @click="viewDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 对比对话框 -->
    <el-dialog v-model="diffDialogVisible" title="配置对比" width="90%" :close-on-click-modal="false" top="5vh">
      <el-row :gutter="12">
        <el-col :span="12">
          <h4 style="margin:0 0 8px;font-size:13px;color:#595959;">{{ diffLeftTitle }}</h4>
          <pre class="config-content">{{ diffLeftContent || '暂无内容' }}</pre>
        </el-col>
        <el-col :span="12">
          <h4 style="margin:0 0 8px;font-size:13px;color:#595959;">{{ diffRightTitle }}</h4>
          <pre class="config-content">{{ diffRightContent || '暂无内容' }}</pre>
        </el-col>
      </el-row>
      <template #footer>
        <el-button @click="diffDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/api/request'
import { rollbackConfig } from '@/api/config'
import { ElMessage, ElMessageBox } from 'element-plus'

const backups = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

let compareTarget: any = null

// 查看对话框
const viewDialogVisible = ref(false)
const viewTitle = ref('')
const viewContent = ref('')

// 对比对话框
const diffDialogVisible = ref(false)
const diffLeftTitle = ref('')
const diffRightTitle = ref('')
const diffLeftContent = ref('')
const diffRightContent = ref('')

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

async function handleView(row: any) {
  viewTitle.value = `配置 - ${row.device?.name || 'Device#' + row.device_id} (#${row.id})`
  viewContent.value = row.content || '暂无配置内容'
  viewDialogVisible.value = true
}

function handleCompare(row: any) {
  if (!compareTarget) {
    compareTarget = row
    ElMessage.info(`已选择 #${row.id} 作为对比基准，请点击另一个版本的"对比"`)
  } else if (compareTarget.id === row.id) {
    ElMessage.warning('请选择不同的版本进行对比')
  } else {
    const b1 = compareTarget
    const b2 = row
    diffLeftTitle.value = `${b1.device?.name || ''} #${b1.id} - ${b1.created_at || ''}`
    diffRightTitle.value = `${b2.device?.name || ''} #${b2.id} - ${b2.created_at || ''}`
    diffLeftContent.value = b1.content || '暂无内容'
    diffRightContent.value = b2.content || '暂无内容'
    diffDialogVisible.value = true
    compareTarget = null
  }
}

async function handleDownload(row: any) {
  try {
    const res = await request.get(`/configs/export/${row.id}`, { responseType: 'blob' })
    const deviceName = row.device?.name || `device_${row.device_id}`
    const filename = `config_${deviceName}_${(row.created_at || '').substring(0, 10)}.txt`
    const url = window.URL.createObjectURL(new Blob([res as any]))
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    a.click()
    window.URL.revokeObjectURL(url)
  } catch {
    ElMessage.error('下载失败')
  }
}

async function handleRollback(row: any) {
  try {
    await ElMessageBox.confirm(
      `确认将设备「${row.device?.name || 'Device#' + row.device_id}」回滚到版本 #${row.id}？`,
      '回滚确认',
      { type: 'warning', confirmButtonText: '确认回滚', cancelButtonText: '取消' }
    )
    await rollbackConfig(row.id)
    ElMessage.success('回滚任务已创建')
  } catch { }
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
.config-content {
  background: #f5f5f5;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  padding: 12px;
  font-size: 12px;
  line-height: 1.6;
  max-height: 500px;
  overflow: auto;
  white-space: pre;
  font-family: 'JetBrains Mono', 'Courier New', monospace;
  margin: 0;
}
</style>
