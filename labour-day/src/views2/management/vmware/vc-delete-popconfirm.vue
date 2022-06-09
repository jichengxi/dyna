<template>
  <n-popconfirm
      :on-positive-click="handlePositiveClick"
      :on-negative-click="handleNegativeClick"
      positive-text="确认"
      negative-text="取消"
  >
    <template #trigger>
      <n-button>删除</n-button>
    </template>
    是否需要删除名为 {{ name }} 的vCenter Cluster
  </n-popconfirm>
</template>

<script lang="ts" setup>

import { defineProps } from 'vue'
import { useMessage } from 'naive-ui'
const props = defineProps<{
  id: number
  name: string
}>()

const emit = defineEmits<{
  (e: 'delCluster', num: number): void
}>()
const message = useMessage()
const handlePositiveClick = () => {
  emit('delCluster', props.id)
  // message.info('是的')
}
const handleNegativeClick = () => {
  message.info('并不')
}
</script>
