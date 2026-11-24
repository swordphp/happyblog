<script setup>
import { onMounted } from 'vue'
import { useSiteStore } from '@/stores/site'
import LangSwitch from '@/components/LangSwitch.vue'

const site = useSiteStore()
onMounted(() => site.load())
</script>

<template>
  <div class="reader-layout">
    <header class="site-header">
      <div class="container">
        <router-link class="title" to="/">{{ site.sitename }}</router-link>
        <div v-if="site.sitedesc" class="desc">{{ site.sitedesc }}</div>
        <nav class="site-nav">
          <router-link to="/">{{ $t('nav.home') }}</router-link>
          <router-link to="/albums">{{ $t('nav.albums') }}</router-link>
          <a href="/api/rss" target="_blank" rel="noopener">{{ $t('nav.rss') }}</a>
          <LangSwitch class="nav-lang" />
        </nav>
      </div>
    </header>

    <main class="main">
      <div class="container">
        <router-view />
      </div>
    </main>

    <footer class="site-footer">
      {{ $t('reader.poweredBy') }} · <a href="/api/rss">RSS</a>
    </footer>
  </div>
</template>

<style scoped>
.nav-lang { margin-left: auto; }
</style>
