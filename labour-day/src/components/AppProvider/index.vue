<template>
<!--  设置主题-->
  <n-config-provider :theme-overrides="useTheme.naiveThemeOverrides">
    <n-loading-bar-provider>
      <LoadingBar />
      <n-dialog-provider>
        <DialogContent />
        <n-message-provider>
          <MessageContent />
          <slot></slot>
        </n-message-provider>
      </n-dialog-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script lang="ts">
import MessageContent from './MessageContent.vue'
import DialogContent from './DialogContent.vue'
import LoadingBar from './LoadingBar.vue'
import { useThemeStore } from '@/store/modules/theme'
import { defineComponent, ref } from 'vue'
export default defineComponent({
  component: {MessageContent, DialogContent, LoadingBar},
  setup(){
    const useTheme = ref(useThemeStore())
    return {
      useTheme
    }
  }
})


</script>
