<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import LangSwitch from '@/components/LangSwitch.vue'

const router = useRouter()
const { t } = useI18n()
const submitting = ref(false)
const error = ref('')
const done = ref(false)

const f = reactive({
  sitename: '',
  siteurl: '',
  sitedesc: '',
  dbDriver: 'sqlite',
  sqlitePath: './happyblog.db',
  dbHost: '127.0.0.1',
  dbPort: '3306',
  dbName: '',
  dbUser: '',
  dbPassword: '',
  adminEmail: '',
  adminPassword: '',
  adminNickname: 'admin',
})

function validate() {
  if (!f.sitename.trim()) return t('install.errSiteName')
  if (!f.adminEmail.trim() || !f.adminPassword.trim()) return t('install.errAdmin')
  if ((f.dbDriver === 'mysql' || f.dbDriver === 'mariadb') && (!f.dbHost || !f.dbName || !f.dbUser)) {
    return t('install.errMysql')
  }
  return ''
}

async function submit() {
  error.value = ''
  const msg = validate()
  if (msg) { error.value = msg; return }
  submitting.value = true
  try {
    await api.install({ ...f })
    done.value = true
    setTimeout(() => router.push({ name: 'admin-login' }), 1200)
  } catch (e) {
    error.value = e.message || t('install.failed')
  } finally {
    submitting.value = false
  }
}
</script>
<template>
  <div class="install-wrap">
    <div class="install-card">
      <div class="install-top"><LangSwitch /></div>
      <h1>{{ $t('install.title') }}</h1>
      <p class="muted">{{ $t('install.intro') }}</p>

      <div v-if="done" class="ok">{{ $t('install.done') }}</div>

      <template v-else>
        <h3>{{ $t('install.siteSection') }}</h3>
        <label>{{ $t('install.siteName') }} *<input v-model="f.sitename" :placeholder="$t('install.siteNamePh')" /></label>
        <label>{{ $t('install.siteUrl') }}<input v-model="f.siteurl" placeholder="https://example.com" /></label>
        <label>{{ $t('install.siteDesc') }}<input v-model="f.sitedesc" :placeholder="$t('install.siteDescPh')" /></label>

        <h3>{{ $t('install.dbSection') }}</h3>
        <label>{{ $t('install.dbType') }}
          <select v-model="f.dbDriver">
            <option value="sqlite">{{ $t('install.dbSqlite') }}</option>
            <option value="mysql">{{ $t('install.dbMysql') }}</option>
          </select>
        </label>
        <template v-if="f.dbDriver === 'sqlite'">
          <label>{{ $t('install.sqlitePath') }}<input v-model="f.sqlitePath" /></label>
        </template>
        <template v-else>
          <label>{{ $t('install.dbHost') }}<input v-model="f.dbHost" /></label>
          <label>{{ $t('install.dbPort') }}<input v-model="f.dbPort" /></label>
          <label>{{ $t('install.dbName') }}<input v-model="f.dbName" /></label>
          <label>{{ $t('install.dbUser') }}<input v-model="f.dbUser" /></label>
          <label>{{ $t('install.dbPassword') }}<input v-model="f.dbPassword" type="password" /></label>
        </template>

        <h3>{{ $t('install.adminSection') }}</h3>
        <label>{{ $t('install.email') }} *<input v-model="f.adminEmail" type="email" placeholder="admin@example.com" /></label>
        <label>{{ $t('install.password') }} *<input v-model="f.adminPassword" type="password" /></label>
        <label>{{ $t('install.nickname') }}<input v-model="f.adminNickname" /></label>

        <p v-if="error" class="err">{{ error }}</p>
        <button class="btn primary" :disabled="submitting" @click="submit">
          {{ submitting ? $t('install.submitting') : $t('install.submit') }}
        </button>
      </template>
    </div>
  </div>
</template>

<style scoped>
.install-wrap { min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 40px 16px; }
.install-card { width: 100%; max-width: 480px; background: #fff; border: 1px solid #ececec; border-radius: 12px; padding: 32px 36px; box-shadow: 0 6px 30px rgba(0,0,0,.05); }
.install-top { display: flex; justify-content: flex-end; }
.install-card h1 { font-size: 22px; margin: 6px 0 6px; }
.install-card h3 { margin: 22px 0 8px; font-size: 15px; color: #333; border-bottom: 1px solid #f0f0f0; padding-bottom: 6px; }
label { display: block; font-size: 13px; color: #555; margin-bottom: 12px; }
input, select { display: block; width: 100%; margin-top: 4px; padding: 9px 11px; border: 1px solid #d9d9d9; border-radius: 6px; font-size: 14px; }
input:focus, select:focus { outline: none; border-color: #0d6efd; }
.btn.primary { width: 100%; background: #0d6efd; color: #fff; border: none; padding: 11px; border-radius: 6px; font-size: 15px; margin-top: 10px; cursor: pointer; }
.btn.primary[disabled] { opacity: .6; cursor: default; }
.err { color: #d93025; font-size: 13px; }
.ok { color: #1a7f37; padding: 20px 0; text-align: center; font-size: 15px; }
</style>