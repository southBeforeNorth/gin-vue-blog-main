// 相对图片地址 => 完整的图片路径, 用于本地文件上传
// - 如果包含 http 说明是 Web 图片资源
// - 否则是服务器上的图片，需要拼接服务器路径
const SERVER_URL = import.meta.env.VITE_BACKEND_URL

/**
 * 将相对地址转换为完整的图片路径
 * @param {string} imgUrl
 * @returns {string} 完整的图片路径
 */
export function convertImgUrl(imgUrl) {
  if (!imgUrl) {
    return 'http://dummyimage.com/400x400'
  }
  // 网络资源
  if (imgUrl.startsWith('http')) {
    return imgUrl
  }
  // 服务器资源
  return `${SERVER_URL}/${imgUrl}`
}

/**
 * 获取用户地理位置信息
 * @returns {Promise<Object>} 包含经纬度等位置信息的对象
 */
export function getCurrentLocation() {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      reject(new Error('浏览器不支持地理位置功能'))
      return
    }

    const options = {
      enableHighAccuracy: true, // 启用高精度
      timeout: 10000, // 10秒超时
      maximumAge: 300000 // 5分钟缓存
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude, accuracy } = position.coords
        resolve({
          latitude,
          longitude,
          accuracy,
          timestamp: position.timestamp
        })
      },
      (error) => {
        let errorMessage = '获取位置失败'
        switch (error.code) {
          case error.PERMISSION_DENIED:
            errorMessage = '用户拒绝了位置访问请求'
            break
          case error.POSITION_UNAVAILABLE:
            errorMessage = '位置信息不可用'
            break
          case error.TIMEOUT:
            errorMessage = '获取位置超时'
            break
        }
        reject(new Error(errorMessage))
      },
      options
    )
  })
}

export * from './local'
export * from './http'
