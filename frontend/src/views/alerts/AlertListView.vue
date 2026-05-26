<template>
  <div class="page-container">
    <!-- Stats Row -->
    <el-row :gutter="12">
      <el-col :span="6" v-for="s in statItems" :key="s.label">
        <el-card shadow="never">
          <div class="mini-stat">
            <span class="mini-stat-label">{{ s.label }}</span>
            <span class="mini-stat-value" :style="{ color: s.color }">{{ s.value }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Alerts Table -->
    <el-card shadow="never">
      <template #header>
        <div class="table-toolbar">
          <span class="toolbar-title">告警列表</span>
          <div class="toolbar-filters">
            <el-select v-model="filters.severity" placeholder="级别" clearable size="small" style="width:110px" @change="loadData">
              <el-option label="严重" value="critical" />
              <el-option label="警告" value="warning" />
              <el-option label="提示" value="info" />
            </el-select>
            <el-select v-model="filters.status" placeholder="状态" clearable size="small" style="width:110px" @change="loadData">
              <el-option label="已触发" value="triggered" />
              <el-option label="已确认" value="confirmed" />
              <el-option label="已解决" value="resolved" />
            </el-select>
          </div>
        </div>
      </template>

      <el-table :data="alerts" v-loading="loading" stripe size="small">
        <el-table-column label="级别" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="sevTagType(row.severity)" size="small">{{ sevLabel(row.severity) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="device?.name" label="设备" min-width="150" />
        <el-table-column prop="alert_type" label="类型" width="120" />
        <el-table-column prop="message" label="消息" min-width="240" show-overflow-tooltip />
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">
              {{ row.status === 'triggered' ? '待处理' : row.status === 'confirmed' ? '已确认' : '已解决' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="触发时间" width="170">
          <template #default="{ row }">
            <span style="font-size:12px;color:#8c8c8c">{{ row.triggered_at?.substring(0,19) || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button v-if="row.status==='triggered'" type="primary" link size="small" @click="handleConfirm(row)">确认</el-button>
            <el-button v-if="row.status!=='resolved'" type="success" link size="small" @click="handleResolve(row)">解决</el-button>
            <span v-if="row.status==='resolved'" style="color:#8c8c8c;font-size:12px">已关闭</span>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { getAlerts, getAlertStats, confirmAlert, resolveAlert } from '@/api/device'
import { ElMessage } from 'element-plus'

const alerts = ref<any[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const stats = ref({ critical: 0, warning: 0, info: 0, triggered: 0, resolved: 0 })
const filters = reactive({ severity: '', status: '' })

const statItems = computed(() => [
  { label: '严重', value: stats.value.critical, color: '#ff4d4f' },
  { label: '警告', value: stats.value.warning, color: '#fa8c16' },
  { label: '提示', value: stats.value.info, color: '#1890ff' },
  { label: '已解决', value: stats.value.resolved, color: '#52c41a' },
])

function sevLabel(s: string) { return { critical: '严重', warning: '警告', info: '提示' }[s] || s }
function sevTagType(s: string): 'danger' | 'warning' | 'info' { return ({ critical: 'danger', warning: 'warning', info: 'info' } as Record<string, 'danger' | 'warning' | 'info'>)[s] || 'info' }
function statusTagType(s: string): 'danger' | 'warning' | 'success' | 'info' { return ({ triggered: 'danger', confirmed: 'warning', resolved: 'success' } as Record<string, 'danger' | 'warning' | 'success' | 'info'>)[s] || 'info' }

async function loadData() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: page.value, page_size: pageSize.value }
    if (filters.severity) params.severity = filters.severity
    if (filters.status) params.status = filters.status
    const res = await getAlerts(params)
    alerts.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch { alerts.value = [] }
  finally { loading.value = false }
}

async function loadStats() {
  try { const res = await getAlertStats(); stats.value = res.data } catch { }
}

async function handleConfirm(row: any) {
  try { await confirmAlert(row.id); ElMessage.success('已确认'); loadData() } catch { }
}

async function handleResolve(row: any) {
  try { await resolveAlert(row.id); ElMessage.success('已解决'); loadData() } catch { }
}

onMounted(() => { loadData(); loadStats() })
</script>

<style scoped>
.mini-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 4px 0;
}

.mini-stat-label {
  font-size: 12px;
  color: #8c8c8c;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 6px;
}

.mini-stat-value {
  font-size: 24px;
  font-weight: 700;
}

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.toolbar-title {
  font-size: 15px;
  font-weight: 600;
  color: #262626;
}

.toolbar-filters {
  display: flex;
  gap: 8px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
