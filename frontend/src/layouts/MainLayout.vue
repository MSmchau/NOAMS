<template>
  <el-container style="height: 100vh">
    <!-- Sidebar -->
    <el-aside width="220px" style="background-color: #1d1e1f">
      <div class="sidebar-header">
        <h2 class="sidebar-title">NOAMS</h2>
        <span class="sidebar-subtitle">网络运维自动化管理</span>
      </div>
      <el-menu
        :default-active="route.path"
        background-color="#1d1e1f"
        text-color="#bfcbd9"
        active-text-color="#409eff"
        router
      >
        <el-menu-item index="/dashboard">
          <el-icon><Monitor /></el-icon>
          <span>仪表盘</span>
        </el-menu-item>
        <el-menu-item index="/devices">
          <el-icon><Server /></el-icon>
          <span>设备管理</span>
        </el-menu-item>
        <el-sub-menu index="inspection">
          <template #title>
            <el-icon><Search /></el-icon>
            <span>自动巡检</span>
          </template>
          <el-menu-item index="/inspections">巡检记录</el-menu-item>
          <el-menu-item index="/inspections/report">巡检报告</el-menu-item>
        </el-sub-menu>
        <el-sub-menu index="config">
          <template #title>
            <el-icon><Document /></el-icon>
            <span>配置管理</span>
          </template>
          <el-menu-item index="/configs/backup">配置备份</el-menu-item>
          <el-menu-item index="/configs/history">版本历史</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/alerts">
          <el-icon><WarningFilled /></el-icon>
          <span>告警管理</span>
        </el-menu-item>
        <el-sub-menu index="tasks">
          <template #title>
            <el-icon><Clock /></el-icon>
            <span>定时任务</span>
          </template>
          <el-menu-item index="/tasks">任务列表</el-menu-item>
          <el-menu-item index="/tasks/logs">执行日志</el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-aside>

    <!-- Main content -->
    <el-container>
      <!-- Header -->
      <el-header style="border-bottom: 1px solid #e4e7ed; display: flex; align-items: center; justify-content: space-between; height: 60px;">
        <div class="header-left">
          <el-breadcrumb>
            <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="route.meta.title">{{ route.meta.title }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right" style="display: flex; align-items: center; gap: 16px;">
          <el-tag v-if="userStore.userInfo?.role === 'admin'" type="danger" size="small">管理员</el-tag>
          <el-dropdown trigger="click">
            <span style="cursor: pointer; display: flex; align-items: center; gap: 6px;">
              <el-avatar :size="28">{{ userStore.userInfo?.nickname?.[0] || 'U' }}</el-avatar>
              <span>{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item disabled>个人中心</el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- Content -->
      <el-main style="background-color: #f5f7fa; padding: 20px;">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

function handleLogout() {
  ElMessageBox.confirm('确认退出登录吗？', '提示', {
    type: 'warning',
  }).then(() => {
    userStore.logout()
    router.push('/login')
  }).catch(() => {})
}
</script>

<style scoped>
.sidebar-header {
  padding: 20px 16px;
  text-align: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}
.sidebar-title {
  color: #fff;
  font-size: 20px;
  margin: 0;
  letter-spacing: 2px;
}
.sidebar-subtitle {
  color: #8a8f9d;
  font-size: 11px;
  display: block;
  margin-top: 4px;
}
.el-aside {
  overflow-y: auto;
  overflow-x: hidden;
}
.el-menu {
  border-right: none;
}
</style>
