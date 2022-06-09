<template>
  <div p10>
    <n-spin :show="show">
      <n-descriptions label-placement="left" title="概要" size="small">
        <n-descriptions-item label="Name">
          {{vmHostInfo.name}}
        </n-descriptions-item>
        <n-descriptions-item label="Cluster">
          {{vmHostInfo.cluster}}
        </n-descriptions-item>
        <n-descriptions-item label="DataCenter">
          {{vmHostInfo.datacenter}}
        </n-descriptions-item>
      </n-descriptions>

      <n-descriptions label-placement="left" title="CPU" size="small">
        <n-descriptions-item label="Total">
          {{vmHostInfo.cpu.total_cpu}}
        </n-descriptions-item>
        <n-descriptions-item label="Used">
          {{vmHostInfo.cpu.used_cpu}}
        </n-descriptions-item>
        <n-descriptions-item label="Free">
          {{vmHostInfo.cpu.free_cpu}}
        </n-descriptions-item>
        <n-descriptions-item label="Usage Percent">
          {{vmHostInfo.cpu.cpu_usage_percent}}
        </n-descriptions-item>
      </n-descriptions>

      <n-descriptions label-placement="left" title="Memory" size="small">
        <n-descriptions-item label="Total">
          {{vmHostInfo.memory.total_memory}}
        </n-descriptions-item>
        <n-descriptions-item label="Used">
          {{vmHostInfo.memory.used_memory}}
        </n-descriptions-item>
        <n-descriptions-item label="Free">
          {{vmHostInfo.memory.free_memory}}
        </n-descriptions-item>
        <n-descriptions-item label="Usage Percent">
          {{vmHostInfo.memory.memory_usage_percent}}
        </n-descriptions-item>
      </n-descriptions>

      <n-descriptions label-placement="left" title="DataStore" size="small" :column="4">
        <template v-for="item in vmHostInfo.datastore">
          <n-descriptions-item label="name">
            {{item.name}}
          </n-descriptions-item>
          <n-descriptions-item label="Total">
            {{item.total_space}}
          </n-descriptions-item>
          <n-descriptions-item label="Free">
            {{item.free_space}}
          </n-descriptions-item>
          <n-descriptions-item label="Free Percent">
            {{item.datastore_free_percent}}
          </n-descriptions-item>
        </template>
      </n-descriptions>

      <n-descriptions label-placement="left" title="Status" size="small">
        <n-descriptions-item label="Overall Status">
          {{vmHostInfo.status.overall_status}}
        </n-descriptions-item>
        <n-descriptions-item label="Connection Status">
          {{vmHostInfo.status.connection_state}}
        </n-descriptions-item>
        <n-descriptions-item label="Power State">
          {{vmHostInfo.status.power_state}}
        </n-descriptions-item>
      </n-descriptions>

      <n-descriptions label-placement="left" title="hardware" size="small">
        <n-descriptions-item label="品牌">
          {{vmHostInfo.hardware.vendor}}
        </n-descriptions-item>
        <n-descriptions-item label="服务器型号">
          {{vmHostInfo.hardware.model}}
        </n-descriptions-item>
        <n-descriptions-item label="cpu型号">
          {{vmHostInfo.hardware.cpu_model}}
        </n-descriptions-item>
        <n-descriptions-item label="物理cpu">
          {{vmHostInfo.hardware.num_cpu_pkgs}}
        </n-descriptions-item>
        <n-descriptions-item label="逻辑cpu">
          {{vmHostInfo.hardware.num_cpu_cores}}
        </n-descriptions-item>
        <n-descriptions-item label="cpu核数">
          {{vmHostInfo.hardware.num_cpu_threads}}
        </n-descriptions-item>
        <n-descriptions-item label="Esxi Full Name">
          {{vmHostInfo.hardware.esxi_full_name}}
        </n-descriptions-item>
        <n-descriptions-item label="Esxi版本">
          {{vmHostInfo.hardware.version}}
        </n-descriptions-item>
      </n-descriptions>

    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import {onBeforeMount, ref} from "vue"
import { getVmHost } from "@/views2/vmware/view/vc-methods"
import type { vmHost } from "@/store/modules/vc/vc-manage"

const route = useRoute()
const vmHostInfo = ref<vmHost>({
  name: "",
  cluster: "",
  vms: 0,
  cpu: {
    total_cpu: 0,
    used_cpu: 0,
    free_cpu: 0,
    cpu_usage_percent: 0
  },
  memory: {
    total_memory: 0,
    used_memory: 0,
    free_memory: 0,
    memory_usage_percent: 0
  },
  datastore: [],
  hardware: {
    vendor: "",
    model: "",
    cpu_model: "",
    num_cpu_pkgs: 0,
    num_cpu_cores: 0,
    num_cpu_threads: 0,
    esxi_full_name: "",
    version: "",
  },
  status: {
    overall_status: "",
    connection_state: "",
    power_state: "",
  },
  datacenter: "",
})
const show = ref<boolean>(false)

const host = route.query.host

onBeforeMount(async () => {
  show.value = true
  vmHostInfo.value = await getVmHost(host as string)
  show.value = false
})

</script>
