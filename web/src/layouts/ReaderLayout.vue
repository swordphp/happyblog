<script setup>
import { onMounted } from 'vue'
import { useSiteStore } from '@/stores/site'

const site = useSiteStore()
onMounted(() => site.load())
</script>

<template>
  <header class="site-header">
    <div class="container">
      <router-link class="title" to="/">{{ site.sitename }}</router-link>
      <div v-if="site.sitedesc" class="desc">{{ site.sitedesc }}</div>
      <nav class="site-nav">
        <router-link to="/">首页</router-link>
        <router-link to="/albums">专辑</router-link>
        <a href="/api/rss" target="_blank" rel="noopener">RSS</a>
      </nav>
    </div>
  </header>

  <main class="main">
    <div class="container">
      <router-view />
    </div>
  </main>

  <footer class="site-footer">
    Powered by HappyBlog · <a href="/api/rss">RSS</a>
  </footer>
</template>
