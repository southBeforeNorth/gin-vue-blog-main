<script setup>
import { onMounted, ref, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import dayjs from 'dayjs'

import BannerPage from '@/components/BannerPage.vue'
import api from '@/api'

const router = useRouter()

const loading = ref(true)
const loadingMore = ref(false) // 加载更多的loading状态
const total = ref(0)
const archiveList = ref([])
const pageSize = 10
const current = ref(1)

// 计算总页数
const totalPages = computed(() => Math.ceil(total.value / pageSize))

// 是否还有更多数据
const hasMore = computed(() => current.value < totalPages.value)

async function getArchives(isLoadMore = false) {
  if (isLoadMore) {
    loadingMore.value = true
  } else {
    loading.value = true
  }
  
  try {
    const resp = await api.getArchives({
      page_num: current.value,
      page_size: pageSize,
    })
    
    if (isLoadMore) {
      // 追加数据
      archiveList.value = [...archiveList.value, ...resp.data.page_data]
    } else {
      // 替换数据
      archiveList.value = resp.data.page_data
    }
    
    total.value = resp.data.total
  } catch (error) {
    console.error('获取归档数据失败:', error)
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

// 加载下一页
async function loadMore() {
  if (loadingMore.value || !hasMore.value) return
  
  current.value++
  await getArchives(true)
}

// 滚动事件处理
function handleScroll() {
  // 防抖处理
  if (loadingMore.value || !hasMore.value) return
  
  const scrollTop = window.pageYOffset || document.documentElement.scrollTop
  const windowHeight = window.innerHeight
  const documentHeight = document.documentElement.scrollHeight
  
  // 距离底部还有100px时开始加载
  if (scrollTop + windowHeight >= documentHeight - 100) {
    loadMore()
  }
}

// 节流函数
function throttle(func, wait) {
  let timeout
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout)
      func(...args)
    }
    clearTimeout(timeout)
    timeout = setTimeout(later, wait)
  }
}

const throttledScroll = throttle(handleScroll, 200)

onMounted(() => {
  getArchives()
  
  // 添加滚动监听
  window.addEventListener('scroll', throttledScroll)
  
  // 组件卸载时移除监听
  onUnmounted(() => {
    window.removeEventListener('scroll', throttledScroll)
  })
})
</script>

<template>
  <BannerPage title="时光轴线" label="archive" :loading="loading" card>
    <p class="pb-5 text-lg lg:text-2xl">
      目前共计 {{ total }} 篇文章
    </p>
    
    <template v-for="(item, idx) of archiveList" :key="item.id">
      <div class="group flex items-center gap-3 py-2 hover:bg-gray-50 rounded-lg px-2 transition-all duration-200">
        <div class="i-mdi:circle bg-blue text-sm group-hover:bg-orange transition-colors" />
        <span class="text-sm color-#666 lg:text-base font-mono min-w-24">
          {{ dayjs(item.created_at).format('YYYY-MM-DD') }}
        </span>
        <a class="color-#666 lg:text-lg hover:text-orange cursor-pointer flex-1 group-hover:translate-x-1 transform transition-transform" 
           @click="router.push(`/article/${item.id}`)">
          {{ item.title }}
        </a>
        <div class="i-mdi:arrow-right text-gray-400 opacity-0 group-hover:opacity-100 group-hover:text-orange transition-all" />
      </div>
      <hr v-if="idx !== archiveList.length - 1" class="my-3 border-1 border-color-#d2ebfd border-dashed">
    </template>

    <!-- 加载更多状态 -->
    <div v-if="loadingMore" class="flex justify-center items-center py-8">
      <div class="flex items-center gap-3 text-blue-500">
        <div class="i-mdi:loading animate-spin text-xl" />
        <span class="text-lg">正在加载更多...</span>
      </div>
    </div>

    <!-- 没有更多数据提示 -->
    <div v-else-if="!hasMore && archiveList.length > 0" class="flex justify-center items-center py-8">
      <div class="flex items-center gap-3 text-gray-500">
        <div class="i-mdi:check-circle text-xl text-green-500" />
        <span class="text-lg">已经到底了，共 {{ total }} 篇文章</span>
      </div>
    </div>

    <!-- 手动加载按钮（备用方案，可选） -->
    <div v-if="hasMore && !loadingMore" class="flex justify-center items-center py-8">
      <button 
        @click="loadMore"
        class="group flex items-center gap-3 px-8 py-4 bg-gradient-to-r from-blue-500 to-blue-600 
               text-white rounded-full shadow-lg hover:shadow-xl hover:from-blue-600 hover:to-blue-700 
               transform hover:scale-105 transition-all duration-300"
      >
        <div class="i-mdi:arrow-down text-xl group-hover:animate-bounce" />
        <span class="font-medium text-lg">加载更多 ({{ current }}/{{ totalPages }})</span>
      </button>
    </div>
  </BannerPage>
</template>
