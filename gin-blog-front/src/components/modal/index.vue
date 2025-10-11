<script setup>
import SearchModal from '@/components/modal/SearchModal.vue'
import LoginModal from '@/components/modal/LoginModal.vue'
import RegisterModal from '@/components/modal/RegisterModal.vue'
import ChangePasswordModal from '@/components/modal/ChangePasswordModal.vue'
</script>

<template>
  <!-- 搜索窗口 -->
  <SearchModal />
  <!-- 登录窗口 -->
  <LoginModal />
  <!-- 注册窗口 -->
  <RegisterModal />
  <!-- 修改密码窗口 -->
  <ChangePasswordModal />
</template>
