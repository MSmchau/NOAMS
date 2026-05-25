<template>
  <div class="page">
    <div class="stats-row">
      <div class="mini-stat" v-for="s in statItems" :key="s.label">
        <span class="mini-stat-label">{{ s.label }}</span>
        <span class="mini-stat-value" :style="{ color: s.color }">{{ s.value }}</span>
      </div>
    </div>

    <div class="table-card">
      <div class="table-header">
        <span class="table-title">告警列表</span>
        <div class="table-filters">
          <el-select v-model="filters.severity" placeholder="级别" clearable size="small" style="width:110px" @change="loadData">
            <el-option label="严重" value="critical" /><el-option label="警告" value="warning" /><el-option label="提示" value="info" />
          </el-select>
          <el-select v-model="filters.status" placeholder="状态" clearable size="small" style="width:110px" @change="loadData">
            <el-option label="已触发" value="triggered" /><el-option label="已确认" value="confirmed" /><el-option label="已解决" value="resolved" />
          </el-select>
        </div>
      </div>
      <table class="tech-table">
        <thead><tr>
          <th>级别</th><th>设备</th><th>类型</th><th>消息</th><th>状态</th><th>触发时间</th><th>操作</th>
        </tr></thead>
        <tbody>
          <tr v-for="row in alerts" :key="row.id">
            <td><span class="sev-tag" :class="'sev-'+row.severity">{{ sevLabel(row.severity) }}</span></td>
            <td><span class="cell-primary">{{ row.device?.name || '-' }}</span></td>
            <td><span class="cell-text">{{ row.alert_type }}</span></td>
            <td><span class="cell-text" style="max-width:240px;display:inline-block;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{{ row.message }}</span></td>
            <td><span class="cell-badge" :class="row.status === 'triggered' ? 'danger' : row.status === 'confirmed' ? 'warn' : 'success'">{{ row.status === 'triggered' ? '待处理' : row.status === 'confirmed' ? '已确认' : '已解决' }}</span></td>
            <td><span class="cell-mono-sm">{{ row.triggered_at?.substring(0,19) || '-' }}</span></td>
            <td>
              <el-button v-if="row.status==='triggered'" type="primary" link size="small" @click="handleConfirm(row)">确认</el-button>
              <el-button v-if="row.status!=='resolved'" type="success" link size="small" @click="handleResolve(row)">解决</el-button>
            </td>
          </tr>
          <tr v-if="!alerts.length"><td colspan="7" class="empty-row">暂无告警</td></tr>
        </tbody>
      </table>
      <div v-if="total > pageSize" class="pagination-wrap">
        <el-pagination v-model:current-page="page" :page-size="pageSize" :total="total" layout="total, prev, pager, next" @current-change="loadData" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { getAlerts, getAlertStats, confirmAlert, resolveAlert } from '@/api/device'
import { ElMessage } from 'element-plus'

const alerts = ref<any[]>([])
const loading = ref(false); const page = ref(1); const pageSize = ref(20); const total = ref(0)
const stats = ref({ critical: 0, warning: 0, info: 0, triggered: 0, resolved: 0 })
const filters = reactive({ severity: '', status: '' })

const statItems = computed(() => [
  { label: '严重', value: stats.value.critical, color: '#ef4444' },
  { label: '警告', value: stats.value.warning, color: '#f59e0b' },
  { label: '提示', value: stats.value.info, color: '#00d4ff' },
  { label: '已解决', value: stats.value.resolved, color: '#10b981' },
])

function sevLabel(s: string) { return { critical: '严重', warning: '警告', info: '提示' }[s] || s }

async function loadData() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: page.value, page_size: pageSize.value }
    if (filters.severity) params.severity = filters.severity
    if (filters.status) params.status = filters.status
    const res = await getAlerts(params); alerts.value = res.data?.list || []; total.value = res.data?.total || 0
  } catch { alerts.value = [] }
  finally { loading.value = false }
}
async function loadStats() { try { const res = await getAlertStats(); stats.value = res.data } catch { } }
async function handleConfirm(row: any) { try { await confirmAlert(row.id); ElMessage.success('已确认'); loadData() } catch { } }
async function handleResolve(row: any) { try { await resolveAlert(row.id); ElMessage.success('已解决'); loadData() } catch { } }

onMounted(() => { loadData(); loadStats() })
</script>

<style scoped>
.page { display: flex; flex-direction: column; gap: 16px; }

.stats-row { display: grid; grid-template-columns: repeat(4,1fr); gap: 12px; }
.mini-stat {
  padding: 16px; background: rgba(26,35,59,0.7); backdrop-filter: blur(16px);
  border: 1px solid rgba(0,212,255,0.06); border-radius: 10px; text-align: center;
}
.mini-stat-label { display: block; font-size: 11px; color: #64748b; text-transform: uppercase; letter-spacing: 0.5px; margin-bottom: 6px; }
.mini-stat-value { font-size: 24px; font-weight: 700; }

.table-card {
  background: rgba(26,35,59,0.7); backdrop-filter: blur(16px);
  border: 1px solid rgba(0,212,255,0.06); border-radius: 12px; padding: 20px;
}
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.table-title { font-size: 15px; font-weight: 600; color: #e2e8f0; }
.table-filters { display: flex; gap: 8px; }

.tech-table { width: 100%; border-collapse: collapse; }
.tech-table th {
  font-size: 11px; text-transform: uppercase; letter-spacing: 0.8px;
  color: #64748b; font-weight: 600; padding: 10px 12px; text-align: left;
  border-bottom: 1px solid rgba(0,212,255,0.06);
}
.tech-table td { padding: 10px 12px; font-size: 13px; color: #cbd5e1; border-bottom: 1px solid rgba(0,212,255,0.04); }
.tech-table tbody tr:hover { background: rgba(0,212,255,0.03); }

.cell-primary { color: #e2e8f0; font-weight: 500; }
.cell-text { color: #94a3b8; }
.cell-mono-sm { font-family: 'JetBrains Mono', monospace; font-size: 11px; color: #475569; }

.sev-tag { font-size: 11px; padding: 2px 10px; border-radius: 4px; font-weight: 500; }
.sev-critical { background: rgba(239,68,68,0.1); color: #ef4444; }
.sev-warning { background: rgba(245,158,11,0.1); color: #f59e0b; }
.sev-info { background: rgba(0,212,255,0.08); color: #00d4ff; }

.cell-badge { font-size: 11px; padding: 2px 10px; border-radius: 4px; font-weight: 500; }
.cell-badge.danger { background: rgba(239,68,68,0.1); color: #ef4444; }
.cell-badge.warn { background: rgba(245,158,11,0.1); color: #f59e0b; }
.cell-badge.success { background: rgba(16,185,129,0.1); color: #10b981; }

.empty-row { text-align: center; color: #475569 !important; padding: 40px !important; }
.pagination-wrap { display: flex; justify-content: flex-end; margin-top: 16px; }
</style>
