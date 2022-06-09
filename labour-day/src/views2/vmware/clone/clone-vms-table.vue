<template>
    <n-button ghost color="#8a2be2" mr-20px @click="handleCloneVmsClick">
      批量克隆虚拟机
    </n-button>
    <n-modal v-model:show="showClonesModal" :mask-closable="false" preset="card" title="克隆虚拟机" :style="{width: '600px'}" :bordered="false" :closable="true">
      <n-button type="primary" ghost @click="handleDownloadTemplate">下载克隆虚拟机模版</n-button>
      <template #footer>
        <n-upload ref="uploadRef" :action=uploadAction @change="handleChangeUploadFile" :default-upload="false" accept=".xlsx" :directory-dnd="true" :max="1" method="POST">
          <n-button>选择模版</n-button>
        </n-upload>
      </template>
      <template #action>
        <n-button :disabled="!uploadFileList" @click="handleUploadClick">上传模版</n-button>
      </template>
    </n-modal>
</template>

<script setup>
import {ref} from "vue"
import {getCloneVmTemplate} from "./clone-methods"
const showClonesModal = ref(false)

const uploadAction = import.meta.env.VITE_APP_GLOB_BASE_API + "/vmware/clone-vm/template"
const uploadRef = ref(null)
const uploadFileList = ref(0)
const handleChangeUploadFile = (options) => {
  uploadFileList.value = options.fileList.length;
}
const handleUploadClick = () => {
  uploadRef.value?.submit();
  handleCloneVmsClick()
}

const handleDownloadTemplate = async () => {
  await getCloneVmTemplate()
}

const handleCloneVmsClick = () => {
  showClonesModal.value = !showClonesModal.value
}
</script>


