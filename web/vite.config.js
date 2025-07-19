import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import vitePluginImp from 'vite-plugin-imp'

// https://vite.dev/config/
export default defineConfig(()=>{
  return {
    base: '/sgridnext/',
    plugins: [
      react(),
      vitePluginImp({
        libList: [
          {
            libName: 'antd',
            style: (name) => `antd/es/${name}/style`,
          },
        ],
      }),
    ],
    build: {
      rollupOptions: {
        external: ['lodash', /^lodash\//]
      }
    },
    css: {
      preprocessorOptions: {
        less: {
          javascriptEnabled: true,
        },
      },
    },
    server: {
      proxy: {
        '/api/': {
          target: 'http://10.169.114.23:15872/api/',
          changeOrigin: true,
          rewrite: (path) => path.replace(/^\/api/, ''),
        },
      },
    },
  }
})
