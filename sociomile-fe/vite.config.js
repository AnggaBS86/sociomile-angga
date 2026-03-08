import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/health': 'http://127.0.0.1:8080',
      '/auth': 'http://127.0.0.1:8080',
      '/channel': 'http://127.0.0.1:8080',
      '/conversations': 'http://127.0.0.1:8080',
      '/tickets': 'http://127.0.0.1:8080'
    }
  }
})
