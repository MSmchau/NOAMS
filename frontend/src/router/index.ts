import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
        meta: { title: '仪表盘' },
      },
      // 设备管理
      {
        path: 'devices',
        name: 'Devices',
        component: () => import('@/views/devices/DeviceListView.vue'),
        meta: { title: '设备管理' },
      },
      {
        path: 'devices/create',
        name: 'DeviceCreate',
        component: () => import('@/views/devices/DeviceFormView.vue'),
        meta: { title: '添加设备' },
      },
      {
        path: 'devices/:id/edit',
        name: 'DeviceEdit',
        component: () => import('@/views/devices/DeviceFormView.vue'),
        meta: { title: '编辑设备' },
      },
      {
        path: 'devices/:id',
        name: 'DeviceDetail',
        component: () => import('@/views/devices/DeviceDetailView.vue'),
        meta: { title: '设备详情' },
      },
      // 自动巡检
      {
        path: 'inspections',
        name: 'Inspections',
        component: () => import('@/views/inspections/InspectionListView.vue'),
        meta: { title: '巡检记录' },
      },
      {
        path: 'inspections/report',
        name: 'InspectionReport',
        component: () => import('@/views/inspections/InspectionReportView.vue'),
        meta: { title: '巡检报告' },
      },
      // 配置管理
      {
        path: 'configs/backup',
        name: 'ConfigBackup',
        component: () => import('@/views/configs/ConfigBackupView.vue'),
        meta: { title: '配置备份' },
      },
      {
        path: 'configs/history',
        name: 'ConfigHistory',
        component: () => import('@/views/configs/ConfigHistoryView.vue'),
        meta: { title: '版本历史' },
      },
      // 告警管理
      {
        path: 'alerts',
        name: 'Alerts',
        component: () => import('@/views/alerts/AlertListView.vue'),
        meta: { title: '告警管理' },
      },
      // 定时任务
      {
        path: 'tasks',
        name: 'Tasks',
        component: () => import('@/views/tasks/TaskListView.vue'),
        meta: { title: '任务列表' },
      },
      {
        path: 'tasks/logs',
        name: 'TaskLogs',
        component: () => import('@/views/tasks/TaskLogsView.vue'),
        meta: { title: '执行日志' },
      },
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  if (to.meta.requiresAuth !== false && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router
