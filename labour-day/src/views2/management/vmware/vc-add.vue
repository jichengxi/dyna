<template>
  <div p24>
    <n-form ref="formRef" :model="model" :rules="rules" label-placement="left" label-width="auto" require-mark-placement="right">
      <n-form-item path="name" label="name" >
        <n-input v-model:value="model.name" @keydown.enter.prevent autosize style="min-width: 10%" placeholder="vCenter集群名"/>
      </n-form-item>
      <n-form-item path="ip" label="ip">
        <n-input v-model:value="model.ip" @keydown.enter.prevent autosize style="min-width: 20%" placeholder="vCenter地址"/>
      </n-form-item>
      <n-form-item path="user" label="user">
        <n-input v-model:value="model.user" @keydown.enter.prevent default-value="administrator@vsphere.local" autosize style="min-width: 20%" placeholder="vCenter用户"/>
      </n-form-item>
      <n-form-item path="password" label="password">
        <n-input v-model:value="model.password" type="password" @keydown.enter.prevent autosize style="min-width: 10%" placeholder="vCenter密码"/>
      </n-form-item>
      <n-form-item path="tags" label="tags">
        <n-dynamic-tags v-model:value="model.tags" />
      </n-form-item>

      <n-row :gutter="[0, 24]">
        <n-col :span="24">
          <div style="display: flex; justify-content: flex-end">
            <n-button
                :disabled="model.age === null"
                round
                type="primary"
                @click="handleValidateButtonClick"
            >
              验证
            </n-button>
          </div>
        </n-col>
      </n-row>
    </n-form>

    <pre>{{ JSON.stringify(model, null, 2) }}</pre>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useMessage } from 'naive-ui'

const formRef = ref(null)
const message = useMessage()
const model = ref({
  name: null,
  ip: null,
  user: null,
  password: null,
  tags: []
})

const rules =  {
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

function handleValidateButtonClick (e) {
  e.preventDefault()
  console.log(formRef.value)
  formRef.value?.validate((errors) => {
    if (!errors) {
      message.success('验证成功')
    } else {
      console.log(errors)
      message.error('验证失败')
    }
  })
}

</script>

