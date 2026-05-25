<template>
  <div class="device-list">
    <!-- Search bar -->
    <div class="search-card">
      <el-form :model="filters" inline size="default">
        <el-form-item label="设备名称">
          <el-input v-model="filters.name" placeholder="搜索设备名称" clearable @clear="search" @keyup.enter="search" />
        </el-form-item>
        <el-form-item label="管理IP">
          <el-input v-model="filters.management_ip" placeholder="搜索IP" clearable @keyup.enter="search" style="width:140px" />
        </el-form-item>
        <el-form-item label="厂商">
          <el-select v-model="filters.vendor" placeholder="厂商" clearable @change="search" style="width:100px">
            <el-option label="H3C" value="h3c" />
            <el-option label="华为" value="huawei" />
            <el-option label="思科" value="cisco" />
            <el-option label="锐捷" value="ruijie" />
          </el-select>
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="filters.role" placeholder="角色" clearable @change="search" style="width:100px">
            <el-option label="核心" value="core" />
            <el-option label="AC" value="ac" />
            <el-option label="本体" value="aggregation" />
            <el-option label="接入" value="access" />
            <el-option label="AP" value="ap" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="状态" clearable @change="search" style="width:90px">
            <el-option label="在线" :value="1" />
            <el-option label="离线" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">查 询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- Device table -->
    <div class="table-card">
      <div class="table-header">
        <div class="table-header-left">
          <span class="table-title">设备列表</span>
          <span class="table-total">共 {{ total }} 台</span>
        </div>
        <div class="table-actions">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>添加设备
          </el-button>
          <el-button @click="loadDevices">
            <el-icon><Refresh /></el-icon>刷新
          </el-button>
        </div>
      </div>

      <div class="table-scroll">
        <table class="tech-table">
          <thead>
            <tr>
              <th>设备名称</th>
              <th>管理IP</th>
              <th>厂商</th>
              <th>型号</th>
              <th>角色</th>
              <th>楼栋</th>
              <th>分组</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="row in devices" :key="row.id" @click="handleDetail(row)" style="cursor:pointer">
              <td><span class="cell-primary">{{ row.name }}</span></td>
              <td><code class="cell-mono">{{ row.management_ip }}</code></td>
              <td><span class="cell-vendor">{{ (row.vendor || '').toUpperCase() || '-' }}</span></td>
              <td><span class="cell-model">{{ row.model || '-' }}</span></td>
              <td>
                <span class="role-badge" :class="'role-' + (row.role || 'access')">
                  {{ roleLabel(row.role) }}
                </span>
              </td>
              <td><span class="cell-text">{{ row.building || '-' }}</span></td>
              <td><span class="cell-text">{{ row.group?.name || '-' }}</span></td>
              <td>
                <span class="status-dot" :class="row.status ? 'online' : 'offline'"></span>
                <span class="status-text" :class="row.status ? 'online' : 'offline'">{{ row.status ? '在线' : '离线' }}</span>
              </td>
              <td @click.stop>
                <el-button type="primary" link size="small" @click="handleDetail(row)">详情</el-button>
                <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
                <el-popconfirm title="确认删除该设备？" @confirm="handleDelete(row)">
                  <template #reference>
                    <el-button type="danger" link size="small">删除</el-button>
                  </template>
                </el-popconfirm>
              </td>
            </tr>
            <tr v-if="!devices.length && !loading">
              <td colspan="9" class="empty-row">暂无设备数据</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="total > pageSize" class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadDevices"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getDevices, deleteDevice } from '@/api/device'
import { ElMessage } from 'element-plus'
import type { Device } from '@/api/device'

const router = useRouter()
const devices = ref<Device[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

const filters = reactive({ name: '', management_ip: '', vendor: '', role: '', status: '' as string | number })

function roleLabel(role?: string) {
  const map: Record<string, string> = { core: '核心', ac: 'AC', aggregation: '本体', access: '接入', ap: 'AP' }
  return map[role || ''] || role || '-'
}

async function loadDevices() {
  loading.value = true
  try {
    const params: Record<string, any> = { page: page.value, page_size: pageSize.value }
    if (filters.name) params.name = filters.name
    if (filters.management_ip) params.management_ip = filters.management_ip
    if (filters.vendor) params.vendor = filters.vendor
    if (filters.role) params.role = filters.role
    if (filters.status !== '') params.status = filters.status
    const res = await getDevices(params)
    devices.value = res.data.list || []
    total.value = res.data.total || 0
  } catch { devices.value = [] }
  finally { loading.value = false }
}

function search() { page.value = 1; loadDevices() }
function resetSearch() { Object.assign(filters, { name: '', management_ip: '', vendor: '', role: '', status: '' }); page.value = 1; loadDevices() }
function handleAdd() { router.push('/devices/create') }
function handleDetail(row: Device) { router.push(`/devices/${row.id}`) }
function handleEdit(row: Device) { router.push(`/devices/${row.id}/edit`) }
async function handleDelete(row: Device) {
  try { await deleteDevice(row.id); ElMessage.success('已删除'); loadDevices() } catch { }
}

onMounted(loadDevices)
</script>

<style scoped>
.device-list { display: flex; flex-direction: column; gap: 16px; }

.search-card {
  padding: 16px 20px;
  background: rgba(26, 35, 59, 0.7);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(0, 212, 255, 0.06);
  border-radius: 12px;
}

.table-card {
  background: rgba(26, 35, 59, 0.7);
  backdrop-filter: blur(16px);
  border: 1px solid rgba(0, 212, 255, 0.06);
  border-radius: 12px;
  padding: 20px;
}

.table-header {
  display: flex; justify-content: space-between; align-items: center;
  margin-bottom: 16px;
}
.table-header-left { display: flex; align-items: center; gap: 12px; }
.table-title { font-size: 15px; font-weight: 600; color: #e2e8f0; }
.table-total { font-size: 12px; color: #64748b; }
.table-actions { display: flex; gap: 8px; }

/* Tech table */
.tech-table { width: 100%; border-collapse: collapse; }
.tech-table th {
  font-size: 11px; text-transform: uppercase; letter-spacing: 0.8px;
  color: #64748b; font-weight: 600; padding: 10px 12px;
  text-align: left; border-bottom: 1px solid rgba(0, 212, 255, 0.06);
}
.tech-table td {
  padding: 10px 12px; font-size: 13px; color: #cbd5e1;
  border-bottom: 1px solid rgba(0, 212, 255, 0.04);
}
.tech-table tbody tr:hover { background: rgba(0, 212, 255, 0.03); }

.cell-primary { color: #e2e8f0; font-weight: 500; }
.cell-mono { font-family: 'JetBrains Mono', monospace; font-size: 12px; color: #64748b; }
.cell-vendor { font-size: 11px; font-weight: 600; letter-spacing: 0.5px; color: #94a3b8; }
.cell-model { font-size: 12px; color: #94a3b8; }
.cell-text { color: #94a3b8; }

/* Role badges */
.role-badge {
  font-size: 11px; padding: 2px 10px; border-radius: 4px; font-weight: 500;
}
.role-core { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
.role-ac { background: rgba(245, 158, 11, 0.1); color: #f59e0b; }
.role-aggregation { background: rgba(0, 212, 255, 0.08); color: #00d4ff; }
.role-access { background: rgba(100, 116, 139, 0.1); color: #94a3b8; }
.role-ap { background: rgba(16, 185, 129, 0.1); color: #10b981; }

/* Status indicators */
.status-dot {
  display: inline-block; width: 7px; height: 7px; border-radius: 50%; margin-right: 6px; vertical-align: middle;
}
.status-dot.online { background: #10b981; box-shadow: 0 0 8px rgba(16,185,129,0.5); }
.status-dot.offline { background: #ef4444; box-shadow: 0 0 8px rgba(239,68,68,0.3); }
.status-text { font-size: 12px; vertical-align: middle; }
.status-text.online { color: #10b981; }
.status-text.offline { color: #ef4444; }

.empty-row { text-align: center; color: #475569 !important; padding: 40px !important; }

.pagination-wrap {
  display: flex; justify-content: flex-end; margin-top: 16px;
}
</style>
