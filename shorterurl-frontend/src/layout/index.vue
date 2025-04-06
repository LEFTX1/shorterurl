<template>
  <div class="layout-container">
    <!-- 左侧导航栏 -->
    <div class="sidebar">
      <div class="logo-container">
        <div class="logo-wrapper">
          <!-- 简洁Logo -->
          <svg class="logo-svg" viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
            <circle cx="20" cy="20" r="20" fill="#1a1a2e" />
            <path d="M12,16 L28,16 L20,28 L12,16" fill="#7b68ee" />
            <circle cx="20" cy="12" r="4" fill="#e94560" />
          </svg>
          <span class="logo-text">LinkPro</span>
        </div>
      </div>
      
      <div class="nav-menu">
        <div 
          class="nav-item" 
          :class="{ active: activeMenu === 'link' }"
          @click="navigateTo('/link')"
        >
          <i class="nav-icon">🔗</i>
          <span class="nav-text">短链管理</span>
        </div>
        
        <div 
          class="nav-item" 
          :class="{ active: activeMenu === 'recycle' }"
          @click="showRecycleBin"
        >
          <i class="nav-icon">🗑️</i>
          <span class="nav-text">回收站</span>
        </div>
        
        <div 
          class="nav-item" 
          :class="{ active: activeMenu === 'stats' }"
          @click="navigateTo('/stats')"
        >
          <i class="nav-icon">📊</i>
          <span class="nav-text">统计分析</span>
        </div>
        
        <div 
          class="nav-item" 
          :class="{ active: activeMenu === 'account' }"
          @click="navigateTo('/account')"
        >
          <i class="nav-icon">👤</i>
          <span class="nav-text">账号管理</span>
        </div>
      </div>
      
      <div class="user-section">
        <div class="user-info">
          <div class="user-avatar">{{ userInitial }}</div>
          <div class="user-name">{{ username }}</div>
        </div>
        <div class="logout-btn" @click="handleLogout">
          <i class="logout-icon">🚪</i>
          <span>退出登录</span>
        </div>
      </div>
    </div>
    
    <!-- 主内容区域 -->
    <div class="main-content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useUserStore } from '../store/user';
import mitt from 'mitt';

// 创建事件总线
export const emitter = mitt();

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

// 标记是否显示回收站
const isRecycleBin = ref(false);

const username = computed(() => userStore.getUsername || '用户');
const userInitial = computed(() => {
  return username.value.charAt(0).toUpperCase();
});

// 当前激活的菜单项
const activeMenu = computed(() => {
  const path = route.path;
  
  if (path.includes('/link')) {
    return isRecycleBin.value ? 'recycle' : 'link';
  }
  if (path.includes('/stats')) return 'stats';
  if (path.includes('/account')) return 'account';
  
  return '';
});

// 导航到指定路由
const navigateTo = (path: string) => {
  isRecycleBin.value = false;
  router.push(path);
};

// 显示回收站
const showRecycleBin = () => {
  isRecycleBin.value = true;
  router.push('/link');
  // 发送事件到link组件，切换到回收站视图
  emitter.emit('show-recycle-bin');
};

// 处理退出登录
const handleLogout = async () => {
  await userStore.logout();
};

onMounted(() => {
  // 监听回收站状态变化
  emitter.on('recycle-bin-status', (status: boolean) => {
    isRecycleBin.value = status;
  });
});
</script>

<style scoped>
.layout-container {
  display: flex;
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}

.sidebar {
  width: 220px;
  background-color: white;
  border-right: 1px solid #eaeaea;
  display: flex;
  flex-direction: column;
  height: 100%;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.05);
  z-index: 10;
}

.logo-container {
  padding: 20px 16px;
  border-bottom: 1px solid #f0f0f0;
}

.logo-wrapper {
  display: flex;
  align-items: center;
}

.logo-svg {
  width: 28px;
  height: 28px;
  margin-right: 10px;
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  color: #333;
  letter-spacing: 0.5px;
}

.nav-menu {
  flex: 1;
  padding: 16px 0;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 4px;
  border-radius: 4px;
  margin: 0 8px 4px 8px;
}

.nav-item:hover {
  background-color: #f5f5f5;
}

.nav-item.active {
  background-color: #e6f7ff;
  color: #1677ff;
}

.nav-icon {
  font-size: 16px;
  margin-right: 12px;
}

.nav-text {
  font-size: 14px;
  font-weight: 500;
}

.user-section {
  padding: 16px;
  border-top: 1px solid #f0f0f0;
}

.user-info {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background-color: #1677ff;
  color: white;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: 500;
  margin-right: 12px;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
}

.logout-btn {
  display: flex;
  align-items: center;
  padding: 8px;
  cursor: pointer;
  color: #f5222d;
  font-size: 14px;
  border-radius: 4px;
  transition: all 0.3s;
}

.logout-btn:hover {
  background-color: #fff1f0;
}

.logout-icon {
  margin-right: 8px;
}

.main-content {
  flex: 1;
  overflow: hidden;
  background-color: #f5f5f5;
}
</style> 