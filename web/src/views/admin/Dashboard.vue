<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'

const router = useRouter()
const stats = ref({ articles: 0, tags: 0, albums: 0, users: 0 })
const loading = ref(true)

onMounted(async () => {
  try {
    const [arts, tags, albums, users] = await Promise.all([
      api.adminArticles({ page: 1 }),
      api.adminTags(),
      api.adminAlbums(),
      api.adminUsers(),
    ])
    stats.value = {
      articles: (arts && arts.total) || 0,
      tags: (tags && tags.length) || 0,
      albums: (albums && albums.length) || 0,
      users: (users && users.length) || 0,
    }
  } catch (e) { /* 忽略 */ } finally {
    loading.value = false
  }
})

const cards = [
  { key: 'articles', i18n: 'dashboard.articles', to: '/admin/articles' },
  { key: 'tags', i18n: 'dashboard.tags', to: '/admin/tags' },
  { key: 'albums', i18n: 'dashboard.albums', to: '/admin/albums' },
  { key: 'users', i18n: 'dashboard.users', to: '/admin/users' },
]
</script>

<template>
  <h2 class="page-title">{{ $t('dashboard.title') }}</h2>
  <el-row :gutter="16" v-loading="loading">
    <el-col :span="6" v-for="c in cards" :key="c.key">
      <el-card shadow="hover" class="stat" @click="router.push(c.to)">
        <div class="num">{{ stats[c.key] }}</div>
        <div class="label">{{ $t(c.i18n) }}</div>
      </el-card>
    </el-col>
  </el-row>
</template>

<style scoped>
.page-title { margin: 0 0 18px; font-size: 20px; }
.stat { cursor: pointer; text-align: center; }
.stat .num { font-size: 30px; font-weight: 700; color: #0d6efd; }
.stat .label { color: #888; margin-top: 6px; }
</style>
