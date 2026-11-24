<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { marked } from 'marked'
import api from '@/api'
import { fmtDate } from '@/utils/format'

const props = defineProps({ id: { type: [String, Number], required: true } })

const article = ref(null)
const tags = ref([])
const album = ref(null)
const loading = ref(true)
const notFound = ref(false)

const contentHtml = computed(() =>
  article.value ? marked.parse(article.value.content || '') : ''
)

async function load(id) {
  loading.value = true
  notFound.value = false
  try {
    const data = await api.pubArticleInfo(id)
    article.value = data.articleInfo
    tags.value = data.tags || []
    album.value = data.album && data.album.id ? data.album : null
    if (article.value && article.value.title) document.title = article.value.title
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
  <router-link to="/" class="back-link">&laquo; 返回首页</router-link>

  <div v-if="loading" class="empty">加载中…</div>
  <div v-else-if="notFound" class="empty">文章不存在,或尚未发布。</div>

  <article v-else class="article">
    <h1 class="article-title">{{ article.title }}</h1>
    <div class="post-meta">
      {{ article.authorInfo && article.authorInfo.nickName }} · {{ fmtDate(article.createTime, true) }}
    </div>
    <div v-if="tags.length" style="margin:10px 0 4px">
      <router-link v-for="t in tags" :key="t.id" class="tag-badge" :to="`/tag/${t.id}`">{{ t.tagName }}</router-link>
    </div>
    <div v-if="album" class="muted" style="margin-bottom:10px;font-size:13px">
      专辑:<router-link :to="`/album/${album.id}`">{{ album.albumName }}</router-link>
    </div>
    <hr />
    <div class="article-content" v-html="contentHtml"></div>
  </article>
</template>
