import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 8081,
    proxy: {
      '/api': {
        target: 'http://localhost:8083',
        changeOrigin: true
      },
      '/uploads': {
        target: 'http://localhost:8083',
        changeOrigin: true
      }
    }
  }
})
