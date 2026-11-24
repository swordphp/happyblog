import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  base: '/',
  plugins: [vue()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  build: {
    outDir: 'dist',      // 产物落在 web/dist,由 Go 二进制内嵌
    assetsDir: 'spa',    // 静态资源放 /spa,避开读者上传目录 /assets
    emptyOutDir: true,
  },
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://127.0.0.1:8888',
      '/admin/api': 'http://127.0.0.1:8888',
      '/assets': 'http://127.0.0.1:8888',
    },
  },
})
