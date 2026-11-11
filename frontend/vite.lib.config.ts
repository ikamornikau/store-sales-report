import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

const external = ['vue', 'chart.js', 'vue-chartjs']

// Library build for embedding ReportScreen in another Vue 3 app.
// The SPA build stays on vite.config.ts and writes to dist for the nginx image.
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  publicDir: false,
  build: {
    outDir: 'dist-lib',
    emptyOutDir: true,
    copyPublicDir: false,
    sourcemap: true,
    lib: {
      entry: fileURLToPath(new URL('./src/report/index.ts', import.meta.url)),
      formats: ['es'],
      fileName: 'store-sales-report',
      cssFileName: 'style',
    },
    rollupOptions: {
      external: (id) => external.some((name) => id === name || id.startsWith(`${name}/`)),
    },
  },
})
