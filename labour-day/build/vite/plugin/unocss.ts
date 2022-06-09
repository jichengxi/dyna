import Unocss from 'unocss/vite'
import {presetUno, presetAttributify, presetIcons} from 'unocss'

// https://github.com/antfu/unocss
export function configUnocssPlugin() {
    return Unocss({
        presets: [presetUno(), presetAttributify(), presetIcons()],
        rules: unocssRules
    })
}

const unocssRules: any = [
    ['f-r', {float: 'right'}]
]

