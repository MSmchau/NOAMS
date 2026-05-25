<template>
  <div class="device-list">
    <!-- Search bar -->
    <el-card shadow="hover" style="margin-bottom: 20px;">
      <el-form :model="filters" inline size="default">
        <el-form-item label="设备名称">
          <el-input v-model="filters.name" placeholder="搜索设备名称" clearable @clear="search" @keyup.enter="search" />
        </el-form-item>
        <el-form-item label="管理IP">
          <el-input v-model="filters.management_ip" placeholder="搜索IP地址" clearable @clear="search" @keyup.enter="search" />
        </el-form-item>
        <el-form-item label="厂商">
          <el-select v-model="filters.vendor" placeholder="选择厂商" clearable @change="search">
            <el-option label="H3C" value="h3c" />
            <el-option label="华为" value="huawei" />
            <el-option label="思科" value="cisco" />
            <el-option label="锐捷" value="ruijie" />
          </el-select>
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="filters.role" placeholder="选择角色" clearable @change="search">
            <el-option label="核心" value="core" />
            <el-option label="AC" value="ac" />
            <el-option label="本体" value="aggregation" />
            <el-option label="接入" value="access" />
            <el-option label="AP" value="ap" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="选择状态" clearable @change="search">
            <el-option label="在线" :value="1" />
            <el-option label="离线" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">查询</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- Device table -->
    <el-card shadow="hover">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <span>设备列表</span>
          <div>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>添加设备
            </el-button>
            <el-button @click="handleRefresh">
              <el-icon><Refresh /></el-icon>刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table :data="devices" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="name" label="设备名称" min-width="160">
          <template #default="{ row }">
            <el-link type="primary" @click="handleDetail(row)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="management_ip" label="管理IP" width="140" />
        <el-table-column prop="vendor" label="厂商" width="80">
          <template #default="{ row }">
            <el-tag size="small">{{ row.vendor?.toUpperCase() || '-' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="model" label="型号" width="160" />
        <el-table-column prop="role" label="角色" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.role === 'core'" type="danger" size="small">核心</el-tag>
            <el-tag v-else-if="row.role === 'ac'" type="warning" size="small">AC</el-tag>
            <el-tag v-else-if="row.role === 'aggregation'" size="small">本体</el-tag>
            <el-tag v-else-if="row.role === 'access'" type="info" size="small">接入</el-tag>
            <el-tag v-else-if="row.role === 'ap'" type="success" size="small">AP</el-tag>
            <span v-else>{{ row.role }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="building" label="楼栋" width="120" />
        <el-table-column prop="group?.name" label="分组" width="120" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.status === 1" type="success" size="small">在线</el-tag>
            <el-tag v-else type="danger" size="small">离线</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleDetail(row)">详情</el-button>
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确认删除该设备？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div style="display: flex; justify-content: flex-end; margin-top: 16px;">
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

const filters = reactive({
  name: '',
  management_ip: '',
  vendor: '',
  role: '',
  status: '' as string | number,
})

async function loadDevices() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: page.value,
      page_size: pageSize.value,
    }
    if (filters.name) params.name = filters.name
    if (filters.management_ip) params.management_ip = filters.management_ip
    if (filters.vendor) params.vendor = filters.vendor
    if (filters.role) params.role = filters.role
    if (filters.status !== '') params.status = filters.status

    const res = await getDevices(params)
    devices.value = res.data.list || []
    total.value = res.data.total || 0
  } catch {
    devices.value = []
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  loadDevices()
}

function resetSearch() {
  filters.name = ''
  filters.management_ip = ''
  filters.vendor = ''
  filters.role = ''
  filters.status = ''
  page.value = 1
  loadDevices()
}

function handleAdd() {
  router.push('/devices/create')
}

function handleDetail(row: Device) {
  router.push(`/devices/${row.id}`)
}

function handleEdit(row: Device) {
  router.push(`/devices/${row.id}/edit`)
}

async function handleDelete(row: Device) {
  try {
    await deleteDevice(row.id)
    ElMessage.success('设备已删除')
    loadDevices()
  } catch {
    // already handled by interceptor
  }
}

function handleRefresh() {
  loadDevices()
}

onMounted(() => {
  loadDevices()
})
</script>
