<template>
  <div class="mobile-publish-page">
    <!-- 顶部导航栏 -->
    <div class="header">
      <button class="back-btn" @click="handleBack">
        <span>返回</span>
      </button>
      <h1>写日记</h1>
      <button 
        class="save-btn" 
        @click="handlePublish"
        :disabled="btnLoading || !form.content.trim()"
      >
        {{ btnLoading ? '发布中...' : '发布' }}
      </button>
    </div>

    <!-- 发布表单 -->
    <div class="publish-form">
      <!-- 内容输入 -->
      <div class="form-group">
        <textarea 
          v-model="form.content"
          placeholder="记录此刻的心情..."
          class="content-textarea"
          maxlength="200"
          @input="adjustTextareaHeight"
          ref="contentTextarea"
        ></textarea>
        <div class="char-count">{{ form.content.length }}/200</div>
      </div>

      <!-- 图片上传 -->
      <div class="form-group">
        <div class="upload-section">
          <UploadMore
            v-model:preview="form.imgs"
          />
        </div>
      </div>

      <!-- 发布设置 -->
      <div class="form-group">
        <div class="publish-settings">
          <div class="setting-item">
            <span class="setting-label">记录时间</span>
            <button class="time-selector" @click="showTimePicker = true">
              {{ formatDateTime(selectedDateTime) }}
              <i class="fas fa-chevron-right"></i>
            </button>
          </div>

          <!-- 发布状态 -->
          <div class="setting-item">
            <span class="setting-label">发布状态</span>
            <div class="status-options">
              <label class="status-option">
                <input 
                  v-model="form.status" 
                  type="radio" 
                  :value="1"
                />
                <span>公开</span>
              </label>
              <label class="status-option">
                <input 
                  v-model="form.status" 
                  type="radio" 
                  :value="2"
                />
                <span>私密</span>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 加载提示 -->
    <div v-if="btnLoading" class="loading-overlay">
      <div class="loading-content">
        <div class="loading-spinner"></div>
        <p>发布中...</p>
      </div>
    </div>

    <!-- 时间选择器弹窗 -->
    <div v-if="showTimePicker" class="time-picker-overlay" @click="showTimePicker = false">
      <div class="time-picker-modal" @click.stop>
        <div class="time-picker-header">
          <button class="cancel-btn" @click="showTimePicker = false">取消</button>
          <h3>选择时间</h3>
          <button class="confirm-btn" @click="confirmTime">确定</button>
        </div>
        
        <div class="time-picker-content">
          <!-- 快捷选择 -->
          <div class="quick-select">
            <button 
              v-for="quick in quickTimeOptions" 
              :key="quick.label"
              class="quick-btn"
              @click="selectQuickTime(quick)"
            >
              {{ quick.label }}
            </button>
          </div>
          
          <!-- 日期时间选择 -->
          <div class="datetime-picker">
            <input 
              type="datetime-local" 
              v-model="tempDateTime"
              class="datetime-input"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Toast 提示 -->
    <div v-if="toastMessage" class="toast">
      {{ toastMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import UploadMore from '@/components/UploadMore.vue'

defineOptions({ name: '移动端日记发布' })

const router = useRouter()

// 响应式数据
const btnLoading = ref(false)
const toastMessage = ref('')

// 时间相关数据
const showTimePicker = ref(false)
const selectedDateTime = ref(new Date())
const tempDateTime = ref('')

// 快捷时间选项
const quickTimeOptions = [
  { label: '现在', offset: 0 },
  { label: '1小时前', offset: -1 * 60 * 60 * 1000 },
  { label: '今天早上', hour: 8 },
  { label: '昨天', offset: -24 * 60 * 60 * 1000 }
]

// 时间相关方法
const formatDateTime = (date) => {
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const targetDate = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  
  const diffDays = Math.floor((targetDate - today) / (24 * 60 * 60 * 1000))
  
  let dateStr = ''
  if (diffDays === 0) {
    dateStr = '今天'
  } else if (diffDays === -1) {
    dateStr = '昨天'
  } else {
    dateStr = date.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit' })
  }
  
  const timeStr = date.toLocaleTimeString('zh-CN', { 
    hour: '2-digit', 
    minute: '2-digit',
    hour12: false 
  })
  
  return `${dateStr} ${timeStr}`
}

const selectQuickTime = (quick) => {
  const now = new Date()
  let targetTime
  
  if (quick.offset !== undefined) {
    targetTime = new Date(now.getTime() + quick.offset)
  } else if (quick.hour !== undefined) {
    targetTime = new Date()
    targetTime.setHours(quick.hour, 0, 0, 0)
  }
  
  selectedDateTime.value = targetTime
  updateTempDateTime()
}

const updateTempDateTime = () => {
  const date = selectedDateTime.value
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  
  tempDateTime.value = `${year}-${month}-${day}T${hours}:${minutes}`
}

const confirmTime = () => {
  if (tempDateTime.value) {
    selectedDateTime.value = new Date(tempDateTime.value)
    form.value.add_time = selectedDateTime.value.getTime()
  }
  showTimePicker.value = false
}

const form = ref({
  content: '',
  imgs: [],
  add_time: Date.now(),
  status: 1, // 1:公开 2:私密
})

// 方法
const handleBack = () => {
  if (form.value.content.trim() || form.value.imgs.length > 0) {
    if (confirm('确定要离开吗？未保存的内容将丢失')) {
      router.back()
    }
  } else {
    router.back()
  }
}

const handlePublish = async () => {
  if (!form.value.content.trim()) {
    showToast('请输入日记内容')
    return
  }
  
  btnLoading.value = true
  
  try {
    // 使用选择的时间
    form.value.add_time = selectedDateTime.value.getTime()
    await api.saveOrUpdateDiary(form.value)
    showToast('发布成功')
    
    // 延迟跳转，让用户看到成功提示
    setTimeout(() => {
      router.push('/diary')
    }, 1500)
    
  } catch (error) {
    console.error('发布失败:', error)
    showToast('发布失败，请重试')
  } finally {
    btnLoading.value = false
  }
}

const adjustTextareaHeight = () => {
  nextTick(() => {
    const textarea = document.querySelector('.content-textarea')
    if (textarea) {
      textarea.style.height = 'auto'
      textarea.style.height = Math.max(120, textarea.scrollHeight) + 'px'
    }
  })
}

const showToast = (message) => {
  toastMessage.value = message
  setTimeout(() => {
    toastMessage.value = ''
  }, 2000)
}

// 生命周期
onMounted(() => {
  // 初始化时间
  const now = new Date()
  selectedDateTime.value = now
  form.value.add_time = now.getTime()
  updateTempDateTime()
})
</script>

<style scoped>
.mobile-publish-page {
  min-height: 100vh;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
  position: relative;
}

/* 顶部导航栏 */
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,0,0,0.1);
  position: sticky;
  top: 0;
  z-index: 100;
  min-height: 56px;
}

.back-btn {
  padding: 10px 20px;
  border-radius: 22px;
  border: none;
  background: #d9d9d9;
  color: white;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 70px;
  height: 44px;
  -webkit-tap-highlight-color: transparent;
}

.back-btn:active {
  background: #e8e8e8;
  transform: scale(0.95);
}

.header h1 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
  flex: 1;
  text-align: center;
}

.save-btn {
  padding: 10px 20px;
  border-radius: 22px;
  border: none;
  background: #1890ff;
  color: white;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  min-width: 70px;
  height: 44px;
  -webkit-tap-highlight-color: transparent;
}

.save-btn:disabled {
  background: #d9d9d9;
  color: #bbb;
  cursor: not-allowed;
}

.save-btn:not(:disabled):active {
  background: #1677cc;
  transform: scale(0.98);
}

/* 发布表单 */
.publish-form {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
}

.form-group {
  margin-bottom: 16px;
}

.form-group:last-child {
  margin-bottom: 0;
}

/* 内容输入框 */
.content-textarea {
  width: 100%;
  min-height: 140px;
  padding: 16px;
  border: none;
  border-radius: 12px;
  background: #fff;
  font-size: 16px;
  line-height: 1.6;
  color: #333;
  resize: none;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  box-sizing: border-box;
  transition: box-shadow 0.3s;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.content-textarea:focus {
  outline: none;
  box-shadow: 0 2px 12px rgba(24,144,255,0.2);
}

.content-textarea::placeholder {
  color: #bbb;
  font-size: 16px;
}

.char-count {
  text-align: right;
  font-size: 12px;
  color: #999;
  margin-top: 8px;
  padding-right: 4px;
}

.char-count.warning {
  color: #fa8c16;
}

.char-count.danger {
  color: #ff4d4f;
  font-weight: 500;
}

/* 图片上传区域 */
.upload-section {
  background: #fff;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}

.uploaded-images {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  margin-bottom: 12px;
}

.uploaded-img-item {
  position: relative;
  aspect-ratio: 1;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f5f5;
}

.uploaded-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.remove-img-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  border: none;
  background: rgba(0,0,0,0.6);
  color: white;
  font-size: 16px;
  line-height: 1;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
  -webkit-tap-highlight-color: transparent;
}

.remove-img-btn:active {
  background: rgba(0,0,0,0.8);
}

.add-image-btn {
  width: 100%;
  min-height: 80px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  background: #fafafa;
  color: #999;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  transition: all 0.2s;
  -webkit-tap-highlight-color: transparent;
}

.add-image-btn:active {
  border-color: #1890ff;
  color: #1890ff;
  background: #f0f8ff;
}

.add-image-btn i {
  font-size: 20px;
}

/* 发布设置 */
.publish-settings {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
  overflow: hidden;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-label {
  font-size: 15px;
  color: #333;
  font-weight: 500;
}

.time-selector {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border: none;
  border-radius: 6px;
  background: #f5f5f5;
  color: #666;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;
  -webkit-tap-highlight-color: transparent;
}

.time-selector:active {
  background: #e8e8e8;
}

.time-selector i {
  font-size: 12px;
  opacity: 0.6;
}

.status-options {
  display: flex;
  gap: 16px;
}

.status-option {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  -webkit-tap-highlight-color: transparent;
}

.status-option input[type="radio"] {
  width: 16px;
  height: 16px;
  accent-color: #1890ff;
}

/* 时间选择器弹窗 */
.time-picker-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  z-index: 1000;
  display: flex;
  align-items: flex-end;
  animation: fadeIn 0.3s ease;
}

.time-picker-modal {
  width: 100%;
  background: #fff;
  border-radius: 16px 16px 0 0;
  animation: slideUp 0.3s ease;
  max-height: 60vh;
  overflow: hidden;
}

.time-picker-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #f0f0f0;
}

.time-picker-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #333;
}

.cancel-btn, .confirm-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: background 0.2s;
  -webkit-tap-highlight-color: transparent;
}

.cancel-btn {
  background: #f5f5f5;
  color: #666;
}

.cancel-btn:active {
  background: #e8e8e8;
}

.confirm-btn {
  background: #1890ff;
  color: white;
}

.confirm-btn:active {
  background: #1677cc;
}

.time-picker-content {
  padding: 20px;
}

.quick-select {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 24px;
}

.quick-btn {
  padding: 12px;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  background: #fff;
  color: #666;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
  -webkit-tap-highlight-color: transparent;
}

.quick-btn:active {
  border-color: #1890ff;
  color: #1890ff;
  background: #f0f8ff;
}

.datetime-picker {
  border-top: 1px solid #f0f0f0;
  padding-top: 20px;
}

.datetime-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #d9d9d9;
  border-radius: 8px;
  font-size: 16px;
  color: #333;
  background: #fff;
  box-sizing: border-box;
}

.datetime-input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 2px rgba(24,144,255,0.2);
}

/* 加载遮罩 */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-content {
  background: #fff;
  padding: 32px;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 8px 32px rgba(0,0,0,0.2);
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #1890ff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

.loading-content p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

/* Toast 提示 */
.toast {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0,0,0,0.8);
  color: white;
  padding: 12px 20px;
  border-radius: 8px;
  font-size: 14px;
  z-index: 3000;
  animation: fadeInOut 2s ease;
}

/* 动画 */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { transform: translateY(100%); }
  to { transform: translateY(0); }
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@keyframes fadeInOut {
  0%, 100% { opacity: 0; }
  20%, 80% { opacity: 1; }
}

/* 移动端适配 */
@media (max-width: 375px) {
  .header {
    padding: 10px 12px;
  }
  
  .publish-form {
    padding: 12px;
  }
  
  .content-textarea {
    padding: 12px;
    font-size: 15px;
  }
  
  .upload-section {
    padding: 12px;
  }
  
  .uploaded-images {
    gap: 6px;
  }
  
  .quick-select {
    grid-template-columns: 1fr;
    gap: 8px;
  }
}

/* 安全区域适配 */
@supports (padding: max(0px)) {
  .mobile-publish-page {
    padding-top: env(safe-area-inset-top);
    padding-bottom: env(safe-area-inset-bottom);
  }
  
  .header {
    padding-top: max(12px, env(safe-area-inset-top) + 12px);
  }
}

/* 深色模式支持 */
@media (prefers-color-scheme: dark) {
  .mobile-publish-page {
    background: #1a1a1a;
  }
  
  .header {
    background: #2d2d2d;
    color: #fff;
  }
  
  .header h1 {
    color: #fff;
  }
  
  .back-btn {
    background: #3d3d3d;
    color: #ccc;
  }
  
  .content-textarea,
  .upload-section,
  .publish-settings {
    background: #2d2d2d;
    color: #fff;
  }
  
  .content-textarea::placeholder {
    color: #666;
  }
  
  .setting-item {
    border-bottom-color: #3d3d3d;
  }
  
  .time-selector {
    background: #3d3d3d;
    color: #ccc;
  }
  
  .add-image-btn {
    background: #2d2d2d;
    border-color: #3d3d3d;
    color: #666;
  }
}
</style>