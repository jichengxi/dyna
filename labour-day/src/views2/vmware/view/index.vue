<template>
  <div p10>
      <div class="f-r">
        <n-dropdown trigger="click" :options="vcStore.dropDownVcList" @select="handleSelect">
          <n-button> {{ vcStore.currentVcCluster.label }}</n-button>
        </n-dropdown>
      </div>
      <n-tabs type="line" animated :on-before-leave="beforeSelect">
        <template #prefix>
          {{ vcStore.currentVcCluster.label }}
        </template>
        <n-tab-pane name="host" tab="主机">
          <hosts-table />
        </n-tab-pane>
        <n-tab-pane name="vm" tab="虚拟机">
          <vms-table />
        </n-tab-pane>
        <n-tab-pane name="template" tab="模板">
          <templates-table />
        </n-tab-pane>
      </n-tabs>
  </div>
</template>

<script setup>
import {useMessage} from 'naive-ui'
import {onBeforeMount, unref} from 'vue'
import {useRouter} from 'vue-router'
import HostsTable from './hosts-table.vue'
import VmsTable from './vms-table.vue'
import TemplatesTable from './templates-table.vue'
import {useVcStore} from "@/store/modules/vc/vc-manage"

const message = useMessage()
const vcStore = useVcStore()
const router = useRouter()
const components = [HostsTable, VmsTable, TemplatesTable]

// 选择切换集群
const handleSelect = (key) => {
  vcStore.dropDownVcList.map(async (vc) => {
    if (vc.key === key) {
      vcStore.currentVcCluster = vc
      await loginVc(vcStore.currentVcCluster.key)
    }
  })
}

// 获取vc列表
const getVcClustersList = async () => {
  await vcStore.getVcClusters()
}

// 登录vc
const loginVc = async (id) => {
  const res = await vcStore.loginVcCluster({data: {id}})
  if (res) {
    message.info("登录成功")
  } else {
    message.error("登录失败")
  }
}

onBeforeMount(async () => {
  await vcStore.changeTableLoading()
  await getVcClustersList()
  vcStore.currentVcCluster = vcStore.dropDownVcList[0]
  await loginVc(vcStore.currentVcCluster.key)
  await vcStore.getVmHosts()
  await vcStore.createVmHostsTableColumns(vcStore.vmHostsTableData, routeHostInfo)
  await vcStore.changeTableLoading()
})

// 路由跳转到主机详情
const routeHostInfo = (host) => {
  router.push({ name: 'VmwareHostInfo', query: { host: host }})
}

// 切换选择
const beforeSelect = async (name) => {
  if ((name === "vm" && vcStore.vmsTableData.length === 0) || (name === "template" && vcStore.vmTemplatesTableData.length === 0)) {
    console.log("beforeSelect")
    await vcStore.changeTableLoading()
    await vcStore.getVms()
    await vcStore.createVmsTableColumns()
    await vcStore.createVmTemplateTableColumns()
    await vcStore.changeTableLoading()
  }
  return true
}

</script>
