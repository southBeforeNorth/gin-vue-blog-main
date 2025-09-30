<script setup>
import { computed, ref } from 'vue'
import config from '@/assets/config'

import UModal from '@/components/ui/UModal.vue'
import api from '@/api'
import { useAppStore } from '@/store'

const appStore = useAppStore()

const registerFlag = computed({
  get: () => appStore.registerFlag,
  set: val => appStore.setRegisterFlag(val),
})

const form = ref({
  email: '', // 修改为 email
  password: '',
})

// // 注册
// async function handleRegister() {
//   const { email, password } = form.value

//   const reg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/
//   if (!reg.test(email)) {
//     window.$message?.warning('请输入正确的邮箱格式')
//     return
//   }

//   if (!password) {
//     window.$message?.warning('请输入密码')
//     return
//   }


// }

async function handleRegister() {
  const { email, password } = form.value

  // 邮箱格式校验
  const emailReg = /^([a-zA-Z]|[0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/
  if (!emailReg.test(email)) {
    window.$message?.warning('请输入正确的邮箱格式')
    return
  }

  // 密码基本校验
  if (!password) {
    window.$message?.warning('请输入密码')
    return
  }

  // 密码长度校验
  if (password.length < 6) {
    window.$message?.warning('密码长度不能少于6位')
    return
  }

  if (password.length > 20) {
    window.$message?.warning('密码长度不能超过20位')
    return
  }

  // 密码复杂度校验
  const hasLetter = /[a-zA-Z]/.test(password)
  const hasNumber = /[0-9]/.test(password)
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
  
  if (weakPasswords.includes(password.toLowerCase())) {
    window.$message?.warning('密码过于简单，请使用更复杂的密码')
    return
  }
  // 连续字符检测
  const hasConsecutive = /(.)\1{2,}/.test(password) // 3个及以上相同字符
  if (hasConsecutive) {
    window.$message?.warning('密码不能包含3个及以上连续相同字符')
    return
  }

  // 发送注册请求
  await api.register({ email, password })
  window.$message?.success('邮件已发送，请在邮箱中确认以完成注册')
  form.value = { email: '', password: '' }
}



// 登录
function openLogin() {
  appStore.setRegisterFlag(false)
  appStore.setLoginFlag(true)
}
</script>

<template>
  <UModal v-model="registerFlag" :width="480">
    <div class="mx-2 my-1">
      <div class="mb-4 text-xl font-bold">
        注册
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
          <span class="mr-4 inline-block w-16 text-right"> 密码 </span>
          <input
            v-model="form.password" required type="password" placeholder="请输入密码"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
        </div>
      </div>
      <div class="my-2 text-center">
        <button
          class="w-full rounded bg-red py-2 text-white hover:bg-orange"
          @click="handleRegister"
        >
          注册
        </button>
        <div class="mb-2 mt-6 text-left">
          已有账号？
          <button class="duration-300 hover:text-emerald" @click="openLogin">
            登录
          </button>
        </div>
      </div>
    </div>
  </UModal>
</template>