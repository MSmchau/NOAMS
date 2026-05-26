<template>
  <div class="page-container">
    <!-- Search Card -->
    <el-card shadow="never">
      <el-form :model="filters" inline size="default" class="search-form">
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
    </el-card>

    <!-- Table Card -->
    <el-card shadow="never">
      <template #header>
        <div class="table-toolbar">
          <div class="toolbar-left">
            <span class="toolbar-title">设备列表</span>
            <el-tag type="info" size="small">共 {{ total }} 台</el-tag>
          </div>
          <div class="toolbar-right">
            <el-button type="primary" size="small" @click="handleAdd">
              <el-icon><Plus /></el-icon>添加设备
            </el-button>
            <el-button size="small" @click="loadDevices">
              <el-icon><Refresh /></el-icon>刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="devices" v-loading="loading" stripe size="small" @row-click="handleDetail" style="cursor:pointer">
        <el-table-column prop="name" label="设备名称" min-width="150">
          <template #default="{ row }">
            <span style="color:#262626;font-weight:500">{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column label="管理IP" width="140">
          <template #default="{ row }">
            <code style="font-size:12px;color:#595959">{{ row.management_ip }}</code>
          </template>
        </el-table-column>
        <el-table-column label="厂商" width="80">
          <template #default="{ row }">
            <span style="font-weight:500;color:#595959">{{ (row.vendor || '').toUpperCase() || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="model" label="型号" width="140" />
        <el-table-column label="角色" width="80">
          <template #default="{ row }">
            <el-tag :type="roleTagType(row.role)" size="small">{{ roleLabel(row.role) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="building" label="楼栋" width="100" />
        <el-table-column label="状态" width="80" align="center">
          <template #default="{ row }">
            <span class="status-badge" :class="row.status ? 'online' : 'offline'">
              <span class="dot"></span>{{ row.status ? '在线' : '离线' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right" @click.stop>
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click.stop="handleDetail(row)">详情</el-button>
            <el-button type="primary" link size="small" @click.stop="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确认删除该设备？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link size="small" @click.stop>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="total > pageSize" class="pagination-wrap">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadDevices"
        />
      </div>
    </el-card>
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

function roleTagType(role?: string) {
  const map: Record<string, 'danger' | 'warning' | 'info' | 'success'> = { core: 'danger', ac: 'warning', aggregation: 'info', access: 'info', ap: 'success' }
  return map[role || ''] || 'info'
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
.search-form {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0;
}

.search-form :deep(.el-form-item) {
  margin-bottom: 0;
}

.table-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-title {
  font-size: 15px;
  font-weight: 600;
  color: #262626;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
