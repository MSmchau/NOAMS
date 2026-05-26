<template>
  <div class="app-layout-wrapper">
    <!-- Sidebar -->
    <div class="app-sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-logo">
        <div class="logo-icon">
          <svg viewBox="0 0 36 36" width="28" height="28">
            <circle cx="18" cy="18" r="16" fill="none" stroke="#1890ff" stroke-width="1.5" opacity="0.3"/>
            <circle cx="18" cy="18" r="10" fill="none" stroke="#1890ff" stroke-width="1.5" opacity="0.5"/>
            <circle cx="18" cy="18" r="4" fill="#1890ff" opacity="0.9"/>
          </svg>
        </div>
        <span v-show="!sidebarCollapsed" class="logo-text">NOAMS</span>
      </div>

      <el-scrollbar class="sidebar-menu-wrap">
        <el-menu
          :default-active="route.path"
          :collapse="sidebarCollapsed"
          background-color="#001529"
          text-color="rgba(255,255,255,0.65)"
          active-text-color="#ffffff"
          router
        >
          <el-menu-item index="/dashboard">
            <el-icon><Monitor /></el-icon>
            <template #title>仪表盘</template>
          </el-menu-item>

          <el-menu-item index="/devices">
            <el-icon><Server /></el-icon>
            <template #title>设备管理</template>
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
            <template #title>告警中心</template>
          </el-menu-item>

          <el-sub-menu index="task">
            <template #title>
              <el-icon><Clock /></el-icon>
              <span>定时任务</span>
            </template>
            <el-menu-item index="/tasks">任务列表</el-menu-item>
            <el-menu-item index="/tasks/logs">执行日志</el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-scrollbar>
    </div>

    <!-- Main Area -->
    <div class="app-main" :class="{ collapsed: sidebarCollapsed }">
      <!-- Header -->
      <header class="app-header">
        <div class="header-left">
          <el-button text class="collapse-btn" @click="toggleSidebar">
            <el-icon :size="18">
              <Fold v-if="!sidebarCollapsed" />
              <Expand v-else />
            </el-icon>
          </el-button>
          <span class="page-title">{{ route.meta?.title || '仪表盘' }}</span>
        </div>
        <div class="header-right">
          <!-- Notification -->
          <el-badge :value="notificationCount" :hidden="notificationCount === 0" class="notification-badge">
            <el-button text class="header-icon-btn" @click="handleNotification">
              <el-icon :size="18"><Bell /></el-icon>
            </el-button>
          </el-badge>

          <!-- User -->
          <el-dropdown trigger="click" @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="28" icon="UserFilled" class="user-avatar" />
              <span class="user-name">{{ userStore.userInfo?.nickname || userStore.userInfo?.username || 'admin' }}</span>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>个人中心
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <el-icon><SwitchButton /></el-icon>退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </header>

      <!-- Content -->
      <main class="app-content">
        <router-view />
      </main>

      <!-- Footer -->
      <footer class="app-footer">
        <span>NOAMS v1.0.0 &copy; {{ new Date().getFullYear() }} 网络运维自动化管理系统</span>
      </footer>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const sidebarCollapsed = ref(false)
const notificationCount = ref(0)

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

function handleNotification() {
  // TODO: navigate to notification center
}

function handleCommand(cmd: string) {
  if (cmd === 'logout') {
    ElMessageBox.confirm('确认退出登录？', '提示', { type: 'warning', confirmButtonText: '确认', cancelButtonText: '取消' })
      .then(() => { userStore.logout(); router.push('/login') })
      .catch(() => {})
  }
}
</script>

<style scoped>
.app-layout-wrapper {
  display: flex;
  height: 100vh;
  overflow: hidden;
}

/* ===== Sidebar ===== */
.app-sidebar {
  width: 240px;
  min-width: 240px;
  background: #001529;
  display: flex;
  flex-direction: column;
  transition: all 0.2s ease;
  z-index: 100;
  overflow: hidden;
}

.app-sidebar.collapsed {
  width: 64px;
  min-width: 64px;
}

.sidebar-logo {
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  letter-spacing: 2px;
  color: #ffffff;
  white-space: nowrap;
}

.sidebar-menu-wrap {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.app-sidebar :deep(.el-menu) {
  border-right: none !important;
}

.app-sidebar :deep(.el-menu-item) {
  margin: 2px 8px;
  width: calc(100% - 16px);
  border-radius: 6px;
  font-size: 14px;
  height: 44px;
  line-height: 44px;
  transition: all 0.2s ease;
}

.app-sidebar :deep(.el-menu-item:hover) {
  background-color: rgba(255, 255, 255, 0.08) !important;
}

.app-sidebar :deep(.el-menu-item.is-active) {
  background-color: #1890ff !important;
  color: #ffffff !important;
}

.app-sidebar :deep(.el-sub-menu__title) {
  margin: 2px 8px;
  width: calc(100% - 16px);
  border-radius: 6px;
  font-size: 14px;
  height: 44px;
  line-height: 44px;
  transition: all 0.2s ease;
}

.app-sidebar :deep(.el-sub-menu__title:hover) {
  background-color: rgba(255, 255, 255, 0.08) !important;
}

.app-sidebar :deep(.el-menu--collapse .el-menu-item) {
  margin: 2px 4px;
  width: calc(100% - 8px);
  border-radius: 6px;
}

.app-sidebar :deep(.el-menu--collapse .el-sub-menu__title) {
  margin: 2px 4px;
  width: calc(100% - 8px);
  border-radius: 6px;
}

.app-sidebar .el-scrollbar__view {
  height: 100%;
}

/* ===== Main Area ===== */
.app-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  transition: all 0.2s ease;
}

/* ===== Header ===== */
.app-header {
  height: 56px;
  min-height: 56px;
  background: #ffffff;
  border-bottom: 1px solid #f0f0f0;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  position: sticky;
  top: 0;
  z-index: 50;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.collapse-btn {
  font-size: 16px;
  color: #595959;
  padding: 4px;
  &:hover { color: #1890ff; }
}

.page-title {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon-btn {
  padding: 6px;
  color: #595959;
  &:hover { color: #1890ff; }
}

.notification-badge :deep(.el-badge__content) {
  border: none;
  font-size: 11px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 6px;
  transition: background 0.2s;

  &:hover {
    background: #f5f5f5;
  }
}

.user-avatar {
  background: #1890ff;
  flex-shrink: 0;
}

.user-name {
  font-size: 14px;
  color: #262626;
  font-weight: 500;
}

/* ===== Content ===== */
.app-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  background: #f0f2f5;
}

/* ===== Footer ===== */
.app-footer {
  height: 40px;
  min-height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #ffffff;
  border-top: 1px solid #f0f0f0;
  font-size: 12px;
  color: #8c8c8c;
}

/* ===== Responsive ===== */
@media (max-width: 768px) {
  .app-sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    transform: translateX(-100%);
  }
  .app-sidebar:not(.collapsed) {
    transform: translateX(0);
  }
  .app-main {
    margin-left: 0;
  }
  .user-name {
    display: none;
  }
}
</style>
