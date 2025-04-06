<script setup lang="ts">
import { useRouter } from 'vue-router';

const router = useRouter();

// 返回首页
const goHome = () => {
  router.push('/');
};
</script>

<template>
  <div class="not-found">
    <div class="content">
      <h1>404</h1>
      <div class="desc">抱歉，您访问的页面不存在</div>
      <el-button type="primary" @click="goHome">返回首页</el-button>
    </div>
  </div>
</template>

<style scoped>
.not-found {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
}

.content {
  text-align: center;
}

h1 {
  font-size: 72px;
  color: #409eff;
  margin-bottom: 20px;
}

.desc {
  font-size: 18px;
  color: #606266;
  margin-bottom: 30px;
}
</style>