<script setup lang="ts">
import { onMounted } from 'vue';
import { useUserStore } from '@/store/user';

const userStore = useUserStore();

onMounted(async () => {
  // 检查登录状态
  if (userStore.isLogin) {
    await userStore.checkLogin();
  }
});
</script>

<template>
  <router-view />
</template>

<style scoped>
</style>
