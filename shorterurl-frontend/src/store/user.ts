import { defineStore } from 'pinia';
import { ElMessage } from 'element-plus';
import router from '../router';
import userApi from '../api/user';
import type { UserLoginReq, UserRegisterReq, UserUpdateReq } from '../api/user';

interface UserState {
  token: string;
  username: string;
  realname: string;
  isLogin: boolean;
  isRecycleBinMode: boolean;
}

export const useUserStore = defineStore('user', {
  state: (): UserState => ({
    token: localStorage.getItem('token') || '',
    username: localStorage.getItem('username') || '',
    realname: localStorage.getItem('realname') || '',
    isLogin: !!localStorage.getItem('token'),
    isRecycleBinMode: localStorage.getItem('viewMode') === 'recycle'
  }),
  
  getters: {
    getToken: (state) => state.token,
    getUsername: (state) => state.username,
    getRealname: (state) => state.realname
  },
  
  actions: {
    // 设置登录信息
    setLoginInfo(token: string, username: string, realname: string) {
      if (!token || !username) {
        throw new Error('无效的登录信息');
      }

      this.token = token;
      this.username = username;
      this.realname = realname;
      this.isLogin = true;
      
      // 保存到本地存储
      localStorage.setItem('token', token);
      localStorage.setItem('username', username);
      localStorage.setItem('realname', realname);

      // 打印调试信息
      console.log('登录信息已保存:', {
        token: this.token,
        username: this.username,
        realname: this.realname,
        isLogin: this.isLogin
      });
    },
    
    // 设置视图模式（正常/回收站）
    setViewMode(mode: 'normal' | 'recycle') {
      this.isRecycleBinMode = mode === 'recycle';
      localStorage.setItem('viewMode', mode);
    },
    
    // 登录
    async login(loginForm: UserLoginReq) {
      try {
        const res = await userApi.login(loginForm);
        // 确保token存在且格式正确
        if (!res.data.token) {
          throw new Error('登录失败：未获取到token');
        }
        
        // 保存登录信息
        this.setLoginInfo(res.data.token, res.data.username, res.data.realname);
        
        // 打印调试信息
        console.log('登录响应:', res.data);
        
        ElMessage.success('登录成功');
        return Promise.resolve(res.data);
      } catch (error: any) {
        console.error('登录失败:', error);
        ElMessage.error(error.message || '登录失败');
        return Promise.reject(error);
      }
    },
    
    // 注册
    async register(registerForm: UserRegisterReq) {
      try {
        const res = await userApi.register(registerForm);
        ElMessage.success('注册成功');
        return Promise.resolve(res.data);
      } catch (error) {
        return Promise.reject(error);
      }
    },
    
    // 获取用户信息
    async fetchUserInfo() {
      try {
        const res = await userApi.getUserInfo(this.username);
        this.realname = res.data.realname;
        return Promise.resolve(res.data);
      } catch (error) {
        return Promise.reject(error);
      }
    },
    
    // 更新用户信息
    async updateUser(updateForm: UserUpdateReq) {
      try {
        const res = await userApi.updateUserInfo(updateForm);
        if (updateForm.realName) {
          this.realname = updateForm.realName;
          localStorage.setItem('realname', updateForm.realName);
        }
        ElMessage.success('更新成功');
        return Promise.resolve(res.data);
      } catch (error) {
        return Promise.reject(error);
      }
    },
    
    // 退出登录
    logout() {
      this.token = '';
      this.username = '';
      this.realname = '';
      this.isLogin = false;
      this.isRecycleBinMode = false;
      
      // 清除本地存储
      localStorage.removeItem('token');
      localStorage.removeItem('username');
      localStorage.removeItem('realname');
      localStorage.removeItem('viewMode');
      
      // 跳转到登录页
      router.push('/login');
    },
    
    // 检查登录状态
    checkLogin(): boolean {
      const token = this.token;
      const username = this.username;
      const isValid = this.isLogin && !!token && !!username;
      
      // 打印调试信息
      console.log('检查登录状态:', {
        token,
        username,
        isLogin: this.isLogin,
        isValid
      });
      
      // 如果状态无效，清除登录信息
      if (!isValid && this.isLogin) {
        console.warn('登录状态无效，执行登出操作');
        this.logout();
      }
      
      return isValid;
    }
  }
}); 