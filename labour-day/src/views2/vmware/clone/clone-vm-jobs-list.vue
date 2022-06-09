<template>
  <div p15>
    <n-space>
      <n-button type="info" @click="handleStart(checkedRowKeys)"> 开始任务 </n-button>
    </n-space>
    <n-space vertical :size="12">
      <n-data-table
          :bordered="false"
          :columns="cloneVmJobColumns"
          :data="cloneVmStore.cloneVmJobList"
          :pagination="cloneVmJobPagination"
          :loading="tableLoading"
          :row-key="rowKey"
          v-model:checked-row-keys="checkedRowKeys"
      />
    </n-space>
  </div>
</template>

<script setup>
import {ref, onMounted, h} from 'vue'
import {getCloneVmJobs, postCloneVmRun} from './clone-methods'
import Socket from '@/utils/websocket'
import {useCloneVmStore} from "@/store/modules/vmware/clone-vm"
import {NButton} from 'naive-ui'

const cloneVmStore = useCloneVmStore()
const checkedRowKeys = ref([])
const rowKey = (row) => row.id

// 表格菜单
const cloneVmJobColumns = [
  {
    type: 'selection',
  },
  {
    title: 'Job Name',
    key: 'name',
  },
  {
    title: 'Vc Name',
    key: 'vc_name',
  },
  {
    title: 'Vm Name',
    key: 'vm_name',
  },
  {
    title: 'Template',
    key: 'template',
  },
  {
    title: 'Host',
    key: 'host',
  },
  {
    title: 'Datastore',
    key: 'datastore',
  },
  {
    title: 'Folder',
    key: 'folder',
  },
  {
    title: 'CPU',
    key: 'cpu_num',
  },
  {
    title: 'Memory',
    key: 'memory',
  },
  {
    title: 'Status',
    key: 'status',
  },
  {
    title: 'Message',
    key: 'message',
  },
  {
    title: 'Create Time',
    key: 'create_time',
  },
  {
    title: 'Update Time',
    key: 'update_time',
  },
  {
    title: 'Action',
    key: 'actions',
    render (row) {
      return h(
          NButton,
          {
            size: 'small',
            type: "info",
            onClick: ()=> {
              postCloneVmRun([row.id])
            }
          },
          { default: () => '开始' }
      )
    }
  }
]
// const cloneVmJobList = ref([])

// 表格换页
const cloneVmJobPagination = ref({pageSize: 10})
const tableLoading = ref(false)

const handleStart = (rowKeys) => {
  console.log("rowKeys", rowKeys)
  postCloneVmRun(rowKeys)
}


const ws = ref(null)

const initWebsocket = () => {
  ws.value = new Socket({url: "/vmware/cluster/clone-vm/message"})
}
initWebsocket()

onMounted(async () => {
  tableLoading.value = true
  const vmJobs = await getCloneVmJobs()
  if (vmJobs !== null) {
    cloneVmStore.cloneVmJobList = vmJobs
  }
  tableLoading.value = false
  ws.value.onmessage((data) => {
    cloneVmStore.cloneVmJobList.some((job) => {
      if (job.id === data.id) {
        job.status = data.status
        job.message = data.message
        job.update_time = data.update_time
        return true
      }
    })
  })
})
</script>

