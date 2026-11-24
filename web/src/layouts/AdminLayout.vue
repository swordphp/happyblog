<script setup>
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import LangSwitch from '@/components/LangSwitch.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const { t } = useI18n()

const activeMenu = computed(() => route.path)
const nickName = computed(() => (auth.user && auth.user.nickName) || t('admin.menu.dashboard'))

async function onLogout() {
  await auth.logout()
  ElMessage.success(t('admin.loggedOut'))
  router.push({ name: 'admin-login' })
}
</script>

<template>
  <el-container class="admin-root">
    <el-aside width="210px" class="admin-aside">
      <div class="brand">HappyBlog</div>
      <el-menu :default-active="activeMenu" router unique-opened>
        <el-menu-item index="/admin"><span>{{ $t('admin.menu.dashboard') }}</span></el-menu-item>
        <el-menu-item index="/admin/articles"><span>{{ $t('admin.menu.articles') }}</span></el-menu-item>
        <el-menu-item index="/admin/tags"><span>{{ $t('admin.menu.tags') }}</span></el-menu-item>
        <el-menu-item index="/admin/albums"><span>{{ $t('admin.menu.albums') }}</span></el-menu-item>
        <el-menu-item index="/admin/settings"><span>{{ $t('admin.menu.settings') }}</span></el-menu-item>
        <el-menu-item index="/admin/users"><span>{{ $t('admin.menu.users') }}</span></el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="admin-header">
        <el-button text @click="router.push({ name: 'admin-article-new' })">{{ $t('admin.writeArticle') }}</el-button>
        <div class="spacer" />
        <LangSwitch />
        <a href="/" target="_blank" class="view-site">{{ $t('admin.viewSite') }}</a>
        <el-dropdown @command="(c) => c === 'logout' && onLogout()">
          <span class="user">{{ nickName }} <el-icon class="el-icon--right"><arrow-down /></el-icon></span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">{{ $t('admin.logout') }}</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>

      <el-main class="admin-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.admin-root { height: 100vh; }
.admin-aside { background: #001529; color: #fff; }
.brand { height: 60px; line-height: 60px; text-align: center; font-size: 18px; font-weight: 700; color: #fff; letter-spacing: 1px; }
.admin-aside :deep(.el-menu) { border-right: none; background: #001529; }
.admin-aside :deep(.el-menu-item) { color: #c9d1d9; }
.admin-aside :deep(.el-menu-item.is-active) { color: #fff; background: #0d6efd; }
.admin-header { display: flex; align-items: center; gap: 14px; background: #fff; border-bottom: 1px solid #eee; }
.spacer { flex: 1; }
.view-site { color: #666; font-size: 14px; }
.user { cursor: pointer; color: #333; display: flex; align-items: center; }
.admin-main { background: #f0f2f5; }
</style>
