<template>
  <div class="dashboard">
    <!-- Stats row -->
    <div class="stats-grid">
      <div class="stat-card" v-for="s in statCards" :key="s.label">
        <div class="stat-icon" :style="{ color: s.color, background: s.bg }">
          <el-icon :size="22"><component :is="s.icon" /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-label">{{ s.label }}</div>
          <div class="stat-value" :style="{ color: s.color }">
            {{ s.value }}
            <span class="stat-unit">{{ s.unit }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Charts row -->
    <div class="charts-row">
      <div class="chart-card">
        <div class="chart-header">
          <span class="chart-title">CPU 使用率 TOP10</span>
          <span class="chart-badge">实时</span>
        </div>
        <div v-if="stats.cpu_top?.length" class="chart-body">
          <div v-for="item in stats.cpu_top" :key="item.device_id" class="bar-row">
            <span class="bar-label">{{ item.device_name || 'Device#' + item.device_id }}</span>
            <div class="bar-track">
              <div class="bar-fill" :style="{
                width: Math.min(item.cpu_usage, 100) + '%',
                background: item.cpu_usage > 80 ? '#ef4444' : item.cpu_usage > 50 ? '#f59e0b' : '#00d4ff'
              }"></div>
            </div>
            <span class="bar-value">{{ item.cpu_usage }}%</span>
          </div>
        </div>
        <el-empty v-else description="暂无数据" :image-size="80" />
      </div>

      <div class="chart-card">
        <div class="chart-header">
          <span class="chart-title">内存使用率 TOP10</span>
          <span class="chart-badge">实时</span>
        </div>
        <div v-if="stats.mem_top?.length" class="chart-body">
          <div v-for="item in stats.mem_top" :key="item.device_id" class="bar-row">
            <span class="bar-label">{{ item.device_name || 'Device#' + item.device_id }}</span>
            <div class="bar-track">
              <div class="bar-fill" :style="{
                width: Math.min(item.mem_usage, 100) + '%',
                background: item.mem_usage > 80 ? '#ef4444' : item.mem_usage > 50 ? '#f59e0b' : '#7c3aed'
              }"></div>
            </div>
            <span class="bar-value">{{ item.mem_usage }}%</span>
          </div>
        </div>
        <el-empty v-else description="暂无数据" :image-size="80" />
      </div>
    </div>

    <!-- Recent inspections -->
    <div class="table-card">
      <div class="table-header">
        <span class="table-title">最近巡检记录</span>
        <span class="table-count">共 {{ stats.recent_checks?.length || 0 }} 条</span>
      </div>
      <div class="table-scroll">
        <table class="tech-table">
          <thead>
            <tr>
              <th>设备名称</th>
              <th>管理IP</th>
              <th>CPU</th>
              <th>内存</th>
              <th>状态</th>
              <th>时间</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in stats.recent_checks || []" :key="row.id">
              <td><span class="cell-primary">{{ row.device?.name || '-' }}</span></td>
              <td><code class="cell-mono">{{ row.device?.management_ip || '-' }}</code></td>
              <td>
                <span v-if="row.cpu_usage != null" class="cell-tag" :class="{ danger: row.cpu_usage > 80, warn: row.cpu_usage > 50 && row.cpu_usage <= 80 }">
                  {{ row.cpu_usage }}%
                </span>
                <span v-else class="cell-muted">-</span>
              </td>
              <td>
                <span v-if="row.memory_usage != null" class="cell-tag" :class="{ danger: row.memory_usage > 80, warn: row.memory_usage > 50 && row.memory_usage <= 80 }">
                  {{ row.memory_usage }}%
                </span>
                <span v-else class="cell-muted">-</span>
              </td>
              <td>
                <span v-if="row.status === 'success'" class="cell-badge success">正常</span>
                <span v-else-if="row.status === 'failed'" class="cell-badge danger">失败</span>
                <span v-else class="cell-badge warn">{{ row.status }}</span>
              </td>
              <td><span class="cell-mono-sm">{{ row.inspected_at ? row.inspected_at.substring(0, 19) : '-' }}</span></td>
            </tr>
            <tr v-if="!(stats.recent_checks?.length)">
              <td colspan="6" class="empty-row">暂无巡检记录</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getDashboard } from '@/api/device'

const stats = ref<any>({})

const statCards = computed(() => [
  {
    label: '设备总数', value: stats.value.device_stats?.total ?? 0, unit: '台',
    icon: 'Server', color: '#00d4ff', bg: 'rgba(0,212,255,0.08)',
  },
  {
    label: '在线设备', value: stats.value.device_stats?.online ?? 0, unit: '台',
    icon: 'Connection', color: '#10b981', bg: 'rgba(16,185,129,0.08)',
  },
  {
    label: '离线设备', value: stats.value.device_stats?.offline ?? 0, unit: '台',
    icon: 'WarningFilled', color: '#ef4444', bg: 'rgba(239,68,68,0.08)',
  },
  {
    label: '未处理告警', value: stats.value.alert_stats?.triggered ?? 0, unit: '条',
    icon: 'BellFilled', color: '#f59e0b', bg: 'rgba(245,158,11,0.08)',
  },
])

async function loadDashboard() {
  try { const res = await getDashboard(); stats.value = res.data || {} } catch { }
}

onMounted(loadDashboard)
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Stats grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: rgba(26, 35, 59, 0.85);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(0, 212, 255, 0.06);
  border-radius: 12px;
  transition: all 0.3s ease;
}
.stat-card:hover {
  transform: translateY(-2px);
  border-color: rgba(0, 212, 255, 0.15);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-info { flex: 1; }
.stat-label {
  font-size: 12px;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 4px;
}
.stat-value {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
}
.stat-unit {
  font-size: 13px;
  color: #475569;
  font-weight: 400;
  margin-left: 4px;
}

/* Charts row */
.charts-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.chart-card {
  background: rgba(26, 35, 59, 0.7);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(0, 212, 255, 0.06);
  border-radius: 12px;
  padding: 20px;
  transition: all 0.3s ease;
}
.chart-card:hover {
  border-color: rgba(0, 212, 255, 0.12);
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.chart-title {
  font-size: 14px;
  font-weight: 600;
  color: #e2e8f0;
}
.chart-badge {
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(0, 212, 255, 0.08);
  color: #00d4ff;
  letter-spacing: 1px;
}

.bar-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}
.bar-label {
  width: 120px;
  font-size: 12px;
  color: #94a3b8;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex-shrink: 0;
}
.bar-track {
  flex: 1;
  height: 6px;
  background: rgba(255, 255, 255, 0.04);
  border-radius: 3px;
  overflow: hidden;
}
.bar-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.6s ease;
}
.bar-value {
  width: 44px;
  text-align: right;
  font-size: 12px;
  color: #94a3b8;
  font-family: 'JetBrains Mono', monospace;
}

/* Tech table */
.table-card {
  background: rgba(26, 35, 59, 0.7);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(0, 212, 255, 0.06);
  border-radius: 12px;
  padding: 20px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}
.table-title {
  font-size: 14px;
  font-weight: 600;
  color: #e2e8f0;
}
.table-count {
  font-size: 12px;
  color: #64748b;
}

.tech-table {
  width: 100%;
  border-collapse: collapse;
}

.tech-table th {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  color: #64748b;
  font-weight: 600;
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid rgba(0, 212, 255, 0.06);
}

.tech-table td {
  padding: 10px 12px;
  font-size: 13px;
  color: #cbd5e1;
  border-bottom: 1px solid rgba(0, 212, 255, 0.04);
}

.tech-table tbody tr:hover {
  background: rgba(0, 212, 255, 0.03);
}

.cell-primary { color: #e2e8f0; font-weight: 500; }
.cell-mono { font-family: 'JetBrains Mono', monospace; font-size: 12px; color: #64748b; }
.cell-mono-sm { font-family: 'JetBrains Mono', monospace; font-size: 11px; color: #475569; }
.cell-muted { color: #475569; }

.cell-tag {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(0, 212, 255, 0.06);
  color: #00d4ff;
}
.cell-tag.warn { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }
.cell-tag.danger { background: rgba(239, 68, 68, 0.1); color: #ef4444; }

.cell-badge {
  font-size: 11px;
  padding: 2px 10px;
  border-radius: 4px;
  font-weight: 500;
}
.cell-badge.success { background: rgba(16, 185, 129, 0.1); color: #10b981; }
.cell-badge.danger { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
.cell-badge.warn { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }

.empty-row {
  text-align: center;
  color: #475569 !important;
  padding: 40px !important;
  font-size: 13px;
}
</style>
