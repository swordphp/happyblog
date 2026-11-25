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
  <router-link to="/" class="back-link">{{ $t('reader.backHome') }}</router-link>
  <h2 style="margin:6px 0 18px">{{ $t('reader.tagTitle', { name: tagName || ('#' + id) }) }}</h2>

  <div v-if="loading" class="empty">{{ $t('common.loading') }}</div>
  <template v-else>
    <router-link v-for="a in articles" :key="a.id" class="post-card" :to="`/article/${a.id}`">
      <h2>{{ a.title }}</h2>
      <div class="post-meta">{{ a.authorInfo && a.authorInfo.nickName }} · {{ fmtDate(a.updateTime) }}</div>
      <div v-if="a.brief" class="post-brief">{{ a.brief }}</div>
    </router-link>
    <p v-if="!articles.length" class="empty">{{ $t('reader.noTagArticles') }}</p>
  </template>
</template>
