import { defineStore } from 'pinia'
import { convertImgUrl } from '@/utils'
import api from '@/api'

export const useAppStore = defineStore('app', {
  state: () => ({
    searchFlag: false,
    loginFlag: false,
    registerFlag: false,
    changePasswordFlag: false,
    collapsed: false, // 侧边栏折叠（移动端）

    page_list: [], // 页面数据
    // TODO: 优化
    blogInfo: {
      article_count: 0,
      category_count: 0,
      tag_count: 0,
      view_count: 0,
      user_count: 0,
    },
    blog_config: {
      website_name: 'Breeze',
      website_author: 'Breeze',
      website_intro: '往事随风而去',
      website_avatar: '',
    },
  }),
  getters: {
    isMobile: () => !!navigator.userAgent.match(/(phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone)/i),
    articleCount: state => state.blogInfo.article_count ?? 0,
    categoryCount: state => state.blogInfo.category_count ?? 0,
    tagCount: state => state.blogInfo.tag_count ?? 0,
    viewCount: state => state.blogInfo.view_count ?? 0,
    pageList: state => state.page_list ?? [],
    blogConfig: state => state.blog_config,
  },
  actions: {
    setCollapsed(flag) { this.collapsed = flag },
    setLoginFlag(flag) { this.loginFlag = flag },
    setRegisterFlag(flag) { this.registerFlag = flag },
    setChangePasswordFlag(flag) { this.changePasswordFlag = flag },
    setSearchFlag(flag) { this.searchFlag = flag },

    async getBlogInfo() {
      try {
        const resp = await api.getHomeData()
        if (resp.code === 0) {
          this.blogInfo = resp.data
          this.blog_config = resp.data.blog_config
          this.blog_config.website_avatar = convertImgUrl(this.blog_config.website_avatar)
        }
        else {
          return Promise.reject(resp)
        }
      }
      catch (err) {
        return Promise.reject(err)
      }
    },

    async getPageList() {
      const resp = await api.getPageList()
      if (resp.code === 0) {
        this.page_list = resp.data
        this.page_list?.forEach(e => (e.cover = convertImgUrl(e.cover)))
      }
    },

    // 获取用户地理位置信息
    getCurrentLocation() {
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
    },


    // 上报位置信息
    async reportLocation() {
      try {
        // 获取位置信息
        const location = await this.getCurrentLocation()
        
        const reportData = {
          location: {
            latitude: location.latitude,
            longitude: location.longitude,
            accuracy: location.accuracy,
            timestamp: location.timestamp
          },
        }
        
        console.info('上报位置信息:', JSON.stringify(reportData, null, 2))
        
        // 发送上报请求
        fetch('/api/report', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(reportData)
        }).catch(error => {
          console.warn('位置上报失败:', error)
        })
        
      } catch (error) {
        console.warn('获取位置失败:', error.message)
        // 即使获取位置失败，也发送基本的上报信息
        const basicReportData = {
          location: null,
          error: error.message
        }
        
        fetch('/api/report', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(basicReportData)
        }).catch(reportError => {
          console.warn('基本上报也失败:', reportError)
        })
      }
    },
  },
})
