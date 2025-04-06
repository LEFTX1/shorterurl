<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/store/user';
import type { UserLoginReq } from '@/api/user';
import { User, Lock } from '@element-plus/icons-vue';

// 获取用户状态管理
const userStore = useUserStore();
// 获取路由
const router = useRouter();

// 登录表单数据
const loginForm = reactive<UserLoginReq>({
  username: '',
  password: ''
});

// 表单校验规则
const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 2, max: 32, message: '用户名长度在2到32个字符之间', trigger: 'blur' },
    { 
      validator: (rule: any, value: string, callback: (error?: Error) => void) => {
        if (value && !/^[\x00-\x7F]+$/.test(value)) {
          callback(new Error('用户名只能包含ASCII字符，不能使用中文'));
        } else {
          callback();
        }
      }, 
      trigger: 'blur' 
    }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 32, message: '密码长度在6到32个字符之间', trigger: 'blur' }
  ]
};

// 表单ref
const loginFormRef = ref();
// 登录中状态
const loading = ref(false);
// 记住密码
const rememberMe = ref(false);

// 处理登录
const handleLogin = async () => {
  // 表单验证
  await loginFormRef.value.validate();
  
  // 设置加载状态
  loading.value = true;
  
  try {
    // 调用登录接口
    await userStore.login(loginForm);
    ElMessage.success('登录成功');
    
    // 跳转到短链接管理页面
    router.push('/link');
  } catch (error: any) {
    ElMessage.error(error.message || '登录失败，请重试');
  } finally {
    loading.value = false;
  }
};

// 前往注册页
const goToRegister = () => {
  router.push('/register');
};

// 忘记密码
const forgotPassword = () => {
  ElMessage.info('忘记密码功能开发中...');
};
</script>

<template>
  <div class="login-container">
    <!-- 左侧宣传区域 -->
    <div class="login-banner">
      <div class="banner-header">
        <div class="logo-wrapper">
          <!-- 更现代的Logo -->
          <svg class="logo-svg" viewBox="0 0 40 40" xmlns="http://www.w3.org/2000/svg">
            <circle cx="20" cy="20" r="20" fill="#1a1a2e" />
            <path d="M12,16 L28,16 L20,28 L12,16" fill="#7b68ee" />
            <circle cx="20" cy="12" r="4" fill="#e94560" />
          </svg>
          <span class="logo-text">LinkPro</span>
        </div>
        
        <div class="banner-actions">
          <a href="javascript:void(0)" class="action-link">功能介绍</a>
          <a href="javascript:void(0)" class="action-link">帮助文档</a>
        </div>
      </div>
      
      <div class="banner-content">
        <div class="banner-illustration">
          <!-- 科技感插图 -->
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 300" class="illustration-svg">
            <!-- 背景网格 -->
            <pattern id="grid" width="20" height="20" patternUnits="userSpaceOnUse">
              <path d="M 20 0 L 0 0 0 20" fill="none" stroke="#2a2a42" stroke-width="0.5" />
            </pattern>
            <rect width="400" height="300" fill="#1a1a2e" />
            <rect width="400" height="300" fill="url(#grid)" />
            
            <!-- 数据可视化元素 -->
            <path d="M100,250 L100,150 L140,180 L180,120 L220,160 L260,100 L300,130" 
                  stroke="#e94560" fill="none" stroke-width="3" stroke-linecap="round" />
            <!-- 圆点标记 -->
            <circle cx="100" cy="150" r="5" fill="#e94560" />
            <circle cx="140" cy="180" r="5" fill="#e94560" />
            <circle cx="180" cy="120" r="5" fill="#e94560" />
            <circle cx="220" cy="160" r="5" fill="#e94560" />
            <circle cx="260" cy="100" r="5" fill="#e94560" />
            <circle cx="300" cy="130" r="5" fill="#e94560" />
            
            <!-- 链接图标 -->
            <g transform="translate(170, 70) scale(0.8)">
              <path d="M30,60 L60,30 C75,15 100,15 115,30 C130,45 130,70 115,85 L85,115" 
                    stroke="#7b68ee" fill="none" stroke-width="10" stroke-linecap="round" />
              <path d="M115,30 L85,0" stroke="#7b68ee" fill="none" stroke-width="10" stroke-linecap="round" />
              <path d="M90,60 L60,90 C45,105 20,105 5,90 C-10,75 -10,50 5,35 L35,5" 
                    stroke="#16213e" fill="none" stroke-width="10" stroke-linecap="round" />
              <path d="M5,90 L35,120" stroke="#16213e" fill="none" stroke-width="10" stroke-linecap="round" />
            </g>
            
            <!-- 装饰线条 -->
            <line x1="50" y1="50" x2="50" y2="250" stroke="#2a2a42" stroke-width="1" />
            <line x1="50" y1="250" x2="350" y2="250" stroke="#2a2a42" stroke-width="1" />
          </svg>
        </div>
        
        <div class="banner-text">
          <h1 class="banner-title">精准短链 · 数据驱动</h1>
          <h2 class="banner-subtitle">短链接管理与分析工具</h2>
          <p class="banner-description">
            简洁高效的短链接服务，为您提供链接管理、数据统计<br>
            和访问分析，让每个链接都能发挥最大价值
          </p>
        </div>
      </div>
    </div>
    
    <!-- 右侧登录表单 -->
    <div class="login-form-container">
      <div class="login-form-wrapper">
        <h2 class="form-title">用户登录</h2>
        
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="rules"
          class="login-form"
          label-position="top"
          hide-required-asterisk
        >
          <el-form-item prop="username">
            <el-input 
              v-model="loginForm.username" 
              placeholder="请输入用户名"
              size="large"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              @keyup.enter="handleLogin"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <div class="remember-forgot">
            <el-checkbox v-model="rememberMe">记住密码</el-checkbox>
            <a href="javascript:void(0)" class="forgot-password" @click="forgotPassword">忘记密码?</a>
          </div>
          
          <el-button 
            type="primary" 
            :loading="loading" 
            @click="handleLogin"
            class="login-button"
            size="large"
          >
            登录系统
          </el-button>
          
          <div class="agreement-text">
            登录即表示您同意 <a href="javascript:void(0)">《服务条款》</a> 和 <a href="javascript:void(0)">《隐私政策》</a>
          </div>
            
          <div class="register-link">
            首次使用? <a href="javascript:void(0)" @click="goToRegister">创建账号</a>
          </div>
        </el-form>
      </div>
      
      <div class="login-footer">
        <p>© {{ new Date().getFullYear() }} LinkPro - 短链接管理平台</p>
      </div>
    </div>
  </div>
</template>

<style>
/* 全局样式，确保页面填满 */
html, body, #app {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
</style>

<style scoped>
.login-container {
  display: flex;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  margin: 0;
  padding: 0;
}

/* 左侧宣传区域样式 */
.login-banner {
  flex: 1;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  padding: 40px;
  display: flex;
  flex-direction: column;
  position: relative;
  overflow: auto;
}

.banner-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 60px;
}

.logo-wrapper {
  display: flex;
  align-items: center;
}

.logo-svg {
  width: 36px;
  height: 36px;
  margin-right: 12px;
}

.logo-text {
  font-size: 24px;
  font-weight: 700;
  color: #ffffff;
  letter-spacing: 0.5px;
}

.banner-actions {
  display: flex;
  gap: 30px;
}

.action-link {
  color: #a0a0c0;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.3s;
  font-weight: 500;
}

.action-link:hover {
  color: #e94560;
}

.banner-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.banner-illustration {
  width: 80%;
  max-width: 600px;
  margin-bottom: 60px;
}

.banner-text {
  text-align: center;
}

.banner-title {
  font-size: 36px;
  font-weight: 800;
  color: #ffffff;
  margin: 0 0 16px;
  letter-spacing: 1px;
}

.banner-subtitle {
  font-size: 24px;
  font-weight: 600;
  color: #7b68ee;
  margin: 0 0 24px;
}

.banner-description {
  font-size: 16px;
  color: #a0a0c0;
  line-height: 1.8;
}

/* 右侧登录表单样式 */
.login-form-container {
  width: 480px;
  background-color: #fff;
  display: flex;
  flex-direction: column;
  padding: 60px 50px;
  box-shadow: -4px 0 30px rgba(0, 0, 0, 0.1);
  overflow: auto;
}

.login-form-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.form-title {
  font-size: 32px;
  font-weight: 800;
  color: #1a1a2e;
  margin: 0 0 50px;
  text-align: center;
}

.login-form {
  width: 100%;
}

.remember-forgot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 16px 0 30px;
}

.forgot-password {
  color: #7b68ee;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.3s;
  font-weight: 500;
}

.forgot-password:hover {
  color: #e94560;
}

.login-button {
  width: 100%;
  height: 52px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 8px;
  background: linear-gradient(90deg, #7b68ee, #e94560);
  border: none;
  margin: 0 0 24px;
  transition: all 0.3s;
  letter-spacing: 1px;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(123, 104, 238, 0.3);
}

.agreement-text {
  font-size: 13px;
  color: #666;
  text-align: center;
  margin-bottom: 24px;
}

.agreement-text a {
  color: #7b68ee;
  text-decoration: none;
  transition: color 0.3s;
}

.agreement-text a:hover {
  color: #e94560;
}

.register-link {
  font-size: 15px;
  color: #444;
  text-align: center;
  margin-bottom: 20px;
  font-weight: 500;
}

.register-link a {
  color: #7b68ee;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.3s;
}

.register-link a:hover {
  color: #e94560;
}

.login-footer {
  text-align: center;
  font-size: 14px;
  color: #666;
  margin-top: 20px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .login-form-container {
    width: 450px;
    padding: 50px 40px;
  }
}

@media (max-width: 992px) {
  .login-container {
    flex-direction: column;
  }
  
  .login-banner {
    display: none;
  }
  
  .login-form-container {
    width: 100%;
    height: 100%;
    padding: 60px 40px;
  }
}

@media (max-width: 480px) {
  .login-form-container {
    padding: 40px 20px;
  }
  
  .form-title {
    font-size: 28px;
    margin-bottom: 40px;
  }
}

.illustration-svg {
  width: 100%;
  height: auto;
  max-height: 350px;
}
</style> 