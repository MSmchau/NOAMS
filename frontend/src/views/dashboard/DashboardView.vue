<template>
  <div class="dashboard-page">
    <!-- Stats Cards -->
    <el-row :gutter="16">
      <el-col :xs="12" :sm="12" :md="6" v-for="s in statCards" :key="s.label">
        <el-card shadow="never" class="stat-card">
          <div class="stat-card-body">
            <div class="stat-icon-wrap" :style="{ background: s.bg }">
              <el-icon :size="22" :style="{ color: s.color }">
                <component :is="s.icon" />
              </el-icon>
            </div>
            <div class="stat-info">
              <span class="stat-label">{{ s.label }}</span>
              <span class="stat-value" :style="{ color: s.color }">
                {{ s.value }}<span class="stat-unit">{{ s.unit }}</span>
              </span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Charts Row -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="12">
        <el-card shadow="never">
          <template #header>
            <div class="chart-header">
              <span>CPU 使用率 TOP10</span>
              <el-tag size="small" type="info">实时</el-tag>
            </div>
          </template>
          <div class="chart-container">
            <v-chart v-if="cpuChartOption" :option="cpuChartOption" autoresize />
            <el-empty v-else description="暂无数据" :image-size="60" />
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card shadow="never">
          <template #header>
            <div class="chart-header">
              <span>内存使用率 TOP10</span>
              <el-tag size="small" type="info">实时</el-tag>
            </div>
          </template>
          <div class="chart-container">
            <v-chart v-if="memChartOption" :option="memChartOption" autoresize />
            <el-empty v-else description="暂无数据" :image-size="60" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Recent Inspections -->
    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="24">
        <el-card shadow="never">
          <template #header>
            <div class="table-header">
              <span>最近巡检记录</span>
              <span class="table-count">共 {{ recentChecks.length }} 条</span>
            </div>
          </template>
          <el-table :data="recentChecks" stripe size="small" max-height="420" v-loading="loading">
            <el-table-column prop="device?.name" label="设备名称" min-width="150" />
            <el-table-column label="管理IP" width="140">
              <template #default="{ row }">
                <code style="font-size:12px;color:#595959">{{ row.device?.management_ip || '-' }}</code>
              </template>
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
            <el-table-column label="状态" width="90" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.status === 'success'" type="success" size="small">正常</el-tag>
                <el-tag v-else-if="row.status === 'failed'" type="danger" size="small">失败</el-tag>
                <el-tag v-else type="warning" size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="巡检时间" width="170">
              <template #default="{ row }">
                <span style="font-size:12px;color:#8c8c8c">{{ row.inspected_at?.substring(0, 19) || '-' }}</span>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getDashboard } from '@/api/device'
import { use } from 'echarts/core'
import { BarChart } from 'echarts/charts'
import { GridComponent, TooltipComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'

use([BarChart, GridComponent, TooltipComponent, CanvasRenderer])

const loading = ref(false)
const stats = ref<any>({})

const statCards = computed(() => [
  {
    label: '设备总数', value: stats.value.device_stats?.total ?? 0, unit: '台',
    icon: 'Server', color: '#1890ff', bg: '#e6f7ff',
  },
  {
    label: '在线设备', value: stats.value.device_stats?.online ?? 0, unit: '台',
    icon: 'Connection', color: '#52c41a', bg: '#f6ffed',
  },
  {
    label: '离线设备', value: stats.value.device_stats?.offline ?? 0, unit: '台',
    icon: 'WarningFilled', color: '#ff4d4f', bg: '#fff2e8',
  },
  {
    label: '未处理告警', value: stats.value.alert_stats?.triggered ?? 0, unit: '条',
    icon: 'BellFilled', color: '#fa8c16', bg: '#fff7e6',
  },
])

function cpuTagType(val: number) {
  if (val > 80) return 'danger'
  if (val > 50) return 'warning'
  return 'success'
}

// ECharts options
const recentChecks = computed(() => stats.value.recent_checks || [])

const cpuChartOption = computed(() => {
  const top = stats.value.cpu_top
  if (!top?.length) return null
  return {
    tooltip: { trigger: 'axis' as const, axisPointer: { type: 'shadow' as const } },
    grid: { left: 10, right: 20, top: 10, bottom: 20, containLabel: true },
    xAxis: { type: 'value' as const, max: 100, axisLabel: { formatter: '{value}%' }, splitLine: { lineStyle: { color: '#f0f0f0' } } },
    yAxis: {
      type: 'category' as const,
      data: top.map((i: any) => i.device_name || 'Device#' + i.device_id).reverse(),
      axisLabel: { fontSize: 11, color: '#595959' },
    },
    series: [{
      type: 'bar',
      data: top.map((i: any) => ({
        value: Math.min(i.cpu_usage, 100),
        itemStyle: {
          color: i.cpu_usage > 80 ? '#ff4d4f' : i.cpu_usage > 50 ? '#fa8c16' : '#1890ff',
          borderRadius: [0, 4, 4, 0],
        },
      })).reverse(),
      barWidth: 18,
    }],
  }
})

const memChartOption = computed(() => {
  const top = stats.value.mem_top
  if (!top?.length) return null
  return {
    tooltip: { trigger: 'axis' as const, axisPointer: { type: 'shadow' as const } },
    grid: { left: 10, right: 20, top: 10, bottom: 20, containLabel: true },
    xAxis: { type: 'value' as const, max: 100, axisLabel: { formatter: '{value}%' }, splitLine: { lineStyle: { color: '#f0f0f0' } } },
    yAxis: {
      type: 'category' as const,
      data: top.map((i: any) => i.device_name || 'Device#' + i.device_id).reverse(),
      axisLabel: { fontSize: 11, color: '#595959' },
    },
    series: [{
      type: 'bar',
      data: top.map((i: any) => ({
        value: Math.min(i.mem_usage, 100),
        itemStyle: {
          color: i.mem_usage > 80 ? '#ff4d4f' : i.mem_usage > 50 ? '#fa8c16' : '#52c41a',
          borderRadius: [0, 4, 4, 0],
        },
      })).reverse(),
      barWidth: 18,
    }],
  }
})

async function loadDashboard() {
  try {
    const res = await getDashboard()
    stats.value = res.data || {}
  } catch {
    stats.value = {}
  }
}

onMounted(loadDashboard)
</script>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
}

.stat-card {
  margin-bottom: 0;
}

.stat-card-body {
  display: flex;
  align-items: center;
  gap: 16px;
}

.stat-icon-wrap {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-label {
  font-size: 13px;
  color: #8c8c8c;
  font-weight: 400;
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  line-height: 1.2;
}

.stat-unit {
  font-size: 12px;
  font-weight: 400;
  margin-left: 2px;
  color: #bfbfbf;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  height: 320px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chart-container > * {
  width: 100%;
  height: 100%;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.table-count {
  font-size: 12px;
  color: #8c8c8c;
}
</style>
