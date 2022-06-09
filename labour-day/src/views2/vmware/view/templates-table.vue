<template>
  <n-data-table
      size="small"
      :columns="vcStore.vmTemplatesTableColumns"
      :data="vcStore.vmTemplatesTableData"
      :single-line="false"
      :pagination="vmsPagination"
      :loading="vcStore.tableLoading"
  />
</template>

<script setup>
import {ref} from 'vue'
import {useVcStore} from "@/store/modules/vc/vc-manage"
const vcStore = useVcStore()

const vmsPagination = ref({pageSize: 10})

</script>
