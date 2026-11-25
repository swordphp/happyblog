import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'

export const useSiteStore = defineStore('site', () => {
  const sitename = ref('HappyBlog')
  const sitedesc = ref('')
  const github = ref('')
  const loaded = ref(false)

  async function load() {
    if (loaded.value) return
    try {
      const data = await api.siteInfo()
      if (data) {
        sitename.value = data.sitename || 'HappyBlog'
        sitedesc.value = data.sitedesc || ''
        github.value = data.github || ''
        document.title = sitename.value
      }
    } catch (e) { /* 忽略,用默认值 */ }
    loaded.value = true
  }

  return { sitename, sitedesc, github, loaded, load }
})
