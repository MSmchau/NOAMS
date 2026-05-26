<template>
  <div class="login-container">
    <div class="login-bg"></div>

    <div class="login-card-wrapper">
      <div class="login-card">
        <div class="login-header">
          <div class="login-logo">
            <svg viewBox="0 0 48 48" width="52" height="52">
              <circle cx="24" cy="24" r="22" fill="none" stroke="#1890ff" stroke-width="1.5" opacity="0.3"/>
              <circle cx="24" cy="24" r="14" fill="none" stroke="#1890ff" stroke-width="1.5" opacity="0.5"/>
              <circle cx="24" cy="24" r="6" fill="#1890ff" opacity="0.9"/>
              <line x1="24" y1="2" x2="24" y2="8" stroke="#1890ff" stroke-width="1.5" opacity="0.5"/>
              <line x1="24" y1="40" x2="24" y2="46" stroke="#1890ff" stroke-width="1.5" opacity="0.5"/>
              <line x1="2" y1="24" x2="8" y2="24" stroke="#1890ff" stroke-width="1.5" opacity="0.5"/>
              <line x1="40" y1="24" x2="46" y2="24" stroke="#1890ff" stroke-width="1.5" opacity="0.5"/>
            </svg>
          </div>
          <h1 class="login-title">NOAMS</h1>
          <p class="login-desc">网络运维自动化管理系统</p>
        </div>

        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          size="large"
          class="login-form"
          @keyup.enter="handleLogin"
        >
          <el-form-item prop="username">
            <el-input
              v-model="form.username"
              placeholder="请输入用户名"
              :prefix-icon="User"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              class="login-btn"
              @click="handleLogin"
            >
              {{ loading ? '登录中...' : '登 录' }}
            </el-button>
          </el-form-item>
        </el-form>

        <div class="login-footer">
          <span class="footer-tip">默认账号: admin / admin123</span>
        </div>

        <div class="login-version">
          Network Operations Automation System v1.0.0
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { User, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  username: 'admin',
  password: 'admin123',
})

const rules: FormRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    await userStore.loginAction(form)
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (err: any) {
    ElMessage.error(err.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
}

/* Grid background overlay */
.login-bg {
  position: absolute;
  inset: 0;
  z-index: 1;
}

.login-bg::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background-image:
    radial-gradient(circle at 25% 25%, rgba(255, 255, 255, 0.08) 1px, transparent 1px),
    radial-gradient(circle at 75% 75%, rgba(255, 255, 255, 0.08) 1px, transparent 1px);
  background-size: 50px 50px;
  animation: bgFloat 30s infinite linear;
}

@keyframes bgFloat {
  0% { transform: translate(-50%, -50%) rotate(0deg); }
  100% { transform: translate(-50%, -50%) rotate(360deg); }
}

/* Card */
.login-card-wrapper {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 400px;
  padding: 20px;
  animation: cardEnter 0.5s ease-out;
}

@keyframes cardEnter {
  from { opacity: 0; transform: translateY(24px); }
  to { opacity: 1; transform: translateY(0); }
}

.login-card {
  background: rgba(255, 255, 255, 0.98);
  border-radius: 12px;
  padding: 40px 36px 28px;
  box-shadow: 0 8px 40px rgba(0, 0, 0, 0.12);
  backdrop-filter: blur(10px);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  margin-bottom: 16px;
  display: flex;
  justify-content: center;
}

.login-title {
  margin: 0 0 8px;
  font-size: 28px;
  font-weight: 700;
  color: #262626;
  letter-spacing: 4px;
}

.login-desc {
  margin: 0;
  font-size: 14px;
  color: #8c8c8c;
  letter-spacing: 1px;
}

.login-form {
  margin-bottom: 16px;
}

.login-form :deep(.el-form-item) {
  margin-bottom: 22px;
}

.login-form :deep(.el-input__wrapper) {
  border-radius: 8px;
  padding: 4px 12px;
  box-shadow: 0 0 0 1px #d9d9d9 inset !important;
}

.login-form :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #1890ff inset !important;
}

.login-form :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #1890ff inset !important;
}

.login-form :deep(.el-input__inner) {
  height: 42px;
  font-size: 15px;
}

.login-btn {
  width: 100%;
  height: 44px;
  font-size: 16px;
  letter-spacing: 4px;
  border-radius: 8px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border: none;
  transition: all 0.3s;
}

.login-btn:hover {
  background: linear-gradient(135deg, #40a9ff, #1890ff);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.4);
}

.login-btn:active {
  transform: translateY(0);
}

.login-footer {
  text-align: center;
  border-top: 1px solid #f0f0f0;
  padding-top: 20px;
}

.footer-tip {
  font-size: 12px;
  color: #8c8c8c;
}

.login-version {
  text-align: center;
  font-size: 11px;
  color: #d9d9d9;
  margin-top: 20px;
}

/* Responsive */
@media (max-width: 576px) {
  .login-card-wrapper {
    padding: 16px;
  }
  .login-card {
    padding: 32px 24px 24px;
  }
  .login-title {
    font-size: 24px;
  }
}
</style>
