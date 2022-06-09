import {createHtmlPlugin} from 'vite-plugin-html'
// @ts-ignore
import {version} from '../../../package.json'
import {GLOB_CONFIG_FILE_NAME} from '../../constant'

export function configHtmlPlugin(viteEnv: Record<string, string>, isBuild: boolean) {
  const { VITE_APP_TITLE, VITE_PUBLIC_PATH } = viteEnv
  const path = VITE_PUBLIC_PATH.endsWith('/') ? VITE_PUBLIC_PATH : `${VITE_PUBLIC_PATH}/`

  const getAppConfigSrc = () => {
    return `${path}${GLOB_CONFIG_FILE_NAME}?v=${version}-${new Date().getTime()}`
  }

  return createHtmlPlugin({
    minify: isBuild,
    inject: {
      data: {
        title: VITE_APP_TITLE,
      },
      tags: isBuild
          ? [
            {
              tag: 'script',
              attrs: {
                src: getAppConfigSrc(),
              },
            },
          ]
          : [],
    },
  })
}
