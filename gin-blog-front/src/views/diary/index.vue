<script setup>
import { onMounted, onUnmounted, ref, reactive, nextTick  } from 'vue'
import BannerPage from '@/components/BannerPage.vue'
import api from '@/api'
import { convertImgUrl } from '@/utils'

// 页面状态
const state = reactive({
  posts: [],
  loading: false,
  hasMore: true,
  page: 1,
  pageSize: 10
})


// 滑动手势状态
const swipeState = reactive({
  startY: 0,
  startX: 0,
  currentX: 0,
  currentY: 0,
  isDragging: false,
  isHorizontalSwipe: false,
  isVerticalSwipe: false,
  scale: 1,
  translateY: 0,
  translateX: 0,
  opacity: 1,
  startTime: 0
})

// 滑动提示状态
const swipeHint = reactive({
  visible: false,
  timer: null
})

// 检测是否为移动设备 - 修复版本
const isMobile = ref(false)

// 检测设备类型
const checkDeviceType = () => {
  const userAgent = navigator.userAgent.toLowerCase()
  const isMobileUA = /android|webos|iphone|ipad|ipod|blackberry|iemobile|opera mini/i.test(userAgent)
  const hasTouch = 'ontouchstart' in window || navigator.maxTouchPoints > 0
  const isSmallScreen = window.innerWidth <= 768
  
  // 只有同时满足触摸支持和小屏幕或移动端UA才认为是移动设备
  isMobile.value = hasTouch && (isSmallScreen || isMobileUA)
}

// 格式化时间
const formatTime = (timestamp) => {
  const now = new Date()
  const time = new Date(timestamp)
  const diff = now - time
  
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  return time.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}


// // 修改 loadPosts 函数
// const loadPosts = async (isRefresh = false) => {
//   if (state.loading) return
  
//   state.loading = true
  
//   try {
//     // 如果是刷新，重置页码
//     if (isRefresh) {
//       state.page = 1
//       state.hasMore = true
//     }
   
//     const resp = await api.getDiarys({
//       page: state.page,  // 使用当前页码
//       page_size: state.pageSize,
//       content:'',
//       add_start_time:0,
//       end_start_time:0
//     })
//     const newPosts = resp.data.page_data || []
//     const total = resp.data.total || 0 // 假设API返回总数

//     if (isRefresh) {
//       state.posts = newPosts
//     } else {
//       state.posts.push(...newPosts)
//     }

//     // 修复：无论是否刷新都要正确处理页码和hasMore
//     state.hasMore = state.posts.length < total
//     if (newPosts.length > 0) {
//       state.page++
//     }
//   } catch (error) {
//     console.error('加载失败:', error)
//   } finally {
//     state.loading = false
//   }
// }

const loadPosts = async (isRefresh = false) => {
  if (state.loading) return
  
  state.loading = true
  
  try {
    // 确定要请求的页码
    const requestPage = isRefresh ? 1 : state.page
    
    console.log(`请求第 ${requestPage} 页数据`)
    
    const resp = await api.getDiarys({
      page_num: requestPage,
      page_size: state.pageSize,
      content:'',
      add_start_time:0,
      end_start_time:0
    })
    
    const newPosts = resp.data.page_data || []
    const total = resp.data.total || 0

    // 更新数据
    if (isRefresh) {
      state.posts = newPosts
    } else {
      state.posts.push(...newPosts)
    }

    // 判断是否还有更多数据
    const hasMoreData = newPosts.length === state.pageSize && state.posts.length < total
    state.hasMore = hasMoreData
    
    // ✅ 正确的页码管理：成功获取数据后，下次请求的页码应该是当前页码+1
    if (newPosts.length > 0) {
      state.page = requestPage + 1  // 设置下次要请求的页码
    } else {
      state.hasMore = false  // 没有数据说明已经到底了
    }
    
    console.log(`成功加载 ${newPosts.length} 条数据，下次请求第 ${state.page} 页`)
    
  } catch (error) {
    console.error('加载失败:', error)
    // 请求失败时不修改页码，保持原状态以便重试
  } finally {
    state.loading = false
  }
}

// 节流滚动处理
// 节流滚动处理 - 优化版本
const handleScroll = () => {
  // 防止重复加载和无数据时的无效请求
  if (!state.hasMore || state.loading) return
  
  const { scrollTop, scrollHeight, clientHeight } = document.documentElement
  
  // 计算滚动进度
  const scrollProgress = (scrollTop + clientHeight) / scrollHeight
  const distanceFromBottom = scrollHeight - (scrollTop + clientHeight)
  
  // 多重判断条件：
  // 1. 距离底部小于100px
  // 2. 或者滚动进度超过95%
  // 3. 确保不是在页面顶部（避免初始化时误触发）
  if ((distanceFromBottom <= 400 || scrollProgress >= 0.70) && scrollTop > 0) {
    console.log('触底加载，距离底部:', distanceFromBottom, 'px')
    loadPosts()
  }
}


// 节流函数
const throttle = (func, wait) => {
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

const throttledHandleScroll = throttle(handleScroll, 200)

// 获取图片网格样式类
const getImageGridClass = (count) => {
  if (count === 1) return 'grid-cols-1'
  if (count === 2) return 'grid-cols-2'
  if (count === 3) return 'grid-cols-3'
  if (count === 4) return 'grid-cols-2'
  if (count >= 5) return 'grid-cols-3'
  return 'grid-cols-3'
}

// 显示滑动提示
const showSwipeHint = () => {
  if (!isMobile.value) return
  
  if (swipeHint.timer) {
    clearTimeout(swipeHint.timer)
  }
  
  swipeHint.visible = true
  
  swipeHint.timer = setTimeout(() => {
    swipeHint.visible = false
  }, 2500)
}

// 隐藏滑动提示
const hideSwipeHint = () => {
  if (swipeHint.timer) {
    clearTimeout(swipeHint.timer)
  }
  swipeHint.visible = false
}

// 图片预览
const previewImage = (url, images = [], index = 0) => {
  imagePreview.url = url
  imagePreview.images = images
  imagePreview.currentIndex = index
  imagePreview.visible = true
  imagePreview.loading = true
  document.body.style.overflow = 'hidden'
  
  // 重置滑动状态
  resetSwipeState()
  
  // 显示滑动提示（仅在移动端多图时）
  if (images.length > 1 && isMobile.value) {
    setTimeout(() => {
      showSwipeHint()
    }, 500)
  }
  
  // 预加载相邻图片
  preloadImages(images, index)
}

// 关闭预览
const closePreview = () => {
  imagePreview.visible = false
  imagePreview.loading = false
  document.body.style.overflow = ''
  resetSwipeState()
  hideSwipeHint()
}



// 重置滑动状态
const resetSwipeState = () => {
  Object.assign(swipeState, {
    startY: 0,
    startX: 0,
    currentX: 0,
    currentY: 0,
    isDragging: false,
    isHorizontalSwipe: false,
    isVerticalSwipe: false,
    translateX: 0,
    translateY: 0,
    scale: 1,
    opacity: 1,
    startTime: 0
  })

  setTimeout(() => {
    swipeState.translateX = 0
    swipeState.translateY = 0
    swipeState.opacity = 1
  }, 10)
  // 注意：不要重置 imagePreview 的缩放状态，除非是关闭预览
}


// 触摸开始处理 - 完全修复版本
const handleTouchStart = (e) => {
  if (!isMobile.value || !imagePreview.visible) return
  
  if (e.touches.length === 2) {
    imagePreview.initialDistance = getDistance(e.touches[0], e.touches[1])
    imagePreview.lastScale = imagePreview.scale
    return
  }
  
  const touch = e.touches[0]
  swipeState.startX = touch.clientX
  swipeState.startY = touch.clientY
  swipeState.currentX = touch.clientX
  swipeState.currentY = touch.clientY
  swipeState.startTime = Date.now()
  swipeState.isDragging = false
  swipeState.isHorizontalSwipe = false
  swipeState.isVerticalSwipe = false
  
  hideSwipeHint()
}



// 触摸移动处理 - 完全修复版本
const handleTouchMove = (e) => {
  if (!isMobile.value || !imagePreview.visible) return
  
  // 双指缩放处理
  if (e.touches.length === 2) {
    const currentDistance = getDistance(e.touches[0], e.touches[1])
    if (imagePreview.initialDistance > 0) {
      const scale = (currentDistance / imagePreview.initialDistance) * imagePreview.lastScale
      imagePreview.scale = Math.max(0.5, Math.min(3, scale))
    }
    return
  }
  
  // 放大状态下的单指拖拽
  if (imagePreview.scale > 1 && e.touches.length === 1) {
    const touch = e.touches[0]
    const deltaX = touch.clientX - swipeState.currentX
    const deltaY = touch.clientY - swipeState.currentY
    
    imagePreview.translateX += deltaX
    imagePreview.translateY += deltaY
    
    swipeState.currentX = touch.clientX
    swipeState.currentY = touch.clientY
    return
  }
  
  // 正常状态下的滑动处理
  const touch = e.touches[0]
  swipeState.currentX = touch.clientX
  swipeState.currentY = touch.clientY
  
  const deltaX = swipeState.currentX - swipeState.startX
  const deltaY = swipeState.currentY - swipeState.startY
  
  const minDistance = 10
  if (Math.abs(deltaX) < minDistance && Math.abs(deltaY) < minDistance) {
    return
  }
  
  if (!swipeState.isDragging) {
    swipeState.isDragging = true
    
    if (Math.abs(deltaX) > Math.abs(deltaY)) {
      swipeState.isHorizontalSwipe = true
      swipeState.isVerticalSwipe = false
    } else {
      swipeState.isVerticalSwipe = true
      swipeState.isHorizontalSwipe = false
    }
  }
  
  if (swipeState.isHorizontalSwipe && imagePreview.images.length > 1) {
    const screenWidth = window.innerWidth
    const maxTranslate = screenWidth * 0.4
    
    let translateX = deltaX
    if (Math.abs(deltaX) > maxTranslate) {
      const excess = Math.abs(deltaX) - maxTranslate
      const resistance = Math.min(excess / screenWidth, 0.5)
      translateX = deltaX > 0 ? 
        maxTranslate + excess * (1 - resistance) : 
        -maxTranslate - excess * (1 - resistance)
    }
    
    swipeState.translateX = translateX
    
    // const opacityFactor = Math.min(Math.abs(deltaX) / (screenWidth * 0.6), 0.4)
    // swipeState.opacity = 1 - opacityFactor
  }
  else if (swipeState.isVerticalSwipe) {
    const screenHeight = window.innerHeight
    swipeState.translateY = deltaY
    
    const scaleFactor = Math.max(0.7, 1 - Math.abs(deltaY) / (screenHeight * 1.5))
    swipeState.scale = scaleFactor
    
    // const opacityFactor = Math.min(Math.abs(deltaY) / (screenHeight * 0.8), 0.8)
    // swipeState.opacity = 1 - opacityFactor
  }
}



//从这时候开始修改
const imagePreview = reactive({
  visible: false,
  url: '',
  loading: false,
  currentIndex: 0,
  images: [],
  isTransitioning: false,
  direction: 0,
  // 新增缩放相关状态
  scale: 1,
  translateX: 0,
  translateY: 0,
  initialDistance: 0,
  lastScale: 1,
  lastTapTime: 0,
  rotation: 0, // 添加这个属性
  error: false  // 添加这个属性
})



// 新增辅助函数
const getDistance = (touch1, touch2) => {
  const dx = touch1.clientX - touch2.clientX
  const dy = touch1.clientY - touch2.clientY
  return Math.sqrt(dx * dx + dy * dy)
}

const getCenter = (touch1, touch2) => {
  return {
    x: (touch1.clientX + touch2.clientX) / 2,
    y: (touch1.clientY + touch2.clientY) / 2
  }
}

const handleDoubleTap = (e) => {
  const now = Date.now()
  const timeDiff = now - imagePreview.lastTapTime
  
  if (timeDiff < 300) {
    if (imagePreview.scale > 1) {
      imagePreview.scale = 1
      imagePreview.translateX = 0
      imagePreview.translateY = 0
    } else {
      imagePreview.scale = 2
      
      const touch = e.changedTouches[0]
      const rect = e.target.getBoundingClientRect()
      const centerX = rect.width / 2
      const centerY = rect.height / 2
      const tapX = touch.clientX - rect.left - centerX
      const tapY = touch.clientY - rect.top - centerY
      
      imagePreview.translateX = -tapX * 0.5
      imagePreview.translateY = -tapY * 0.5
    }
  }
  
  imagePreview.lastTapTime = now
}




// const switchImage = async (direction) => {
//   const newIndex = imagePreview.currentIndex + direction
//   if (newIndex < 0 || newIndex >= imagePreview.images.length || imagePreview.isTransitioning) {
//     return
//   }
//   imagePreview.isTransitioning = true
//   const screenWidth = window.innerWidth
//   try {
//     if (direction === 1) {
//       swipeState.translateX = screenWidth
//     } else if (direction === -1) {
//       swipeState.translateX = -screenWidth
//     } 
//     // 等待 Vue 更新 DOM
//     await nextTick()
//     // 等待退出动画完成
//     await new Promise(resolve => setTimeout(resolve, 300))
//     // 切换到新图片
//     imagePreview.currentIndex = newIndex
//     imagePreview.url = imagePreview.images[newIndex]
//     imagePreview.loading = true
//     imagePreview.error = false
//     await nextTick()
    
//     // 重置状态
//     swipeState.opacity = 1
//     swipeState.scale = 1
//     swipeState.translateY = 0
//     imagePreview.scale = 1
//     imagePreview.translateX = 0
//     imagePreview.translateY = 0
//     imagePreview.rotation = 0
//     swipeState.isDragging = false
//     swipeState.isHorizontalSwipe = false
//     swipeState.isVerticalSwipe = false
//     // 设置新图片初始位置（从屏幕外进入）
//     const initialPosition = direction === 1 ? -screenWidth : screenWidth
//     swipeState.translateX = initialPosition
//     swipeState.isDragging = true // 暂时禁用过渡
//     // 等待 Vue 更新 DOM
//     await nextTick()
//     // 触发进入动画
//     swipeState.translateX = 0
//     swipeState.isDragging = false
//     // 等待 Vue 更新 DOM
//     await nextTick()
//     // 等待进入动画完成
//     await new Promise(resolve => setTimeout(resolve, 300))
//   } catch (error) {
//     console.error('切换图片失败:', error)
//   } finally {
//     imagePreview.isTransitioning = false
//     hideSwipeHint()
//   }
// }


// const switchImageComputer = async (direction) => {
//   const newIndex = imagePreview.currentIndex + direction
//   if (newIndex < 0 || newIndex >= imagePreview.images.length || imagePreview.isTransitioning) {
//     return
//   }
//   imagePreview.isTransitioning = true
//   const screenWidth = window.innerWidth
//   try {
//     swipeState.isDragging = false
//     swipeState.isHorizontalSwipe = false
//     swipeState.isVerticalSwipe = false
//     swipeState.translateX = 0
//     swipeState.translateY = 0
//     swipeState.scale = 1
//     swipeState.opacity = 1
    
//     await nextTick()
//     // 第一阶段：旧图片退出动画
//     if (direction === 1) {
//       swipeState.translateX = screenWidth
//     } else if (direction === -1) {
//       swipeState.translateX = -screenWidth
//     }
//     // 等待退出动画完成
//     await new Promise(resolve => setTimeout(resolve, 300))
//     // 切换到新图片
//     imagePreview.currentIndex = newIndex
//     imagePreview.url = imagePreview.images[newIndex]
//     imagePreview.loading = true
//     imagePreview.error = false
    
//     await nextTick()
  
//     // 完全重置所有状态
//     swipeState.opacity = 1
//     swipeState.scale = 1
//     swipeState.translateY = 0
//     imagePreview.scale = 1
//     imagePreview.translateX = 0
//     imagePreview.translateY = 0
//     imagePreview.rotation = 0
//     swipeState.isDragging = false
//     swipeState.isHorizontalSwipe = false
//     swipeState.isVerticalSwipe = false

//     // 设置新图片初始位置（从屏幕外进入）
//     const initialPosition = direction === 1 ? -screenWidth : screenWidth
//     swipeState.translateX = initialPosition
//     swipeState.isDragging = true // 暂时禁用过渡
    
//     await nextTick()
    
//     // 小延迟确保状态已应用
//     await new Promise(resolve => setTimeout(resolve, 50))
//     // 触发进入动画
//     swipeState.translateX = 0
//     swipeState.isDragging = false
    
//     await nextTick()
    
//     // 等待进入动画完成
//     await new Promise(resolve => setTimeout(resolve, 300))
    
//   } catch (error) {
//     console.error('切换图片失败:', error)
//   } finally {
//     imagePreview.isTransitioning = false
//     hideSwipeHint()
//   }
// }


const switchImage = async (direction) => {
  const newIndex = imagePreview.currentIndex + direction
  if (newIndex < 0 || newIndex >= imagePreview.images.length || imagePreview.isTransitioning) {
    return
  }
  
  imagePreview.isTransitioning = true
  const screenWidth = window.innerWidth
  
  try {
    if (isMobile.value) {
      // 移动端逻辑（原 switchImage 逻辑）
      if (direction === 1) {
        swipeState.translateX = screenWidth
      } else if (direction === -1) {
        swipeState.translateX = -screenWidth
      } 
      
      await nextTick()
      //await new Promise(resolve => setTimeout(resolve, 10))
      
      // 切换到新图片
      imagePreview.currentIndex = newIndex
      imagePreview.url = imagePreview.images[newIndex]
      imagePreview.loading = true
      imagePreview.error = false
      await nextTick()
      
      // 重置状态
      swipeState.opacity = 1
      swipeState.scale = 1
      swipeState.translateY = 0
      imagePreview.scale = 1
      imagePreview.translateX = 0
      imagePreview.translateY = 0
      imagePreview.rotation = 0
      swipeState.isDragging = false
      swipeState.isHorizontalSwipe = false
      swipeState.isVerticalSwipe = false
      
      // 设置新图片初始位置（从屏幕外进入）
      const initialPosition = direction === 1 ? -screenWidth : screenWidth
      swipeState.translateX = initialPosition
      swipeState.isDragging = true // 暂时禁用过渡
      
      await nextTick()
      
      // 触发进入动画
      swipeState.translateX = 0
      swipeState.isDragging = false
      
      await nextTick()
      await new Promise(resolve => setTimeout(resolve, 300))
      
    } else {
      // 电脑端逻辑（原 switchImageComputer 逻辑）
      // 预先重置所有状态
      swipeState.isDragging = false
      swipeState.isHorizontalSwipe = false
      swipeState.isVerticalSwipe = false
      swipeState.translateX = 0
      swipeState.translateY = 0
      swipeState.scale = 1
      swipeState.opacity = 1
      
      await nextTick()
      
      // 第一阶段：旧图片退出动画
      if (direction === 1) {
        swipeState.translateX = screenWidth
      } else if (direction === -1) {
        swipeState.translateX = -screenWidth
      }
      
      await new Promise(resolve => setTimeout(resolve, 300))
      
      // 切换到新图片
      imagePreview.currentIndex = newIndex
      imagePreview.url = imagePreview.images[newIndex]
      imagePreview.loading = true
      imagePreview.error = false
      
      await nextTick()
    
      // 完全重置所有状态
      swipeState.opacity = 1
      swipeState.scale = 1
      swipeState.translateY = 0
      imagePreview.scale = 1
      imagePreview.translateX = 0
      imagePreview.translateY = 0
      imagePreview.rotation = 0
      swipeState.isDragging = false
      swipeState.isHorizontalSwipe = false
      swipeState.isVerticalSwipe = false

      // 设置新图片初始位置（从屏幕外进入）
      const initialPosition = direction === 1 ? -screenWidth : screenWidth
      swipeState.translateX = initialPosition
      swipeState.isDragging = true // 暂时禁用过渡
      
      await nextTick()
      
      // 小延迟确保状态已应用
      await new Promise(resolve => setTimeout(resolve, 50))
      
      // 触发进入动画
      swipeState.translateX = 0
      swipeState.isDragging = false
      
      await nextTick()
      
      // 等待进入动画完成
      await new Promise(resolve => setTimeout(resolve, 300))
    }
    
  } catch (error) {
    console.error('切换图片失败:', error)
  } finally {
    imagePreview.isTransitioning = false
    hideSwipeHint()
  }
}



const getImageTransform = () => {
  let transform = ''
  let transition = 'none'
  
  // 基础旋转
  const baseRotation = `rotate(${imagePreview.rotation || 0}deg)`
  
  if (imagePreview.scale > 1) {
    // 放大状态下的拖拽 - 使用 imagePreview 的位移
    transform = `${baseRotation} translate(${imagePreview.translateX}px, ${imagePreview.translateY}px) scale(${imagePreview.scale})`
    if (!swipeState.isDragging) {
      transition = 'transform 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94)'
    }
  } else {
    // 正常状态下的滑动和缩放 - 使用 swipeState 的位移
    transform = `${baseRotation} translate(${swipeState.translateX}px, ${swipeState.translateY}px) scale(${swipeState.scale})`
    transition = swipeState.isDragging ? 'none' : 'transform 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94)'
  }
  
  return {
    transform,
    opacity: swipeState.opacity,
    transition
  }
}




// 触摸结束处理 - 完全修复版本
const handleTouchEnd = (e) => {
  if (!isMobile.value || !imagePreview.visible) return
  
  // 双指缩放结束
  if (e.touches.length === 0 && imagePreview.initialDistance > 0) {
    imagePreview.lastScale = imagePreview.scale
    imagePreview.initialDistance = 0
    
    // 如果缩放太小，重置为1
    if (imagePreview.scale < 1) {
      imagePreview.scale = 1
      imagePreview.translateX = 0
      imagePreview.translateY = 0
      imagePreview.lastScale = 1
    }
  }
  
  // 单指操作
  if (!swipeState.isDragging && e.changedTouches.length === 1) {
    handleDoubleTap(e)
    return
  }
  
  // 如果是放大状态，不处理滑动切换
  if (imagePreview.scale > 1) {
    return
  }
  
  if (!swipeState.isDragging) {
    resetSwipeState()
    return
  }
  
  const deltaX = swipeState.currentX - swipeState.startX
  const deltaY = swipeState.currentY - swipeState.startY
  const deltaTime = Date.now() - swipeState.startTime
  
  const velocityX = Math.abs(deltaX) / deltaTime
  const velocityY = Math.abs(deltaY) / deltaTime
  
  const distanceThreshold = 80
  const velocityThreshold = 0.3
  
  if (swipeState.isHorizontalSwipe && imagePreview.images.length > 1) {
    const shouldSwitch = Math.abs(deltaX) > distanceThreshold || velocityX > velocityThreshold
    
    if (shouldSwitch) {
      const direction = deltaX > 0 ? -1 : 1
      const newIndex = imagePreview.currentIndex + direction
      
      if (newIndex >= 0 && newIndex < imagePreview.images.length) {
        switchImage(direction)
        return
      }
    }
  }
  else if (swipeState.isVerticalSwipe) {
    const shouldClose = Math.abs(deltaY) > distanceThreshold || velocityY > velocityThreshold
    
    if (shouldClose) {
      closePreview()
      return
    }
  }

  resetSwipeState()
}


// 键盘事件处理
const handleKeydown = (e) => {
  if (!imagePreview.visible) return
  
  switch (e.key) {
    case 'Escape':
      closePreview()
      break
    case 'ArrowLeft':
      if (imagePreview.currentIndex > 0) {
        switchImage(-1)
      }
      break
    case 'ArrowRight':
      if (imagePreview.currentIndex < imagePreview.images.length - 1) {
        switchImage(1)
      }
      break
  }
}

// 预加载图片
const preloadImages = (images, currentIndex) => {
  const preloadIndexes = []
  
  for (let i = Math.max(0, currentIndex - 2); i <= Math.min(images.length - 1, currentIndex + 2); i++) {
    if (i !== currentIndex) {
      preloadIndexes.push(i)
    }
  }
  
  preloadIndexes.forEach(index => {
    const img = new Image()
    img.src = images[index]
  })
}

// 预览图片加载完成
const onPreviewImageLoad = () => {
  imagePreview.loading = false
}

// 预览图片加载错误
const onPreviewImageError = () => {
  imagePreview.loading = false
}

// 获取背景透明度样式
const getBackgroundOpacity = () => {
  return {
    backgroundColor: `rgba(0, 0, 0, ${0.95 * swipeState.opacity})`
  }
}



// 窗口大小变化处理
const handleResize = () => {
  checkDeviceType()
}

// 组件挂载
onMounted(() => {
  checkDeviceType()
  loadPosts(true)
  window.addEventListener('scroll', throttledHandleScroll, { passive: true })
  window.addEventListener('resize', handleResize, { passive: true })
  document.addEventListener('keydown', handleKeydown)

  //强行修改上层组件样式
  if (window.innerWidth <= 768) {
    const cardElement = document.querySelector('.card-fade-up')
    if (cardElement) {
      cardElement.style.paddingLeft = '0'
      cardElement.style.paddingRight = '0'
      cardElement.style.paddingTop = '20px'
      cardElement.style.paddingBottom = '10px'
    }
  }

})

// 组件卸载
onUnmounted(() => {
  window.removeEventListener('scroll', throttledHandleScroll)
  window.removeEventListener('resize', handleResize)
  document.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = ''
  if (swipeHint.timer) {
    clearTimeout(swipeHint.timer)
  }
})
</script>

<template>
  <BannerPage title="印象" label="diary" :loading="loading">
    <!-- 动态列表容器 -->
    <div class="diary-container">
      <!-- 动态卡片 -->
      <div 
        v-for="post in state.posts" 
        :key="post.id"
        class="diary-card"
      >
        <!-- 卡片内容 -->
        <div class="card-content">
          <!-- 时间戳 -->
          <div class="timestamp-section">
            <div class="timestamp">
              <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ formatTime(post.add_time) }}
            </div>
          </div>

          <!-- 文字内容 -->
          <div class="text-content" v-if="post.content">
            <p class="content-text">
              {{ post.content }}
            </p>
          </div>

          <!-- 图片展示 -->
          <div 
            v-if="post.imgs && post.imgs.length > 0"
            class="image-grid"
            :class="getImageGridClass(post.imgs.length)"
          >
            <div 
              v-for="(image, index) in post.imgs.slice(0, 9)" 
              :key="index"
              class="image-item"
              :class="{
                'single-image': post.imgs.length === 1,
                'multi-image': post.imgs.length > 1
              }"
              @click="previewImage(convertImgUrl(image), post.imgs.map(img => convertImgUrl(img)), index)"
            >
              <img 
                :src="convertImgUrl(image)"
                :alt="`图片 ${index + 1}`"
                class="image"
                loading="lazy"
                @error="$event.target.src = '/placeholder-image.jpg'"
              />
              
              <!-- 多图数量指示器 -->
              <div 
                v-if="post.imgs.length > 9 && index === 8"
                class="image-overlay"
              >
                <span class="overlay-text">+{{ post.imgs.length - 9 }}</span>
              </div>
              
              <!-- 悬停遮罩 -->
              <div class="hover-overlay"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- 加载更多指示器 -->
      <div v-if="state.loading" class="loading-indicator">
        <div class="loading-content">
          <svg class="loading-spinner" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span>加载中...</span>
        </div>
      </div>

      <!-- 没有更多数据提示 -->
      <div v-else-if="!state.hasMore && state.posts.length > 0" class="no-more-hint">
        <p>没有更多内容了</p>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!state.loading && state.posts.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
        </div>
        <p class="empty-title">暂无动态</p>
        <p class="empty-desc">还没有发布任何动态</p>
      </div>
    </div>

    <!-- 图片预览弹窗 -->
    <Teleport to="body">
      <div 
        v-if="imagePreview.visible"
        class="image-preview-modal"
        :style="getBackgroundOpacity()"
        @click="closePreview"
        @touchstart.prevent="handleTouchStart"
        @touchmove.prevent="handleTouchMove"
        @touchend.prevent="handleTouchEnd"
      >
        <!-- 关闭按钮 -->
        <button 
           v-if="!isMobile"
          class="preview-close-btn"
          @click.stop="closePreview"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <!-- 图片计数器 -->
        <div 
          v-if="imagePreview.images.length > 1"
          class="image-counter"
        >
          {{ imagePreview.currentIndex + 1 }} / {{ imagePreview.images.length }}
        </div>

        <!-- 预览图片容器 -->
        <div 
          class="preview-image-container"
          @click.stop
        >
          <!-- 加载指示器 -->
          <div v-if="imagePreview.loading" class="preview-loading">
            <svg class="loading-spinner" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </div>

          <!-- 预览图片 -->
          <img 
            :src="imagePreview.url"
            class="preview-image"
            :style="getImageTransform()"
            @load="onPreviewImageLoad"
            @error="onPreviewImageError"
            @click.stop
          />
        </div>

        <!-- 左右切换按钮（电脑端显示） -->
        <template v-if="!isMobile && imagePreview.images.length > 1">
          <button 
            v-if="imagePreview.currentIndex > 0"
            class="nav-btn nav-btn-left"
            @click.stop="switchImage(-1)"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          
          <button 
            v-if="imagePreview.currentIndex < imagePreview.images.length - 1"
            class="nav-btn nav-btn-right"
            @click.stop="switchImage(1)"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </template>

        <!-- 滑动提示（仅移动端显示） -->
        <div 
          v-if="isMobile && swipeHint.visible && imagePreview.images.length > 1"
          class="swipe-hint"
        >
          <div class="swipe-hint-content">
            <span>滑动切换 拖拽关闭</span>
          </div>
        </div>
      </div>
    </Teleport>
  </BannerPage>
</template>

<style scoped>
/* 主容器 */
.diary-container {
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
  padding: 20px;
}

/* 日记卡片 */
.diary-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
  overflow: hidden;
  width: 100%;
  transition: all 0.3s ease;
}

.diary-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.card-content {
  padding: 20px;
}

/* 时间戳区域 */
.timestamp-section {
  margin-bottom: 12px;
}

.timestamp {
  display: flex;
  align-items: center;
  color: #666;
  font-size: 14px;
}

/* 文字内容 */
.text-content {
  margin-bottom: 16px;
}

.content-text {
  color: #1e1d1d;
  font-size: 20px;
  line-height: 1.6;
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* 图片网格 */
.image-grid {
  display: grid;
  gap: 8px;
  border-radius: 8px;
  overflow: hidden;
}

.grid-cols-1 {
  grid-template-columns: 1fr;
}

.grid-cols-2 {
  grid-template-columns: repeat(2, 1fr);
}

.grid-cols-3 {
  grid-template-columns: repeat(3, 1fr);
}

/* 图片项 */
.image-item {
  position: relative;
  cursor: pointer;
  border-radius: 6px;
  overflow: hidden;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.image-item.single-image {
  aspect-ratio: 16/9;
  max-height: 400px;
}

.image-item.multi-image {
  aspect-ratio: 1;
}

.image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
  display: block;
}

.image-item:hover .image {
  transform: scale(1.05);
}

/* 悬停遮罩 */
.hover-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0);
  transition: background-color 0.3s ease;
  pointer-events: none;
}

.image-item:hover .hover-overlay {
  background: rgba(0, 0, 0, 0.1);
}

/* 图片数量遮罩 */
.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
  font-size: 18px;
}

.overlay-text {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.5);
}

/* 加载指示器 */
.loading-indicator {
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px;
  color: #666;
}

.loading-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* 没有更多数据提示 */
.no-more-hint {
  text-align: center;
  padding: 40px;
  color: #999;
  font-size: 14px;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  text-align: center;
  color: #999;
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-title {
  font-size: 18px;
  font-weight: 500;
  margin: 0 0 8px 0;
  color: #666;
}

.empty-desc {
  font-size: 14px;
  margin: 0;
  color: #999;
}

/* 图片预览弹窗 */
.image-preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  transition: opacity 0.3s ease;
}

/* 关闭按钮 */
.preview-close-btn {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 50%;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  z-index: 10001;
  backdrop-filter: blur(10px);
}

.preview-close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.1);
}

/* 图片计数器 */
.image-counter {
  position: absolute;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.5);
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  z-index: 10001;
  backdrop-filter: blur(10px);
}

.preview-image-container {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  /* 移动端和电脑端使用不同的padding */
  padding: 60px 20px 40px 20px; /* 移动端默认 */
}


/* 预览图片 */
.preview-image {
  max-width: 100%;
  max-height: 100%;
  width: auto;
  height: auto;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  user-select: none;
  -webkit-user-drag: none;
}



/* 大屏电脑 */
@media (min-width: 1200px) {
  .preview-image-container {
    padding: 100px 150px 80px 150px;
  }
}


/* 预览加载指示器 */
.preview-loading {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  z-index: 10001;
}


/* 在 .preview-image 样式后添加 */
.preview-image {
  max-width: 100%;
  max-height: 100%;
  width: auto;
  height: auto;
  object-fit: contain;
  border-radius: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  user-select: none;
  -webkit-user-drag: none;
  /* 添加切换动画 */
  transition: transform 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
}

/* 切换动画状态 */
.preview-image.transitioning-out-left {
  transform: translateX(-100vw);
}

.preview-image.transitioning-out-right {
  transform: translateX(100vw);
}

.preview-image.transitioning-in-left {
  transform: translateX(-100vw);
}

.preview-image.transitioning-in-right {
  transform: translateX(100vw);
}


/* 导航按钮 */
.nav-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 50%;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  z-index: 10001;
  backdrop-filter: blur(10px);
}

.nav-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-50%) scale(1.1);
}

.nav-btn-left {
  left: 20px;
}

.nav-btn-right {
  right: 20px;
}

/* 滑动提示 */
.swipe-hint {
  position: absolute;
  bottom: 40px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(0, 0, 0, 0.7);
  color: white;
  padding: 12px 20px;
  border-radius: 25px;
  font-size: 14px;
  z-index: 10001;
  backdrop-filter: blur(10px);
  animation: fadeInOut 3s ease-in-out;
}

.swipe-hint-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

@keyframes fadeInOut {
  0%, 100% {
    opacity: 0;
    transform: translateX(-50%) translateY(10px);
  }
  10%, 90% {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}


/* 响应式设计 */
@media (max-width: 768px) {

  .diary-container {
    padding: 0px;
  }

  .diary-card {
    margin-bottom: 16px;
    border-radius: 8px;
  }

  .card-content {
    padding: 6px;
  }

  .content-text {
    font-size: 30px;
  }

  .timestamp {
    font-size: 20px;
  }

  .image-grid {
    gap: 6px;
  }

  .image-item.single-image {
    max-height: 300px;
  }

  /* 预览弹窗移动端适配 */
  .preview-close-btn {
    top: 10px;
    right: 10px;
    width: 36px;
    height: 36px;
  }

  .image-counter {
    top: 10px;
    font-size: 12px;
    padding: 4px 10px;
  }

  .preview-image-container {
    padding: 60px 5px 40px 5px; /* 进一步减少左右内边距 */
    max-width: 100vw;            /* 完全占满屏幕宽度 */
    max-height: 85vh;            /* 保留顶部和底部空间给按钮和提示 */
  }
  

  .swipe-hint {
    bottom: 30px;
    font-size: 13px;
    padding: 10px 16px;
  }

  .empty-state {
    padding: 60px 20px;
  }

  .empty-icon {
    width: 48px;
    height: 48px;
  }

  .empty-title {
    font-size: 16px;
  }

  .empty-desc {
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .diary-container {
    padding: 12px;
  }

  .card-content {
    padding: 12px;
  }

  .content-text {
    font-size: 14px;
  }

  .timestamp {
    font-size: 12px;
  }

  .image-grid {
    gap: 4px;
  }

  .image-item.single-image {
    max-height: 250px;
  }

  .loading-content {
    gap: 8px;
    font-size: 14px;
  }

  .loading-spinner {
    width: 20px;
    height: 20px;
  }
}

/* 深色模式支持 */
@media (prefers-color-scheme: dark) {
  .diary-card {
    background: #1a1a1a;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  }

  .diary-card:hover {
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.4);
  }

  .content-text {
    color: #e0e0e0;
  }

  .timestamp {
    color: #999;
  }

  .loading-indicator {
    color: #999;
  }

  .no-more-hint {
    color: #666;
  }

  .empty-state {
    color: #666;
  }

  .empty-title {
    color: #999;
  }

  .empty-desc {
    color: #666;
  }
}

/* 打印样式 */
@media print {
  .image-preview-modal {
    display: none !important;
  }

  .diary-card {
    break-inside: avoid;
    box-shadow: none;
    border: 1px solid #ddd;
  }

  .image-item {
    break-inside: avoid;
  }

  .loading-indicator,
  .no-more-hint {
    display: none;
  }
}

/* 高对比度模式支持 */
@media (prefers-contrast: high) {
  .diary-card {
    border: 2px solid #000;
  }

  .image-item {
    border: 1px solid #000;
  }

  .preview-close-btn,
  .nav-btn {
    border: 2px solid #fff;
  }
}

/* 减少动画模式支持 */
@media (prefers-reduced-motion: reduce) {
  .diary-card,
  .image,
  .preview-close-btn,
  .nav-btn,
  .image-preview-modal {
    transition: none;
  }

  .loading-spinner {
    animation: none;
  }

  .swipe-hint {
    animation: none;
  }
}

/* 焦点可见性 */
.preview-close-btn:focus-visible,
.nav-btn:focus-visible {
  outline: 2px solid #fff;
  outline-offset: 2px;
}

/* 触摸优化 */
@media (hover: none) and (pointer: coarse) {
  .diary-card:hover {
    transform: none;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .image-item:hover .image {
    transform: none;
  }

  .image-item:hover .hover-overlay {
    background: rgba(0, 0, 0, 0);
  }

  .preview-close-btn:hover,
  .nav-btn:hover {
    transform: none;
    background: rgba(255, 255, 255, 0.1);
  }

  .nav-btn:hover {
    transform: translateY(-50%);
  }
}

/* 滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* Firefox 滚动条 */
* {
  scrollbar-width: thin;
  scrollbar-color: #c1c1c1 #f1f1f1;
}
</style>
