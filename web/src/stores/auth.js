import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)       // { accountEmail, nickName, headImageUri, ... }
  const sitename = ref('')
  const checked = ref(false)   // 是否已向后端确认过登录态

  async function fetchMe() {
    try {
      const data = await api.me()
      user.value = data.user && data.user.accountEmail ? data.user : null
      sitename.value = data.sitename || ''
    } catch (e) {
      user.value = null
    } finally {
      checked.value = true
    }
    return user.value
  }

  async function login(email, password) {
    const data = await api.login(email, password)
    user.value = data
    checked.value = true
    return data
  }

  async function logout() {
    try { await api.logout() } catch (e) { /* ignore */ }
    user.value = null
  }

  const isLoggedIn = () => !!(user.value && user.value.accountEmail)

  return { user, sitename, checked, fetchMe, login, logout, isLoggedIn }
})
