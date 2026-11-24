import { createI18n } from 'vue-i18n'
import zh from './zh'
import en from './en'

const saved = typeof localStorage !== 'undefined' ? localStorage.getItem('locale') : null
const browser = (navigator.language || '').toLowerCase().startsWith('en') ? 'en' : 'zh'
export const initialLocale = saved || browser

const i18n = createI18n({
  legacy: false,      // Composition API 模式,组件里用 useI18n()
  globalInjection: true, // 模板里可直接用 $t
  locale: initialLocale,
  fallbackLocale: 'zh',
  messages: { zh, en },
})

// 切换语言:更新 vue-i18n、持久化、同步 <html lang>
export function setLocale(l) {
  i18n.global.locale.value = l
  try { localStorage.setItem('locale', l) } catch (e) { /* ignore */ }
  document.documentElement.setAttribute('lang', l === 'en' ? 'en' : 'zh-CN')
}

export default i18n
