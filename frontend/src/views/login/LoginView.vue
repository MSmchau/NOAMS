<template>
  <div class="login-container">
    <!-- Animated background -->
    <div class="bg-grid"></div>
    <div class="bg-orb bg-orb-1"></div>
    <div class="bg-orb bg-orb-2"></div>
    <div class="bg-orb bg-orb-3"></div>

    <div class="login-card">
      <div class="login-header">
        <div class="logo-icon">
          <svg viewBox="0 0 48 48" width="48" height="48">
            <circle cx="24" cy="24" r="22" fill="none" stroke="#00d4ff" stroke-width="1.5" opacity="0.3"/>
            <circle cx="24" cy="24" r="14" fill="none" stroke="#00d4ff" stroke-width="1.5" opacity="0.5"/>
            <circle cx="24" cy="24" r="6" fill="#00d4ff" opacity="0.8"/>
            <line x1="24" y1="2" x2="24" y2="10" stroke="#00d4ff" stroke-width="1.5" opacity="0.6"/>
            <line x1="24" y1="38" x2="24" y2="46" stroke="#00d4ff" stroke-width="1.5" opacity="0.6"/>
            <line x1="2" y1="24" x2="10" y2="24" stroke="#00d4ff" stroke-width="1.5" opacity="0.6"/>
            <line x1="38" y1="24" x2="46" y2="24" stroke="#00d4ff" stroke-width="1.5" opacity="0.6"/>
          </svg>
        </div>
        <h1 class="login-title"><span class="gradient-text">NOAMS</span></h1>
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
            placeholder="用户名"
            :prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
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
            {{ loading ? '验证中...' : '进入系统' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <span>默认账号 admin / admin123</span>
      </div>

      <div class="login-version">v1.0.0 · Network Operations Automation System</div>
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
  background: #0a0e17;
  overflow: hidden;
}

/* Grid background */
.bg-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(0, 212, 255, 0.03) 1px, transparent 1px),
    linear-gradient(90deg, rgba(0, 212, 255, 0.03) 1px, transparent 1px);
  background-size: 60px 60px;
}

/* Glowing orbs */
.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
}
.bg-orb-1 {
  width: 400px; height: 400px;
  background: rgba(0, 212, 255, 0.08);
  top: -100px; right: -100px;
  animation: orbFloat 12s ease-in-out infinite;
}
.bg-orb-2 {
  width: 300px; height: 300px;
  background: rgba(124, 58, 237, 0.06);
  bottom: -80px; left: -80px;
  animation: orbFloat 16s ease-in-out infinite reverse;
}
.bg-orb-3 {
  width: 200px; height: 200px;
  background: rgba(0, 212, 255, 0.04);
  top: 40%; left: 50%;
  animation: orbFloat 10s ease-in-out infinite;
}

@keyframes orbFloat {
  0%, 100% { transform: translate(0, 0) scale(1); }
  33% { transform: translate(30px, -30px) scale(1.1); }
  66% { transform: translate(-20px, 20px) scale(0.9); }
}

/* Login card */
.login-card {
  position: relative;
  width: 420px;
  padding: 48px 40px 36px;
  background: rgba(19, 26, 43, 0.85);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border: 1px solid rgba(0, 212, 255, 0.08);
  border-radius: 16px;
  box-shadow:
    0 0 40px rgba(0, 0, 0, 0.4),
    0 0 80px rgba(0, 212, 255, 0.04);
  animation: cardAppear 0.6s ease-out;
}

@keyframes cardAppear {
  from { opacity: 0; transform: translateY(20px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.login-card::before {
  content: '';
  position: absolute;
  inset: -1px;
  border-radius: 17px;
  padding: 1px;
  background: linear-gradient(135deg, rgba(0, 212, 255, 0.2), transparent 40%, transparent 60%, rgba(124, 58, 237, 0.2));
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  pointer-events: none;
}

.login-header {
  text-align: center;
  margin-bottom: 36px;
}

.logo-icon {
  margin-bottom: 16px;
  animation: pulseIcon 3s ease-in-out infinite;
}
@keyframes pulseIcon {
  0%, 100% { opacity: 0.8; }
  50% { opacity: 1; }
}

.login-title {
  margin: 0;
  font-size: 32px;
  letter-spacing: 6px;
  font-weight: 800;
}
.gradient-text {
  background: linear-gradient(135deg, #00d4ff, #7c3aed, #60a5fa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.login-desc {
  margin: 10px 0 0;
  color: #64748b;
  font-size: 14px;
  letter-spacing: 2px;
}

.login-form {
  margin-bottom: 20px;
}

.login-btn {
  width: 100%;
  height: 44px !important;
  font-size: 15px !important;
  letter-spacing: 4px !important;
  border-radius: 8px !important;
}

.login-footer {
  text-align: center;
  color: #475569;
  font-size: 12px;
}

.login-version {
  text-align: center;
  color: rgba(100, 116, 139, 0.4);
  font-size: 11px;
  margin-top: 24px;
  letter-spacing: 0.5px;
}
</style>
