<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { fmtDate } from '@/utils/format'

const route = useRoute()
const router = useRouter()
const articles = ref([])
const page = ref(1)
const loading = ref(false)
const pageSize = 10

async function load(p) {
  loading.value = true
  try {
    page.value = p
    articles.value = (await api.pubArticleList(p)) || []
  } catch (e) {
    articles.value = []
  } finally {
    loading.value = false
  }
}

function goPage(p) {
  if (p < 1) return
  router.push({ name: 'home', query: { page: p } })
}

watch(() => route.query.page, (q) => load(Number(q) || 1))
onMounted(() => load(Number(route.query.page) || 1))
</script>

<template>
  <div v-if="loading" class="empty">{{ $t('common.loading') }}</div>
  <template v-else>
    <router-link v-for="a in articles" :key="a.id" class="post-card" :to="`/article/${a.id}`">
      <h2>{{ a.title }}</h2>
      <div class="post-meta">
        {{ a.authorInfo && a.authorInfo.nickName }} · {{ fmtDate(a.updateTime) }}
      </div>
      <div v-if="a.brief" class="post-brief">{{ a.brief }}</div>
    </router-link>

    <p v-if="!articles.length" class="empty">{{ $t('reader.noArticles') }}</p>

    <nav v-if="page > 1 || articles.length === pageSize" class="pager">
      <button class="btn" :disabled="page <= 1" @click="goPage(page - 1)">{{ $t('reader.prevPage') }}</button>
      <span class="muted">{{ $t('reader.pageN', { page }) }}</span>
      <button class="btn" :disabled="articles.length < pageSize" @click="goPage(page + 1)">{{ $t('reader.nextPage') }}</button>
    </nav>
  </template>
</template>
