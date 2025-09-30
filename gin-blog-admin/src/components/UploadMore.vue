<script setup>
import { computed, ref, watch } from 'vue'
import { NIcon, NText, NUpload, NUploadDragger, NSpace, NButton } from 'naive-ui'
import { useAuthStore } from '@/store'
import { convertImgUrl } from '@/utils'

const props = defineProps({
  preview: {
    type: Array,
    default: () => [],
  },
  width: {
    type: Number,
    default: 120,
  },
  max: {
    type: Number,
    default: 9,
  },
})

const emit = defineEmits(['update:preview'])

const { token } = useAuthStore()
const previewImgs = ref([...props.preview])

watch(() => props.preview, val => previewImgs.value = [...val])

// 上传图片
function handleImgUpload({ event }) {
  const respStr = (event?.target).response
  const res = JSON.parse(respStr)
  if (res.code !== 0) {
    $message?.error(res.message)
    return
  }
  previewImgs.value.push(res.data)
  emit('update:preview', previewImgs.value)
}

// 删除图片
function removeImg(index) {
  previewImgs.value.splice(index, 1)
  emit('update:preview', previewImgs.value)
}

// 判断是否可以继续上传
const canUpload = computed(() => previewImgs.value.length < props.max)

defineExpose({ previewImgs })
</script>

<template>
  <div>
    <NSpace>
      <!-- 已上传的图片 -->
      <div
        v-for="(img, index) in previewImgs"
        :key="index"
        class="relative"
      >
        <img
          border-color="#d9d9d9"
          class="border-2 rounded-lg border-dashed"
          :style="{ width: `${props.width}px`, height: `${props.width}px` }"
          :src="convertImgUrl(img)"
          alt="上传图片"
          object-fit="cover"
        >
        <NButton
          size="tiny"
          type="error"
          circle
          class="absolute -top-2 -right-2"
          @click="removeImg(index)"
        >
          ×
        </NButton>
      </div>

      <!-- 上传按钮 -->
      <NUpload
        v-if="canUpload"
        action="/api/upload"
        :headers="{ Authorization: `Bearer ${token}` }"
        :show-file-list="false"
        @finish="handleImgUpload"
      >
        <NUploadDragger
          :style="{ width: `${props.width}px`, height: `${props.width}px` }"
        >
          <div class="mb-3">
            <NIcon size="30" :depth="3">
              <span class="i-mdi:upload" />
            </NIcon>
          </div>
          <NText depth="3" style="font-size: 12px">
            {{ previewImgs.length }}/{{ props.max }}
          </NText>
        </NUploadDragger>
      </NUpload>
    </NSpace>
  </div>
</template>

<style scoped>
.relative {
  position: relative;
}
.absolute {
  position: absolute;
}
.-top-2 {
  top: -8px;
}
.-right-2 {
  right: -8px;
}
</style>
