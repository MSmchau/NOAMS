<template>
  <div class="page-container">
    <!-- 操作栏 -->
    <el-card shadow="never">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <span style="font-size:15px;font-weight:600;">巡检记录</span>
        <div style="display: flex; gap: 8px;">
          <el-button type="primary" size="small" @click="showManualDialog = true">
            <el-icon><Monitor /></el-icon> 手动巡检
          </el-button>
          <el-button size="small" @click="handleSchedule">
            <el-icon><Clock /></el-icon> 定时巡检
          </el-button>
          <el-button size="small" @click="handleRefresh">
            <el-icon><Refresh /></el-icon> 刷新
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 巡检记录表 -->
    <el-card shadow="never">
      <el-table :data="inspections" v-loading="loading" stripe size="small">
        <el-table-column label="设备" min-width="150">
          <template #default="{ row }">{{ row.device?.name || '-' }}</template>
        </el-table-column>
        <el-table-column label="CPU" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.cpu_usage != null" :type="cpuTagType(row.cpu_usage)" size="small">{{ row.cpu_usage }}%</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column label="内存" width="90" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.memory_usage != null" :type="cpuTagType(row.memory_usage)" size="small">{{ row.memory_usage }}%</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="uptime" label="运行时长" min-width="140" />
        <el-table-column label="异常" width="70" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.is_anomaly" type="danger" size="small">是</el-tag>
            <el-tag v-else type="success" size="small">否</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="inspected_at" label="巡检时间" width="175" />
      </el-table>

      <div v-if="total > pageSize" style="display:flex;justify-content:flex-end;margin-top:16px;">
        <el-pagination
          v-model:current-page="page" :page-size="pageSize" :total="total"
          layout="total, prev, pager, next"
          @current-change="loadData"
        />
      </div>
    </el-card>

    <!-- 手动巡检对话框 -->
    <el-dialog v-model="showManualDialog" title="手动巡检" width="520px" :close-on-click-modal="false">
      <div style="margin-bottom:12px;font-size:13px;color:#595959;">选择要执行巡检的设备：</div>
      <el-table
        :data="allDevices"
        ref="deviceTableRef"
        @selection-change="onSelectionChange"
        stripe
        size="small"
        max-height="360"
      >
        <el-table-column type="selection" width="44" />
        <el-table-column prop="name" label="设备名称" min-width="140" />
        <el-table-column prop="management_ip" label="管理IP" width="140" />
        <el-table-column label="角色" width="80">
          <template #default="{ row }">
            <el-tag :type="roleTagType(row.role)" size="small">{{ roleLabel(row.role) }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div v-if="selectedIds.length" style="margin-top:8px;font-size:12px;color:#8c8c8c;">
        已选 {{ selectedIds.length }} 台设备
      </div>
      <template #footer>
        <el-button @click="showManualDialog = false">取消</el-button>
        <el-button type="primary" :loading="inspecting" :disabled="!selectedIds.length" @click="handleManualInspect">
          开始巡检 ({{ selectedIds.length }})
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getInspectionReport, getAllDevices, batchInspect } from '@/api/device'
import { ElMessage } from 'element-plus'

const router = useRouter()
const inspections = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 手动巡检
const showManualDialog = ref(false)
const allDevices = ref<any[]>([])
const selectedIds = ref<number[]>([])
const inspecting = ref(false)

function cpuTagType(val: number) {
  if (val > 80) return 'danger'
  if (val > 50) return 'warning'
  return 'success'
}

function roleLabel(role?: string) {
  const map: Record<string, string> = { core: '核心', ac: 'AC', aggregation: '本体', access: '接入', ap: 'AP' }
  return map[role || ''] || role || '-'
}

function roleTagType(role?: string) {
  const map: Record<string, 'danger' | 'warning' | 'info' | 'success'> = { core: 'danger', ac: 'warning', aggregation: 'info', access: 'info', ap: 'success' }
  return map[role || ''] || 'info'
}

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

function handleSchedule() {
  router.push('/tasks')
}

async function loadDevices() {
  try {
    const res = await getAllDevices()
    allDevices.value = res.data || []
  } catch { allDevices.value = [] }
}

function onSelectionChange(rows: any[]) {
  selectedIds.value = rows.map((r: any) => r.id)
}

async function handleManualInspect() {
  if (!selectedIds.value.length) return
  inspecting.value = true
  try {
    await batchInspect(selectedIds.value)
    ElMessage.success(`已为 ${selectedIds.value.length} 台设备创建巡检任务`)
    showManualDialog.value = false
    // 延迟后刷新巡检记录
    setTimeout(() => loadData(), 2000)
  } catch {
    ElMessage.error('创建巡检任务失败')
  } finally {
    inspecting.value = false
  }
}

onMounted(() => { loadData() })
</script>
