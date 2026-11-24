<script setup>
import { useI18n } from 'vue-i18n'
import { setLocale } from '@/i18n'

const { locale } = useI18n()
</script>

<template>
  <select
    class="lang-switch"
    :value="locale"
    @change="(e) => setLocale(e.target.value)"
    :aria-label="$t('lang.label')"
  >
    <option value="zh">中文</option>
    <option value="en">English</option>
  </select>
</template>

<style scoped>
.lang-switch {
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  padding: 4px 8px;
  font-size: 13px;
  color: #555;
  background: #fff;
  cursor: pointer;
}
.lang-switch:focus { outline: none; border-color: #0d6efd; }
</style>
