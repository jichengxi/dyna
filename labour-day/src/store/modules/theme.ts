import { defineStore } from 'pinia'
import { themeSettings } from '@/settings'
export const useThemeStore = defineStore('theme', {
  state() {
    return themeSettings
  },
  getters: {},
  actions: {
    setTabVisible(visible: boolean) {
      this.tags.visible = visible
    },
  },
})
