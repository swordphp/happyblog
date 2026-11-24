<script setup>
import { ref, watch, onMounted } from 'vue'
import api from '@/api'
import { fmtDate } from '@/utils/format'

const props = defineProps({ id: { type: [String, Number], required: true } })
const articles = ref([])
const tagName = ref('')
const loading = ref(true)

async function load(id) {
  loading.value = true
  try {
    const [list, tags] = await Promise.all([api.pubTagArticles(id), api.pubTags()])
    articles.value = list || []
    const t = (tags || []).find((x) => String(x.id) === String(id))
    tagName.value = t ? t.tagName : ''
  } catch (e) {
    articles.value = []
  } finally {
    loading.value = false
  }
}

watch(() => props.id, (id) => load(id))
onMounted(() => load(props.id))
</script>

<template>
  <router-link to="/" class="back-link">&laquo; 返回首页</router-link>
  <h2 style="margin:6px 0 18px">标签:{{ tagName || '#' + id }}</h2>

  <div v-if="loading" class="empty">加载中…</div>
  <template v-else>
    <article v-for="a in articles" :key="a.id" class="post-card">
      <h2><router-link :to="`/article/${a.id}`">{{ a.title }}</router-link></h2>
      <div class="post-meta">{{ a.authorInfo && a.authorInfo.nickName }} · {{ fmtDate(a.updateTime) }}</div>
      <div v-if="a.brief" class="post-brief">{{ a.brief }}</div>
    </article>
    <p v-if="!articles.length" class="empty">该标签下还没有文章。</p>
  </template>
</template>
