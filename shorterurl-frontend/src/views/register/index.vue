<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/store/user';
import type { UserRegisterReq } from '@/api/user';
import { User, Lock, Message, Phone } from '@element-plus/icons-vue';

// 获取用户状态管理
const userStore = useUserStore();
// 获取路由
const router = useRouter();

// 注册表单数据
const registerForm = reactive<UserRegisterReq>({
  username: '',
  password: '',
  realname: '',
  phone: '',
  mail: ''
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
  ],
  realname: [
    { required: true, message: '请输入真实姓名', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  mail: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
};

// 表单ref
const registerFormRef = ref();
// 注册中状态
const loading = ref(false);

// 处理注册
const handleRegister = async () => {
  // 表单验证
  await registerFormRef.value.validate();
  
  // 设置加载状态
  loading.value = true;
  
  try {
    // 调用注册接口
    await userStore.register(registerForm);
    // 不再显示成功消息，由store中处理
    
    // 跳转到登录页
    router.push('/login');
  } catch (error: any) {
    ElMessage.error(error.message || '注册失败，请重试');
  } finally {
    loading.value = false;
  }
};

// 前往登录页
const goToLogin = () => {
  router.push('/login');
};
</script>

<template>
  <div class="register-container">
    <!-- 左侧宣传区域 -->
    <div class="register-banner">
      <div class="banner-header">
        <div class="logo-wrapper">
          <!-- 现代Logo -->
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
            
            <!-- 用户图标 -->
            <circle cx="200" cy="100" r="40" fill="#16213e" stroke="#7b68ee" stroke-width="3" />
            <circle cx="200" cy="85" r="15" fill="#7b68ee" />
            <path d="M160,180 C160,140 240,140 240,180" stroke="#7b68ee" fill="none" stroke-width="3" />
            
            <!-- 装饰元素 -->
            <circle cx="120" cy="200" r="25" fill="#e94560" opacity="0.7" />
            <circle cx="280" cy="220" r="20" fill="#7b68ee" opacity="0.5" />
            <rect x="250" y="80" width="60" height="40" rx="5" fill="#16213e" stroke="#e94560" stroke-width="2" />
            <line x1="260" y1="90" x2="300" y2="90" stroke="#e94560" stroke-width="2" />
            <line x1="260" y1="100" x2="290" y2="100" stroke="#e94560" stroke-width="2" />
            <line x1="260" y1="110" x2="280" y2="110" stroke="#e94560" stroke-width="2" />
          </svg>
        </div>
        
        <div class="banner-text">
          <h1 class="banner-title">加入 LinkPro</h1>
          <h2 class="banner-subtitle">开启您的短链接管理之旅</h2>
          <p class="banner-description">
            创建账号后，您将获得完整的短链接管理体验<br>
            数据分析、访问统计、链接优化，尽在掌握
          </p>
        </div>
      </div>
    </div>
    
    <!-- 右侧注册表单 -->
    <div class="register-form-container">
      <div class="register-form-wrapper">
        <h2 class="form-title">创建账号</h2>
        
        <el-form
          ref="registerFormRef"
          :model="registerForm"
          :rules="rules"
          class="register-form"
          label-position="top"
          hide-required-asterisk
        >
          <el-form-item prop="username">
            <el-input 
              v-model="registerForm.username" 
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
              v-model="registerForm.password"
              type="password"
              placeholder="请输入密码"
              size="large"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="realname">
            <el-input
              v-model="registerForm.realname"
              placeholder="请输入真实姓名"
              size="large"
            >
              <template #prefix>
                <el-icon><User /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="phone">
            <el-input
              v-model="registerForm.phone"
              placeholder="请输入手机号"
              size="large"
            >
              <template #prefix>
                <el-icon><Phone /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-form-item prop="mail">
            <el-input
              v-model="registerForm.mail"
              placeholder="请输入邮箱"
              size="large"
            >
              <template #prefix>
                <el-icon><Message /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          
          <el-button 
            type="primary" 
            :loading="loading" 
            @click="handleRegister"
            class="register-button"
            size="large"
          >
            创建账号
          </el-button>
          
          <div class="agreement-text">
            注册即表示您同意 <a href="javascript:void(0)">《服务条款》</a> 和 <a href="javascript:void(0)">《隐私政策》</a>
          </div>
            
          <div class="login-link">
            已有账号? <a href="javascript:void(0)" @click="goToLogin">返回登录</a>
          </div>
        </el-form>
      </div>
      
      <div class="register-footer">
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
.register-container {
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
.register-banner {
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

/* 右侧注册表单样式 */
.register-form-container {
  width: 500px;
  background-color: #fff;
  display: flex;
  flex-direction: column;
  padding: 40px 50px;
  box-shadow: -4px 0 30px rgba(0, 0, 0, 0.1);
  overflow: auto;
}

.register-form-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.form-title {
  font-size: 32px;
  font-weight: 800;
  color: #1a1a2e;
  margin: 0 0 40px;
  text-align: center;
}

.register-form {
  width: 100%;
}

.register-button {
  width: 100%;
  height: 52px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 8px;
  background: linear-gradient(90deg, #7b68ee, #e94560);
  border: none;
  margin: 16px 0 24px;
  transition: all 0.3s;
  letter-spacing: 1px;
}

.register-button:hover {
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

.login-link {
  font-size: 15px;
  color: #444;
  text-align: center;
  margin-bottom: 20px;
  font-weight: 500;
}

.login-link a {
  color: #7b68ee;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.3s;
}

.login-link a:hover {
  color: #e94560;
}

.register-footer {
  text-align: center;
  font-size: 14px;
  color: #666;
  margin-top: 20px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .register-form-container {
    width: 480px;
    padding: 40px;
  }
}

@media (max-width: 992px) {
  .register-container {
    flex-direction: column;
  }
  
  .register-banner {
    display: none;
  }
  
  .register-form-container {
    width: 100%;
    height: 100%;
    padding: 40px;
  }
}

@media (max-width: 480px) {
  .register-form-container {
    padding: 30px 20px;
  }
  
  .form-title {
    font-size: 28px;
    margin-bottom: 30px;
  }
}

.illustration-svg {
  width: 100%;
  height: auto;
  max-height: 350px;
}
</style> 