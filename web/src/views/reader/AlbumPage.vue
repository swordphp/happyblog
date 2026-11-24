<script setup>
import { ref, watch, onMounted } from 'vue'
import api from '@/api'
import { fmtDate } from '@/utils/format'

const props = defineProps({ id: { type: [String, Number], required: true } })
const album = ref(null)
const articles = ref([])
const loading = ref(true)
const notFound = ref(false)

async function load(id) {
  loading.value = true
  notFound.value = false
  try {
    const data = await api.pubAlbum(id)
    album.value = data.album
    articles.value = data.articles || []
    if (album.value && album.value.albumName) document.title = album.value.albumName
  } catch (e) {
    notFound.value = true
  } finally {
    loading.value = false
  }
}

watch(() => props.id, (id) => load(id))
onMounted(() => load(props.id))
</script>

<template>
  <router-link to="/albums" class="back-link">&laquo; 全部专辑</router-link>

  <div v-if="loading" class="empty">加载中…</div>
  <div v-else-if="notFound" class="empty">专辑不存在或未公开。</div>

  <template v-else>
    <h2 style="margin:6px 0 18px">{{ album.albumName }}</h2>
    <article v-for="a in articles" :key="a.id" class="post-card">
      <h2><router-link :to="`/article/${a.id}`">{{ a.title }}</router-link></h2>
      <div class="post-meta">{{ a.authorInfo && a.authorInfo.nickName }} · {{ fmtDate(a.updateTime) }}</div>
      <div v-if="a.brief" class="post-brief">{{ a.brief }}</div>
    </article>
    <p v-if="!articles.length" class="empty">该专辑下还没有已发布文章。</p>
  </template>
</template>
