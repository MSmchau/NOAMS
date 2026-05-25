<template>
  <div class="dashboard">
    <!-- Stats cards -->
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-info">
              <div class="stat-label">设备总数</div>
              <div class="stat-value">{{ stats.device_stats?.total || 0 }}</div>
            </div>
            <el-icon class="stat-icon" :size="48" color="#409eff"><Server /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-info">
              <div class="stat-label">在线设备</div>
              <div class="stat-value" style="color: #67c23a">{{ stats.device_stats?.online || 0 }}</div>
            </div>
            <el-icon class="stat-icon" :size="48" color="#67c23a"><Connection /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-info">
              <div class="stat-label">离线设备</div>
              <div class="stat-value" style="color: #f56c6c">{{ stats.device_stats?.offline || 0 }}</div>
            </div>
            <el-icon class="stat-icon" :size="48" color="#f56c6c"><WarningFilled /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <div class="stat-card">
            <div class="stat-info">
              <div class="stat-label">未处理告警</div>
              <div class="stat-value" style="color: #e6a23c">{{ stats.alert_stats?.triggered || 0 }}</div>
            </div>
            <el-icon class="stat-icon" :size="48" color="#e6a23c"><BellFilled /></el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- CPU & Memory Top -->
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span>CPU 使用率 TOP 10</span>
          </template>
          <div v-if="stats.cpu_top?.length">
            <div v-for="item in stats.cpu_top" :key="item.device_id" class="resource-bar">
              <div class="resource-name">{{ item.device_name || 'Device#' + item.device_id }}</div>
              <el-progress
                :percentage="Math.round(item.cpu_usage)"
                :color="item.cpu_usage > 80 ? '#f56c6c' : item.cpu_usage > 50 ? '#e6a23c' : '#67c23a'"
              />
            </div>
          </div>
          <el-empty v-else description="暂无数据" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span>内存使用率 TOP 10</span>
          </template>
          <div v-if="stats.mem_top?.length">
            <div v-for="item in stats.mem_top" :key="item.device_id" class="resource-bar">
              <div class="resource-name">{{ item.device_name || 'Device#' + item.device_id }}</div>
              <el-progress
                :percentage="Math.round(item.mem_usage)"
                :color="item.mem_usage > 80 ? '#f56c6c' : item.mem_usage > 50 ? '#e6a23c' : '#67c23a'"
              />
            </div>
          </div>
          <el-empty v-else description="暂无数据" />
        </el-card>
      </el-col>
    </el-row>

    <!-- Recent inspections -->
    <el-card shadow="hover">
      <template #header>
        <span>最近巡检记录</span>
      </template>
      <el-table :data="stats.recent_checks || []" stripe style="width: 100%">
        <el-table-column prop="device?.name" label="设备名称" min-width="150" />
        <el-table-column prop="device?.management_ip" label="管理IP" width="140" />
        <el-table-column prop="cpu_usage" label="CPU" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.cpu_usage" :type="row.cpu_usage > 80 ? 'danger' : 'success'" size="small">
              {{ row.cpu_usage }}%
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="memory_usage" label="内存" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.memory_usage" :type="row.memory_usage > 80 ? 'danger' : 'success'" size="small">
              {{ row.memory_usage }}%
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.status === 'success'" type="success" size="small">正常</el-tag>
            <el-tag v-else-if="row.status === 'failed'" type="danger" size="small">失败</el-tag>
            <el-tag v-else type="warning" size="small">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="inspected_at" label="巡检时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getDashboard } from '@/api/device'

const stats = ref<any>({})

async function loadDashboard() {
  try {
    const res = await getDashboard()
    stats.value = res.data || {}
  } catch {
    // silently fail
  }
}

onMounted(() => {
  loadDashboard()
})
</script>

<style scoped>
.stat-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.stat-info {
  flex: 1;
}
.stat-label {
  font-size: 14px;
  color: #8a8f9d;
  margin-bottom: 8px;
}
.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
}
.stat-icon {
  margin-left: 16px;
}
.resource-bar {
  margin-bottom: 16px;
}
.resource-name {
  font-size: 13px;
  color: #606266;
  margin-bottom: 4px;
}
</style>
