import { createRouter, createWebHistory } from 'vue-router'
import api from '@/api'
import { useAuthStore } from '@/stores/auth'

import ReaderLayout from '@/layouts/ReaderLayout.vue'
import Home from '@/views/reader/Home.vue'
import Article from '@/views/reader/Article.vue'
import TagPage from '@/views/reader/TagPage.vue'
import Albums from '@/views/reader/Albums.vue'
import AlbumPage from '@/views/reader/AlbumPage.vue'
import Install from '@/views/Install.vue'

const routes = [
  {
    path: '/',
    component: ReaderLayout,
    children: [
      { path: '', name: 'home', component: Home },
      { path: 'article/:id', name: 'reader-article', component: Article, props: true },
      { path: 'tag/:id', name: 'reader-tag', component: TagPage, props: true },
      { path: 'albums', name: 'reader-albums', component: Albums },
      { path: 'album/:id', name: 'reader-album', component: AlbumPage, props: true },
    ],
  },
  { path: '/install', name: 'install', component: Install },

  // 管理端:异步分包,Element Plus 只打进这些 chunk
  { path: '/admin/login', name: 'admin-login', component: () => import('@/views/admin/Login.vue') },
  {
    path: '/admin',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'admin-dashboard', component: () => import('@/views/admin/Dashboard.vue') },
      { path: 'articles', name: 'admin-articles', component: () => import('@/views/admin/Articles.vue') },
      { path: 'article/new', name: 'admin-article-new', component: () => import('@/views/admin/ArticleEdit.vue') },
      { path: 'article/:id', name: 'admin-article-edit', component: () => import('@/views/admin/ArticleEdit.vue'), props: true },
      { path: 'tags', name: 'admin-tags', component: () => import('@/views/admin/Tags.vue') },
      { path: 'albums', name: 'admin-albums', component: () => import('@/views/admin/AlbumsAdmin.vue') },
      { path: 'settings', name: 'admin-settings', component: () => import('@/views/admin/Settings.vue') },
      { path: 'users', name: 'admin-users', component: () => import('@/views/admin/Users.vue') },
    ],
  },
]

const router = createRouter({
  history: createWebHistory('/'),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

// 安装状态缓存:只需确认一次「已安装」。
let installedCache = false
async function ensureInstalled() {
  if (installedCache) return true
  try {
    const data = await api.installStatus()
    installedCache = !!(data && data.installed)
  } catch (e) {
    installedCache = false
  }
  return installedCache
}

router.beforeEach(async (to) => {
  // 安装向导本身不做拦截
  if (to.name === 'install') return true

  const installed = await ensureInstalled()
  if (!installed) return { name: 'install' }

  if (to.meta.requiresAuth) {
    const auth = useAuthStore()
    if (!auth.checked) await auth.fetchMe()
    if (!auth.isLoggedIn()) {
      return { name: 'admin-login', query: { redirect: to.fullPath } }
    }
  }
  return true
})

export default router
