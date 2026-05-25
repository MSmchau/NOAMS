import { defineStore } from 'pinia'
import { ref } from 'vue'
import { login, getProfile, type LoginData, type LoginResult } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref<LoginResult['user'] | null>(
    JSON.parse(localStorage.getItem('user') || 'null')
  )

  async function loginAction(data: LoginData) {
    const res = await login(data)
    const result = res.data as unknown as LoginResult
    token.value = result.token
    userInfo.value = result.user
    localStorage.setItem('token', result.token)
    localStorage.setItem('user', JSON.stringify(result.user))
    return result
  }

  async function fetchProfile() {
    const res = await getProfile()
    userInfo.value = res.data as unknown as LoginResult['user']
    localStorage.setItem('user', JSON.stringify(userInfo.value))
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    token,
    userInfo,
    loginAction,
    fetchProfile,
    logout,
  }
})
