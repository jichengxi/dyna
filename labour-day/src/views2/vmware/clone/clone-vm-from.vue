<template>
    <n-button ghost color="#8a2be2" mr-20px @click="handleCloneVmClick">
      克隆虚拟机
    </n-button>
    <n-modal v-model:show="showCloneModal" :mask-closable="false" preset="dialog" title="克隆虚拟机" @esc="handleModelCloseClick" @after-leave="handleModelCloseClick">
      <template #header>
        <div>克隆虚拟机</div>
      </template>
      <div>
        <n-spin size="medium" :show="spinLoading">
          <n-form ref="formRef" :model="cloneVmForm" :rules="cloneVmRules" label-placement="left" label-width="auto"
                  require-mark-placement="right">
            <n-form-item path="vc" label="vCenter">
              <n-select v-model:value="cloneVmForm.vc_name"
                        :loading="selectLoading"
                        placeholder="请选择vCenter集群"
                        :options="vcOptions"
                        remote
                        @update:show="clickSelect"
                        @update:value="updateSelectVc"
              />
            </n-form-item>

            <div v-show="cloneVmForm.vc_name.length !== 0">
              <n-form-item path="template" label="模版">
                <n-select v-model:value="cloneVmForm.template" placeholder="请选择模版" :options="templateOptions" @update:value="updateSelectTemplate" />
              </n-form-item>

              <n-form-item path="vm_name" label="虚拟机名称">
                <n-input v-model:value="cloneVmForm.vm_name" placeholder="请输入虚拟机名称" type="text"/>
              </n-form-item>

              <n-form-item path="host" label="目标主机">
                <n-select v-model:value="cloneVmForm.host" placeholder="请选择目标主机" :options="hostOptions" @update:value="updateSelectHost" />
              </n-form-item>
              <n-form-item path="host" label="目标存储">
                <n-select v-model:value="cloneVmForm.datastore" placeholder="请选择目标存储" :options="datastoreOptions" />
              </n-form-item>
              <n-form-item path="folder" label="目标文件夹">
                <n-select v-model:value="cloneVmForm.folder" placeholder="请选择目标文件夹" :options="folderOptions" />
              </n-form-item>

              <n-form-item path="cpu_num" label="CPU">
                <n-input-number v-model:value="cloneVmForm.cpu_num" placeholder="请输入CPU" clearable/>
              </n-form-item>
              <n-form-item path="memory" label="内存">
                <n-input-number v-model:value="cloneVmForm.memory" placeholder="内存" clearable/>
              </n-form-item>
              <n-form-item path="disk" label="硬盘">
                <div>
                  <n-input-number readonly v-for="disk in cloneVmForm.current_disk" :value="disk" :show-button="false"/>
                  <n-dynamic-input v-model:value="cloneVmForm.disks" :min="1" :max="5" :on-create="diskOnCreate">
                    <template #default="{ value }">
                      <div style="display: flex; align-items: center; width: 100%">
                        <n-input-number
                            v-model:value="value.disk_capacity_gb"
                            style="margin-right: 12px; width: 160px"
                        />
                        <n-switch v-model:value="value.disk_thin" />
                      </div>
                    </template>
                  </n-dynamic-input>
                </div>
              </n-form-item>
              <n-form-item path="network" label="网络">
                <n-dynamic-input v-model:value="cloneVmForm.networks" :min="1" :max="3" :on-create="networkOnCreate">
                  <template #default="{ value, index }">
                    <n-select v-model:value="value.network_name" placeholder="请选择网络" :options="networkOptions" @update:value="networkOnUpdate(index)"/>
                  </template>
                </n-dynamic-input>
              </n-form-item>
              <n-form-item path="ip_addr" label="业务网IP">
                <n-input v-model:value="cloneVmForm.ip.ip_addr" placeholder="请输入业务网IP" type="text"/>
              </n-form-item>
              <n-form-item path="net_mask" label="子网掩码">
                <n-input v-model:value="cloneVmForm.ip.net_mask" placeholder="请输入子网掩码" type="text"/>
              </n-form-item>
              <n-form-item path="gateway" label="业务网网关">
                <n-input v-model:value="cloneVmForm.ip.gateway" placeholder="请输入业务网网关" type="text"/>
              </n-form-item>
              <n-form-item path="dns" label="DNS">
                <n-select v-model:value="cloneVmForm.ip.dns" multiple :options="dnsOptions" />
              </n-form-item>
            </div>
          </n-form>
        </n-spin>
        <pre>{{ JSON.stringify(cloneVmForm, null, 2) }}</pre>
      </div>
      <template #action>
        <n-button
            :disabled="cloneVmForm.vm_name === null"
            round
            type="primary"
            @click="handleOkClick"
        >
          确定
        </n-button>
      </template>
    </n-modal>
</template>

<script setup>
import {ref} from 'vue'
import {useVcStore} from "@/store/modules/vc/vc-manage"
import {useCloneVmStore} from "@/store/modules/vmware/clone-vm"
import {postCloneVmJobs} from "./clone-methods"

const vcStore = useVcStore()
const cloneVmStore = useCloneVmStore()

const showCloneModal = ref(false)
const spinLoading = ref(false)

const vmNames = ref([])
const templates = ref([])
const hosts = ref([])
const hostNetworks = ref({})
const folders = ref([])


// 选择相关
const selectLoading = ref(false)  // 选择器加载
const vcOptions = ref([])
const hostOptions = ref([])
const templateOptions = ref([])
const folderOptions = ref([])
const datastoreOptions = ref([])
const networkOptions = ref([])
// const dnsOptions = ref([])

// 初始化cloneVmForm表单的数据
const initCloneVmForm = () => ({
  vc_name: "",
  template: "",
  vm_name: "",
  host: "",
  datastore: "",
  folder: "",
  cpu_num: 0,
  memory: 0,
  current_disk: [],
  disks: [],
  networks: [],
  ip: {
    ip_addr: "",
    net_mask: "",
    gateway: "",
    dns: []
  }
})
const formRef = ref(null)

// 选择vc点击时
const clickSelect = async (show) => {
  if (show && vcOptions.value.length === 0) {
    selectLoading.value = true
    await vcStore.getVcClusters()
    console.log(vcStore.vcClustersList)
    vcOptions.value = vcStore.vcClustersList.map((v) => ({
      label: v.name,
      value: v.name
    }))
    selectLoading.value = false
  }
}

// 选择vc时触发
const updateSelectVc = (value) => {
  vcStore.vcClustersList.some(async (v) => {
    if (v.name === value) {
      spinLoading.value = true
      const loginRes = await vcStore.loginVcCluster({data: {id: v.id}})
      if (loginRes) {
        const resData = await vcStore.getCloneVmMetadata()
        vmNames.value = resData.vms
        templates.value = resData.templates
        hosts.value = resData.hosts
        folders.value = resData.folders
        hostOptions.value = hosts.value.map((v) => ({
          label: v.name,
          value: v.name
        }))
        folderOptions.value = folders.value.map((v) => ({
          label: v,
          value: v
        }))
        templateOptions.value = templates.value.map((v) => ({
          label: v.name,
          value: v.name
        }))
        spinLoading.value = false
      }
      return true
    }
  })
}
// 选择主机时触发
const updateSelectHost = (value) => {
  hostNetworks.value = {}
  hosts.value.some((v) => {
    if (v.name === value) {
      datastoreOptions.value = v.storages.map((storage) => ({
        label: storage.name,
        value: storage.name
      }))
      networkOptions.value = [
        {
          type: 'group',
          label: 'VSS',
          key: 'vss',
          children: []
        },
        {
          type: 'group',
          label: 'VDS',
          key: 'vds',
          children: []
        },
      ]
      v.networks.forEach((network) => {
        hostNetworks.value[network.network_name] = network.network_type
        if (network.network_type === "vss"){
          networkOptions.value[0].children.push({
            label: network.network_name,
            value: network.network_name
          })
        }
        if (network.network_type === "vds"){
          networkOptions.value[1].children.push({
            label: network.network_name,
            value: network.network_name
          })
        }
      })
      console.log("networkOptions.value", networkOptions.value)
      return true
    }
  })
}
// 选择模版时时触发
const updateSelectTemplate = (value) => {
  templates.value.some((v) => {
    if (v.name === value) {
      console.log("v: ", v)
      cloneVmForm.value.cpu_num = v.cpu_num
      cloneVmForm.value.memory = v.memory
      cloneVmForm.value.current_disk = v.storage
      cloneVmForm.value.networks = v.network
      return true
    }
  })
}


// 表单校验
const cloneVmRules = {
  vm_name: {
    required: true,
    trigger: ['blur'],
    validator: (rule, value) => {
      console.log(value, vmNames)
      if (vmNames.value.length > 0) {
        const valRes = vmNames.value.some((vm_name) => vm_name === value)
        if (valRes) {
          return Error("虚拟机名字重复.")
        } else {
          return true
        }
      } else {
        return true
      }
    }
  },
  template: {
    required: true,
    trigger: ['blur'],
  },
  host: {
    required: true,
    trigger: ['blur'],
  },
  datastore: {
    required: true,
    trigger: ['blur'],
  },
  cpu_num: {
    required: true,
    trigger: ['blur'],
  },
  memory: {
    required: true,
    trigger: ['blur'],
  }
}

// 新增磁盘时添加默认值
const diskOnCreate = () => ({
  disk_capacity_gb: 0,
  disk_thin: true
})

// 新增网络时添加默认值
const networkOnCreate = () => ({
  network_name: "",
  network_type: ""
})

// 更改网络时自动选择网卡类型
const networkOnUpdate = (index) => {
  const name = cloneVmForm.value.networks[index].network_name
  cloneVmForm.value.networks[index].network_type = hostNetworks.value[name]
}

const handleOkClick = async (e) => {
  e.preventDefault()
  cloneVmForm.value["memory_mb"] = cloneVmForm.value.memory * 1024
  console.log("cloneVmForm.value", cloneVmForm.value)
  cloneVmStore.cloneVmJobList = await postCloneVmJobs(cloneVmForm.value)
  handleCloneVmClick()
  cloneVmForm.value = initCloneVmForm()
}

// 初始化表单
const cloneVmForm = ref({
  vc_name: "",
  template: "",
  vm_name: "",
  host: "",
  datastore: "",
  folder: "",
  cpu_num: 0,
  memory: 0,
  current_disk: [],
  disks: [],
  networks: [],
  ip: {
    ip_addr: "",
    net_mask: "",
    gateway: "",
    dns: []
  }
})

// 点击克隆虚拟机按钮时
const handleCloneVmClick = () => {
  showCloneModal.value = !showCloneModal.value
}

// 模态框关闭调用
const handleModelCloseClick = () => {
  cloneVmForm.value = initCloneVmForm()
}
const dnsOptions = ref([
  {
    label: '172.28.0.46',
    key: '172.28.0.46',
  },
  {
    label: '172.28.25.34',
    key: '172.28.25.34',
  },
  {
    label: '172.28.25.35',
    key: '172.28.25.35',
  },
  {
    label: '172.28.25.36',
    key: '172.28.25.36',
  }
])
</script>

