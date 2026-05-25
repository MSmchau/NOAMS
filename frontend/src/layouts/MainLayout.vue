<template>
  <div class="layout">
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <div class="brand">
          <svg viewBox="0 0 36 36" width="32" height="32">
            <circle cx="18" cy="18" r="16" fill="none" stroke="#00d4ff" stroke-width="1.2" opacity="0.3"/>
            <circle cx="18" cy="18" r="10" fill="none" stroke="#00d4ff" stroke-width="1.2" opacity="0.5"/>
            <circle cx="18" cy="18" r="4" fill="#00d4ff" opacity="0.9"/>
          </svg>
          <div class="brand-text">
            <span class="brand-name">NOAMS</span>
            <span class="brand-sub">运维管理平台</span>
          </div>
        </div>
      </div>

      <el-scrollbar class="sidebar-menu-wrap">
        <el-menu
          :default-active="route.path"
          text-color="#94a3b8"
          active-text-color="#00d4ff"
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
          <el-menu-item index="/inspections">
            <el-icon><Search /></el-icon>
            <span>自动巡检</span>
          </el-menu-item>
          <el-menu-item index="/inspections/report">
            <el-icon><DataBoard /></el-icon>
            <span>巡检报告</span>
          </el-menu-item>
          <el-menu-item index="/configs/backup">
            <el-icon><Document /></el-icon>
            <span>配置管理</span>
          </el-menu-item>
          <el-menu-item index="/alerts">
            <el-icon><WarningFilled /></el-icon>
            <span>告警中心</span>
          </el-menu-item>
          <el-menu-item index="/tasks">
            <el-icon><Clock /></el-icon>
            <span>定时任务</span>
          </el-menu-item>
        </el-menu>
      </el-scrollbar>

      <div class="sidebar-footer">
        <div class="status-dot"></div>
        <span class="user-name">{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</span>
        <el-dropdown trigger="click" @command="handleCommand">
          <el-icon class="more-btn"><MoreFilled /></el-icon>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人中心</el-dropdown-item>
              <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </aside>

    <!-- Main -->
    <div class="main-area">
      <header class="topbar">
        <div class="topbar-left">
          <span class="page-title">{{ route.meta.title || 'Dashboard' }}</span>
        </div>
        <div class="topbar-right">
          <div class="header-time">{{ currentTime }}</div>
          <el-tag v-if="userStore.userInfo?.role === 'admin'" size="small" class="role-tag">ADMIN</el-tag>
        </div>
      </header>

      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const currentTime = ref('')
let timer: ReturnType<typeof setInterval>

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleString('zh-CN', {
    hour12: false,
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  })
}

function handleCommand(cmd: string) {
  if (cmd === 'logout') {
    ElMessageBox.confirm('确认退出登录？', '提示', { type: 'warning' })
      .then(() => { userStore.logout(); router.push('/login') })
      .catch(() => {})
  }
}

onMounted(() => { updateTime(); timer = setInterval(updateTime, 1000) })
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.layout {
  display: flex;
  height: 100vh;
  background: #0a0e17;
}

/* === Sidebar === */
.sidebar {
  width: 220px;
  min-width: 220px;
  background: #0d111c;
  border-right: 1px solid rgba(0, 212, 255, 0.06);
  display: flex;
  flex-direction: column;
  position: relative;
  z-index: 10;
}

.sidebar-header {
  padding: 20px 16px;
  border-bottom: 1px solid rgba(0, 212, 255, 0.06);
  margin-bottom: 4px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-text {
  display: flex;
  flex-direction: column;
}
.brand-name {
  font-size: 18px;
  font-weight: 800;
  letter-spacing: 3px;
  background: linear-gradient(135deg, #00d4ff, #7c3aed);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  line-height: 1.2;
}
.brand-sub {
  font-size: 10px;
  color: #475569;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.sidebar-menu-wrap {
  flex: 1;
  overflow-y: auto;
  padding: 4px 8px;
}

.sidebar :deep(.el-menu) {
  background: transparent;
  border: none;
}

.sidebar :deep(.el-menu-item) {
  border-radius: 8px;
  margin: 2px 0;
  height: 42px;
  line-height: 42px;
  font-size: 13px;
  transition: all 0.25s ease;

  &:hover {
    background: rgba(0, 212, 255, 0.06) !important;
    color: #e2e8f0 !important;
  }

  &.is-active {
    background: rgba(0, 212, 255, 0.1) !important;
    color: #00d4ff !important;
    position: relative;

    &::before {
      content: '';
      position: absolute;
      left: -8px;
      top: 50%;
      transform: translateY(-50%);
      width: 3px;
      height: 20px;
      background: #00d4ff;
      border-radius: 0 3px 3px 0;
      box-shadow: 0 0 10px rgba(0, 212, 255, 0.5);
    }
  }

  .el-icon { font-size: 18px; margin-right: 10px; }
}

/* Footer */
.sidebar-footer {
  padding: 14px 16px;
  border-top: 1px solid rgba(0, 212, 255, 0.06);
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 13px;
  color: #64748b;
}

.status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.5);
  animation: pulse 2s ease-in-out infinite;
}
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.user-name {
  flex: 1;
  color: #94a3b8;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.more-btn {
  cursor: pointer;
  color: #64748b;
  &:hover { color: #94a3b8; }
}

/* === Main Area === */
.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.topbar {
  height: 56px;
  min-height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 28px;
  border-bottom: 1px solid rgba(0, 212, 255, 0.06);
  background: rgba(10, 14, 23, 0.8);
  backdrop-filter: blur(12px);
}

.topbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title {
  font-size: 15px;
  font-weight: 600;
  color: #e2e8f0;
  letter-spacing: 0.5px;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-time {
  font-family: 'JetBrains Mono', monospace;
  font-size: 13px;
  color: #64748b;
  letter-spacing: 0.5px;
}

.role-tag {
  background: rgba(0, 212, 255, 0.1) !important;
  border: 1px solid rgba(0, 212, 255, 0.2) !important;
  color: #00d4ff !important;
  font-size: 10px !important;
  letter-spacing: 2px !important;
}

.content {
  flex: 1;
  padding: 24px 28px;
  overflow-y: auto;
  background:
    radial-gradient(ellipse at 0% 0%, rgba(0, 212, 255, 0.015) 0%, transparent 50%),
    radial-gradient(ellipse at 100% 100%, rgba(124, 58, 237, 0.015) 0%, transparent 50%);
}
</style>
