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
  <div v-if="loading" class="empty">加载中…</div>
  <template v-else>
    <article v-for="a in articles" :key="a.id" class="post-card">
      <h2><router-link :to="`/article/${a.id}`">{{ a.title }}</router-link></h2>
      <div class="post-meta">
        {{ a.authorInfo && a.authorInfo.nickName }} · {{ fmtDate(a.updateTime) }}
      </div>
      <div v-if="a.brief" class="post-brief">{{ a.brief }}</div>
    </article>

    <p v-if="!articles.length" class="empty">还没有已发布的文章。</p>

    <nav v-if="page > 1 || articles.length === pageSize" class="pager">
      <button class="btn" :disabled="page <= 1" @click="goPage(page - 1)">&laquo; 上一页</button>
      <span class="muted">第 {{ page }} 页</span>
      <button class="btn" :disabled="articles.length < pageSize" @click="goPage(page + 1)">下一页 &raquo;</button>
    </nav>
  </template>
</template>
