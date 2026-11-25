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
  <router-link to="/albums" class="back-link">{{ $t('reader.allAlbums') }}</router-link>

  <div v-if="loading" class="empty">{{ $t('common.loading') }}</div>
  <div v-else-if="notFound" class="empty">{{ $t('reader.notFoundAlbum') }}</div>

  <template v-else>
    <h2 style="margin:6px 0 18px">{{ album.albumName }}</h2>
    <router-link v-for="a in articles" :key="a.id" class="post-card" :to="`/article/${a.id}`">
      <h2>{{ a.title }}</h2>
      <div class="post-meta">{{ a.authorInfo && a.authorInfo.nickName }} · {{ fmtDate(a.updateTime) }}</div>
      <div v-if="a.brief" class="post-brief">{{ a.brief }}</div>
    </router-link>
    <p v-if="!articles.length" class="empty">{{ $t('reader.noAlbumArticles') }}</p>
  </template>
</template>
