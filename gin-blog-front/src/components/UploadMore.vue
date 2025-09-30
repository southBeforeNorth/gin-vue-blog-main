<script setup>
import { computed, ref, watch } from 'vue'
import { useUserStore } from '@/store'
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

const { token } = useUserStore()
const previewImgs = ref([...props.preview])
const uploading = ref(false)

watch(() => props.preview, val => previewImgs.value = [...val])

// 上传图片
async function handleFileSelect(event) {
  const file = event.target.files[0]
  if (!file) return

  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    alert('请选择图片文件')
    return
  }

  // 检查文件大小 (5MB)
  if (file.size > 50 * 1024 * 1024) {
    alert('图片大小不能超过50MB')
    return
  }

  uploading.value = true

  try {
    const formData = new FormData()
    formData.append('file', file)

    const response = await fetch('/api/upload', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
      body: formData,
    })

    const res = await response.json()
    
    if (res.code !== 0) {
      alert(res.message || '上传失败')
      return
    }

    previewImgs.value.push(res.data)
    emit('update:preview', previewImgs.value)
  } catch (error) {
    console.error('上传失败:', error)
    alert('上传失败，请重试')
  } finally {
    uploading.value = false
    // 清空文件输入
    event.target.value = ''
  }
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
    <div class="upload-container">
      <!-- 已上传的图片 -->
      <div
        v-for="(img, index) in previewImgs"
        :key="index"
        class="image-item"
      >
        <img
          class="preview-image"
          :src="convertImgUrl(img)"
          alt="上传图片"
        >
        <button
          class="remove-btn"
          @click="removeImg(index)"
        >
          ×
        </button>
      </div>

      <!-- 上传按钮 -->
      <div
        v-if="canUpload"
        class="upload-area"
        @click="$refs.fileInput.click()"
      >
        <div class="upload-content">
          <div class="upload-icon">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
              <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z" />
            </svg>
          </div>
          <div class="upload-text">
            {{ uploading ? '上传中...' : `${previewImgs.length}/${props.max}` }}
          </div>
        </div>
      </div>

      <!-- 隐藏的文件输入 -->
      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        style="display: none"
        @change="handleFileSelect"
      >
    </div>
  </div>
</template>

<style scoped>
.upload-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
  max-width: 100%;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .upload-container {
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
  }
}

.image-item {
  position: relative;
  aspect-ratio: 1;
}

.preview-image {
  width: 100%;
  height: 100%;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  object-fit: cover;
  display: block;
}

.remove-btn {
  position: absolute;
  top: -6px;
  right: -6px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #ff4d4f;
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  line-height: 1;
  box-shadow: 0 2px 4px rgba(0,0,0,0.2);
}

.remove-btn:hover {
  background: #ff7875;
}

.upload-area {
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: border-color 0.3s;
  background: #fafafa;
  aspect-ratio: 1;
}

.upload-area:hover {
  border-color: #40a9ff;
  background: #f0f8ff;
}

.upload-content {
  text-align: center;
  color: #999;
  padding: 8px;
}

.upload-icon {
  margin-bottom: 4px;
}

.upload-text {
  font-size: 10px;
}

/* 移动端进一步优化 */
@media (max-width: 480px) {
  .upload-container {
    gap: 6px;
  }
  
  .remove-btn {
    width: 16px;
    height: 16px;
    font-size: 10px;
    top: -4px;
    right: -4px;
  }
  
  .upload-icon svg {
    width: 20px;
    height: 20px;
  }
  
  .upload-text {
    font-size: 9px;
  }
  
  .upload-content {
    padding: 4px;
  }
}
</style>
