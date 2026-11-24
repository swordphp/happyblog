<script setup>
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import LangSwitch from '@/components/LangSwitch.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const { t } = useI18n()

const form = reactive({ email: '', password: '' })
const loading = ref(false)

async function onSubmit() {
  if (!form.email || !form.password) {
    ElMessage.warning(t('login.needInput'))
    return
  }
  loading.value = true
  try {
    await auth.login(form.email, form.password)
    ElMessage.success(t('login.ok'))
    router.push(route.query.redirect || { name: 'admin-dashboard' })
  } catch (e) {
    ElMessage.error(e.message || t('login.failed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-wrap">
    <el-card class="login-card">
      <div class="login-top"><LangSwitch /></div>
      <h2 class="login-title">{{ $t('login.title') }}</h2>
      <el-form label-position="top" @submit.prevent="onSubmit">
        <el-form-item :label="$t('login.email')">
          <el-input v-model="form.email" placeholder="admin@example.com" @keyup.enter="onSubmit" />
        </el-form-item>
        <el-form-item :label="$t('login.password')">
          <el-input v-model="form.password" type="password" show-password @keyup.enter="onSubmit" />
        </el-form-item>
        <el-button type="primary" style="width:100%" :loading="loading" @click="onSubmit">{{ $t('login.submit') }}</el-button>
      </el-form>
    </el-card>
  </div>
</template>

<style scoped>
.login-wrap { min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f0f2f5; }
.login-card { width: 360px; }
.login-top { display: flex; justify-content: flex-end; margin-bottom: 6px; }
.login-title { text-align: center; margin: 0 0 18px; font-size: 20px; }
</style>
