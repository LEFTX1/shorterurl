import { defineStore } from 'pinia';
import { user as userApi } from '@/api';
import type { UserInfoResp, UserLoginReq, UserRegisterReq, UserUpdateReq } from '@/api/user';
import { ElMessage } from 'element-plus';
import router from '@/router';

export const useUserStore = defineStore('user', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    username: localStorage.getItem('username') || '',
    realname: localStorage.getItem('realname') || '',
    userInfo: null as UserInfoResp | null,
    isLogin: !!localStorage.getItem('token'),
    viewMode: localStorage.getItem('viewMode') || 'normal' // 'normal' 或 'recycle'
  }),
  
  getters: {
    getUserInfo: (state) => state.userInfo,
    getUsername: (state) => state.username,
    getRealname: (state) => state.realname,
    getToken: (state) => state.token,
    getLoginStatus: (state) => state.isLogin,
    getViewMode: (state) => state.viewMode,
    isRecycleBinMode: (state) => state.viewMode === 'recycle'
  },
  
  actions: {
    // 设置视图模式
    setViewMode(mode: 'normal' | 'recycle') {
      this.viewMode = mode;
      localStorage.setItem('viewMode', mode);
    },
    
    // 登录
    async login(loginForm: UserLoginReq) {
      try {
        const res = await userApi.login(loginForm);
        this.token = res.data.token;
        this.username = res.data.username;
        this.realname = res.data.realname;
        this.isLogin = true;
        
        // 保存到本地存储
        localStorage.setItem('token', res.data.token);
        localStorage.setItem('username', res.data.username);
        localStorage.setItem('realname', res.data.realname);
        
        // 获取用户信息
        await this.fetchUserInfo();
        
        return Promise.resolve(res.data);
      } catch (error) {
        return Promise.reject(error);
      }
    },
    
    // 注册
    async register(registerForm: UserRegisterReq) {
      try {
        const res = await userApi.register(registerForm);
        ElMessage.success('注册成功，请登录');
        return Promise.resolve(res.data);
      } catch (error) {
        return Promise.reject(error);
      }
    },
    
    // 获取用户信息
    async fetchUserInfo() {
      if (!this.username) return;
      
      try {
        const res = await userApi.getUserInfo(this.username);
        this.userInfo = res.data;
        return Promise.resolve(res.data);
      } catch (error) {
        return Promise.reject(error);
      }
    },
    
    // 更新用户信息
    async updateUserInfo(updateForm: UserUpdateReq) {
      try {
        const res = await userApi.updateUserInfo(updateForm);
        ElMessage.success('更新成功');
        
        // 重新获取用户信息
        await this.fetchUserInfo();
        
        return Promise.resolve(res.data);
      } catch (error) {
        return Promise.reject(error);
      }
    },
    
    // 检查登录状态
    async checkLogin() {
      if (!this.token || !this.username) {
        this.logout();
        return false;
      }
      
      try {
        const res = await userApi.checkLogin(this.username, this.token);
        return res.data.success;
      } catch (error) {
        this.logout();
        return false;
      }
    },
    
    // 退出登录
    async logout() {
      if (this.token && this.username) {
        try {
          await userApi.logout(this.username, this.token);
        } catch (error) {
          console.error('退出登录错误:', error);
        }
      }
      
      // 清除状态
      this.token = '';
      this.username = '';
      this.realname = '';
      this.userInfo = null;
      this.isLogin = false;
      this.viewMode = 'normal';
      
      // 清除本地存储
      localStorage.removeItem('token');
      localStorage.removeItem('username');
      localStorage.removeItem('realname');
      localStorage.removeItem('viewMode');
      
      // 跳转到登录页
      router.push('/login');
    }
  }
}); 