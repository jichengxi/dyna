<template>
  <div p15>
<!--    按钮-->
    <div>
      <n-button ghost color="#8a2be2" mr-20px @click="handleAddClusterClick">
        录入vcenter集群
      </n-button>
      <n-modal v-model:show="showModal" :mask-closable="false" preset="dialog" title="录入vCenter集群">
        <template #header>
          <div>录入vCenter集群</div>
        </template>
        <div>
          <n-form ref="formRef" :model="vcClusterModel" :rules="rules" label-placement="left" label-width="auto" require-mark-placement="right">
            <n-form-item path="name" label="name">
              <n-input v-model:value="vcClusterModel.name" @keydown.enter.prevent placeholder="vCenter集群名" clearable/>
            </n-form-item>
            <n-form-item path="host" label="host">
              <n-input v-model:value="vcClusterModel.host" @keydown.enter.prevent placeholder="vCenter地址" clearable/>
            </n-form-item>
            <n-form-item path="user" label="user">
              <n-input v-model:value="vcClusterModel.user" @keydown.enter.prevent placeholder="vCenter用户" clearable/>
            </n-form-item>
            <n-form-item path="password" label="password">
              <n-input v-model:value="vcClusterModel.password" type="password" @keydown.enter.prevent placeholder="vCenter密码" clearable/>
            </n-form-item>
            <n-form-item path="tags" label="tags">
              <n-dynamic-tags v-model:value="vcClusterModel.tags"/>
            </n-form-item>
          </n-form>
          <pre>{{ JSON.stringify(vcClusterModel, null, 2) }}</pre>
        </div>
        <template #action>
          <n-button
              :disabled="vcClusterModel.name === null"
              round
              type="primary"
              @click="handleValidateClick"
          >
            验证
          </n-button>

          <n-button
              :disabled="vcClusterModel.name === null"
              round
              type="primary"
              @click="handleOkClick"
          >
            确定
          </n-button>
        </template>
      </n-modal>
    </div>

<!--    展示-->
    <div>
      <n-space vertical :size="12">
        <n-data-table
            :bordered="false"
            :columns="vcColumns"
            :data="vcStore.vcClustersList"
            :pagination="pagination"
            :loading="loading"
        />
      </n-space>
    </div>
  </div>
</template>

<script setup>
import {h, onMounted, ref} from 'vue'
import { useVcStore } from "@/store/modules/vc/vc-manage"
import {NTag, useMessage} from 'naive-ui'
import VcDeletePopConfirm from './vc-delete-popconfirm.vue'

const message = useMessage()
const vcStore = useVcStore()
const components = [VcDeletePopConfirm]

// 表格数据
const vcColumns = createColumns()
const pagination = {pageSize: 10}
const loading = ref(false)
function createColumns() {
  return [
    {
      title: 'name',
      key: 'name',
    },
    {
      title: 'host',
      key: 'host',
    },
    {
      title: 'user',
      key: 'user',
    },
    {
      title: 'password',
      key: 'password',
    },
    {
      title: 'tags',
      key: 'tags',
      render(row) {
        return row.tags.map((tagKey) => {
          return h(
              NTag,
              {
                style: {
                  marginRight: '6px'
                },
                type: 'info'
              },
              {
                default: () => tagKey
              }
          )
        })
      }
    },
    {
      title: 'Action',
      key: 'actions',
      render(row) {
        return h(
            VcDeletePopConfirm,
            {
                "id": row.id,
                "name": row.name,
                onDelCluster: delCluster
            }
        )
      }
    }
  ]
}

// 弹窗数据
const formRef = ref(null)
const showModal = ref(false)
const vcClusterModel = ref({
  name: "",
  host: "",
  user: "",
  password: "",
  tags: []
})

const rules = {
  tags: {
    trigger: ['change'],
    validator(rule, value) {
      console.log("rule,", rule)
      console.log("value,", value)
      if (value.length >= 5) return new Error('不得超过四个标签')
      return true
    }
  }
}

const delCluster = async (id) => {
  loading.value = true
  const res = await vcStore.delVcCluster({data: {id}})
  if (res) {
    vcStore.vcClustersList.map((value, index ) => {
      if (value.id === id) {
        vcStore.vcClustersList.splice(index, 1)
      }
    })
    message.success("删除成功")
  } else {
    message.error("删除失败")
  }
  loading.value = false
}

// 点击添加集群按钮触发事件
const handleAddClusterClick = () => {
  showModal.value = true
}

// 点击验证触发事件
const handleValidateClick = (e) => {
  e.preventDefault()
  formRef.value?.validate((errors) => {
    if (!errors) {
      message.success('验证成功')
    } else {
      console.log(errors)
      message.error('验证失败')
    }
  })
}

// 点击确定触发事件
const handleOkClick = (e) => {
  e.preventDefault()
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      let {name, host, user, password, tags} = vcClusterModel.value
      tags.push(name)
      loading.value = true
      await vcStore.addVcCluster({data: {name, host, user, password: password.toString(), tags}})
      vcClusterModel.value = {
        name: "",
        host: "",
        user: "",
        password: "",
        tags: []
      }
      loading.value = false
      showModal.value = false
    } else {
      console.log(errors)
      message.error('添加失败')
    }
  })
}

onMounted( async () => {
  loading.value = true
  await vcStore.getVcClusters()
  loading.value = false
})

</script>
