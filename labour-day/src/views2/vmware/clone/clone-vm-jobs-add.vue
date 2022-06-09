<template>
  <div p15>
    <clone-vm-from />
    <clone-vms-from />
  </div>
</template>

<script setup>
import CloneVmFrom from './clone-vm-from.vue'
import CloneVmsFrom from './clone-vms-table.vue'

const components = [CloneVmFrom, CloneVmsFrom]
</script>

