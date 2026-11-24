<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'

const router = useRouter()
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
  if (!f.sitename.trim()) return '请填写站点名称'
  if (!f.adminEmail.trim() || !f.adminPassword.trim()) return '请填写管理员邮箱与密码'
  if ((f.dbDriver === 'mysql' || f.dbDriver === 'mariadb') && (!f.dbHost || !f.dbName || !f.dbUser)) {
    return 'MySQL 连接信息不完整(主机/库名/用户)'
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
    error.value = e.message || '安装失败'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="install-wrap">
    <div class="install-card">
      <h1>HappyBlog 安装向导</h1>
      <p class="muted">填写以下信息完成初始化,无需重启即可使用。</p>

      <div v-if="done" class="ok">安装成功,正在跳转到登录页…</div>

      <template v-else>
        <h3>站点信息</h3>
        <label>站点名称 *<input v-model="f.sitename" placeholder="我的博客" /></label>
        <label>站点地址<input v-model="f.siteurl" placeholder="https://example.com" /></label>
        <label>站点描述<input v-model="f.sitedesc" placeholder="一句话介绍" /></label>

        <h3>数据库</h3>
        <label>类型
          <select v-model="f.dbDriver">
            <option value="sqlite">SQLite(推荐,免外置数据库)</option>
            <option value="mysql">MySQL / MariaDB</option>
          </select>
        </label>
        <template v-if="f.dbDriver === 'sqlite'">
          <label>数据文件路径<input v-model="f.sqlitePath" /></label>
        </template>
        <template v-else>
          <label>主机<input v-model="f.dbHost" /></label>
          <label>端口<input v-model="f.dbPort" /></label>
          <label>数据库名<input v-model="f.dbName" /></label>
          <label>用户名<input v-model="f.dbUser" /></label>
          <label>密码<input v-model="f.dbPassword" type="password" /></label>
        </template>

        <h3>管理员账号</h3>
        <label>邮箱 *<input v-model="f.adminEmail" type="email" placeholder="admin@example.com" /></label>
        <label>密码 *<input v-model="f.adminPassword" type="password" /></label>
        <label>昵称<input v-model="f.adminNickname" /></label>

        <p v-if="error" class="err">{{ error }}</p>
        <button class="btn primary" :disabled="submitting" @click="submit">
          {{ submitting ? '安装中…' : '开始安装' }}
        </button>
      </template>
    </div>
  </div>
</template>

<style scoped>
.install-wrap { min-height: 100vh; display: flex; align-items: center; justify-content: center; padding: 40px 16px; }
.install-card { width: 100%; max-width: 480px; background: #fff; border: 1px solid #ececec; border-radius: 12px; padding: 32px 36px; box-shadow: 0 6px 30px rgba(0,0,0,.05); }
.install-card h1 { font-size: 22px; margin: 0 0 6px; }
.install-card h3 { margin: 22px 0 8px; font-size: 15px; color: #333; border-bottom: 1px solid #f0f0f0; padding-bottom: 6px; }
label { display: block; font-size: 13px; color: #555; margin-bottom: 12px; }
input, select { display: block; width: 100%; margin-top: 4px; padding: 9px 11px; border: 1px solid #d9d9d9; border-radius: 6px; font-size: 14px; }
input:focus, select:focus { outline: none; border-color: #0d6efd; }
.btn.primary { width: 100%; background: #0d6efd; color: #fff; border: none; padding: 11px; border-radius: 6px; font-size: 15px; margin-top: 10px; cursor: pointer; }
.btn.primary[disabled] { opacity: .6; cursor: default; }
.err { color: #d93025; font-size: 13px; }
.ok { color: #1a7f37; padding: 20px 0; text-align: center; font-size: 15px; }
</style>
