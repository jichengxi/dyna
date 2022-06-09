<template>
  <n-data-table
      size="small"
      :data="vcStore.vmHostsTableData"
      :columns="vcStore.vmHostsTableColumns"
      :single-line="false"
      :pagination="vmHostsPagination"
      :loading="vcStore.tableLoading"
  />
</template>

<script setup>
import {ref} from 'vue'
import {useMessage} from 'naive-ui'
import {useVcStore} from "@/store/modules/vc/vc-manage"
const message = useMessage()
const vcStore = useVcStore()

const vmHostsPagination = ref({pageSize: 10})

</script>
