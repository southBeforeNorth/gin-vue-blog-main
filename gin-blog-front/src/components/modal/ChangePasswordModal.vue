<script setup>
import { computed, ref, nextTick } from 'vue'
import UModal from '@/components/ui/UModal.vue'
import api from '@/api'
import { useAppStore } from '@/store'

const appStore = useAppStore()

const changePasswordFlag = computed({
  get: () => appStore.changePasswordFlag,
  set: val => appStore.setChangePasswordFlag(val),
})

const form = ref({
  email: '',
  code: '',
  newPassword: '',
  confirmPassword: '',
})

// 验证码倒计时
const countdown = ref(0)
const isCodeSending = ref(false)

// 发送验证码
async function sendCode() {
  const { email } = form.value
  
  // 邮箱格式校验
  const emailReg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/
  if (!emailReg.test(email)) {
    window.$message?.warning('请输入正确的邮箱格式')
    return
  }

  if (countdown.value > 0) {
    window.$message?.warning('请等待倒计时结束后再发送')
    return
  }

  try {
    isCodeSending.value = true
    console.log('准备发送验证码请求:', { email })
    const result = await api.sendchangePasswordCode({email})
    console.log('验证码发送成功:', result)
    window.$message?.success('验证码已发送，请查收邮箱')
    
    // 开始倒计时
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error) {
    console.error('验证码发送失败:', error)
    window.$message?.error('发送验证码失败，请稍后重试')
  } finally {
    isCodeSending.value = false
  }
}

// 修改密码
async function handleChangePassword() {
  console.log('开始修改密码流程')
  const { email, code, newPassword, confirmPassword } = form.value
  console.log('表单数据:', { email, code, newPassword: '***', confirmPassword: '***' })

  // 邮箱格式校验
  const emailReg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/
  if (!emailReg.test(email)) {
    console.log('邮箱格式校验失败')
    window.$message?.warning('请输入正确的邮箱格式')
    return
  }

  // 验证码校验
  if (!code) {
    window.$message?.warning('请输入验证码')
    return
  }

  if (code.length !== 6) {
    window.$message?.warning('验证码格式不正确')
    return
  }

  // 新密码校验
  if (!newPassword) {
    window.$message?.warning('请输入新密码')
    return
  }

  // 连续字符检测
  const hasConsecutive = /(.)\1{2,}/.test(newPassword) // 3个及以上相同字符
  if (hasConsecutive) {
    window.$message?.warning('密码不能包含3个及以上连续相同字符')
    return
  }

  // 密码长度校验
  if (newPassword.length < 6) {
    window.$message?.warning('密码长度不能少于6位')
    return
  }

  if (newPassword.length > 20) {
    window.$message?.warning('密码长度不能超过20位')
    return
  }

  // 密码复杂度校验
  const hasLetter = /[a-zA-Z]/.test(newPassword)
  const hasNumber = /[0-9]/.test(newPassword)
  if (!hasLetter) {
    window.$message?.warning('密码必须包含字母')
    return
  }

  if (!hasNumber) {
    window.$message?.warning('密码必须包含数字')
    return
  }

  // 弱密码检测
  const weakPasswords = [
    '123456', '123456789', 'password', '12345678', 'qwerty', 
    'abc123', '111111', '123123', 'admin', 'root'
  ]
  
  if (weakPasswords.includes(newPassword.toLowerCase())) {
    window.$message?.warning('密码过于简单，请使用更复杂的密码')
    return
  }


  // 确认密码校验
  if (!confirmPassword) {
    window.$message?.warning('请确认新密码')
    return
  }

  if (newPassword !== confirmPassword) {
    window.$message?.warning('两次输入的密码不一致')
    return
  }

  try {
    console.log('准备发送修改密码请求')
    const result = await api.changePassword({ email, code, newPassword })
    console.log('修改密码请求成功:', result)
    window.$message?.success('密码修改成功，请重新登录')
    form.value = { email: '', code: '', newPassword: '', confirmPassword: '' }
    
    // 直接切换模态框
    console.log('准备切换到登录页面')
    appStore.setChangePasswordFlag(false)
    // 使用 nextTick 确保 DOM 更新后再打开登录框
    nextTick(() => {
      console.log('打开登录模态框')
      appStore.setLoginFlag(true)
    })
  } catch (error) {
    console.error('修改密码请求失败:', error)
  }
}

// 返回登录
function openLogin() {
  console.log('执行 openLogin 函数')
  console.log('当前状态 - changePasswordFlag:', appStore.changePasswordFlag)
  console.log('当前状态 - loginFlag:', appStore.loginFlag)
  
  appStore.setChangePasswordFlag(false)
  appStore.setLoginFlag(true)
  
  console.log('设置后状态 - changePasswordFlag:', appStore.changePasswordFlag)
  console.log('设置后状态 - loginFlag:', appStore.loginFlag)
}
</script>

<template>
  <UModal v-model="changePasswordFlag" :width="480">
    <div class="mx-2 my-1">
      <div class="mb-4 text-xl font-bold">
        修改密码
      </div>
      <div class="my-7 space-y-4">
        <div class="flex items-center">
          <span class="mr-4 inline-block w-16 text-right"> 邮箱 </span>
          <input
            v-model="form.email" required placeholder="请输入邮箱地址"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
        </div>
        <div class="flex items-center">
          <span class="mr-4 inline-block w-16 text-right"> 验证码 </span>
          <div class="flex w-full gap-2">
            <input
              v-model="form.code" required placeholder="请输入验证码"
              class="flex-1 border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
            >
            <button
              :disabled="countdown > 0 || isCodeSending"
              class="px-4 py-2 text-sm rounded-md border border-blue-500 text-white bg-blue-500 hover:bg-blue-600 disabled:opacity-50 disabled:cursor-not-allowed disabled:bg-gray-400 disabled:border-gray-400"
              @click="sendCode"
            >
              {{ countdown > 0 ? `${countdown}s` : '发送' }}
            </button>
          </div>
        </div>
        <div class="flex items-center">
          <span class="mr-4 inline-block w-16 text-right"> 新密码 </span>
          <input
            v-model="form.newPassword" required type="password" placeholder="请输入新密码"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
        </div>
        <div class="flex items-center">
          <span class="mr-4 inline-block w-16 text-right"> 确认密码 </span>
          <input
            v-model="form.confirmPassword" required type="password" placeholder="请再次输入新密码"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
        </div>
      </div>
      <div class="my-2 text-center">
        <button
          class="w-full rounded bg-blue py-2 text-white hover:bg-light-blue"
          @click="handleChangePassword"
        >
          修改密码
        </button>
        <div class="mb-2 mt-6 text-left">
          <button class="duration-300 hover:text-emerald" @click="openLogin">
            想起密码了？返回登录
          </button>
        </div>
      </div>
    </div>
  </UModal>
</template>
