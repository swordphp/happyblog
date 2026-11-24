<script setup>
import { ref, onMounted } from 'vue'
import api from '@/api'

const albums = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    albums.value = (await api.pubAlbums(1)) || []
  } catch (e) {
    albums.value = []
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <router-link to="/" class="back-link">&laquo; 返回首页</router-link>
  <h2 style="margin:6px 0 18px">专辑</h2>

  <div v-if="loading" class="empty">加载中…</div>
  <template v-else>
    <article v-for="al in albums" :key="al.id" class="post-card">
      <h2><router-link :to="`/album/${al.id}`">{{ al.albumName }}</router-link></h2>
      <div class="post-meta">{{ al.articleTotal }} 篇文章</div>
    </article>
    <p v-if="!albums.length" class="empty">还没有公开的专辑。</p>
  </template>
</template>
