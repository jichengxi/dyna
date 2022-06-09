import {ConfigEnv, defineConfig, loadEnv} from 'vite'
import { resolve } from 'path'
import { createVitePlugins } from './build/vite/plugin'
import { createProxy } from './build/vite/proxy'
import { OUTPUT_DIR } from './build/constant'
import {wrapperEnv} from './build/utils'

// https://vitejs.dev/config/


// @ts-ignore
export default defineConfig(({ command, mode }) => {
  const root = process.cwd()
  const isBuild = command === 'build'
  const env = loadEnv(mode, process.cwd() + "/env")
  wrapperEnv(env)
  return {
    root,
    base: env['VITE_PUBLIC_PATH'] || '/',
    plugins: [createVitePlugins(env, isBuild)],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src'),
      },
    },
    css: {
      preprocessorOptions: {
        //define global scss variable
        scss: {
          additionalData: `@import '@/styles/variables.scss';`,
        },
      },
    },
    build: {
      target: 'es2015',
      outDir: OUTPUT_DIR,
      brotliSize: false,
      chunkSizeWarningLimit: 2000,
      assetsInlineLimit: 8192 // 小于此阈值的导入或引用资源将内联为 base64 编码，以避免额外的 http 请求。设置为 0 可以完全禁用此项。

    },
    server: {
      host: '127.0.0.1',
      strictPort: true,
      port: env['VITE_PORT'],
      // proxy: createProxy(env['VITE_PROXY']),
      hmr: {
        overlay: false
      }
    }
  }
})
