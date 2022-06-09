<template>
  <AppProvider>
    <router-view v-slot="{ Component }">
      <component :is="Component" />
    </router-view>
  </AppProvider>
</template>

<script lang="ts">
import AppProvider from '@/components/AppProvider/index.vue'
import { defineComponent } from 'vue'
export default defineComponent ({
  name: "app",
  component: {AppProvider}
})
</script>

<style lang="scss">
#app {
  height: 100%;
  .n-config-provider {
    height: inherit;
  }
}
</style>
