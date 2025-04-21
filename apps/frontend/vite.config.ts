import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// コンテナ内アクセス用
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  server: {
    host: true,
    port: 5173,
    watch: {
        usePolling: true,
    }
  }
})