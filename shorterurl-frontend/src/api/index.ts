import axios from 'axios';
import { ElMessage } from 'element-plus';
import router from '../router';
import user from './user';
import link from './link';
import group from './group';
import stats from './stats';
import recycle from './recycle';

// 创建axios实例
const request = axios.create({
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
});

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 从本地存储获取token
    const token = localStorage.getItem('token');
    if (token) {
      // 将token添加到请求头
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  response => {
    return response;
  },
  error => {
    if (error.response) {
      // 处理401未授权错误
      if (error.response.status === 401) {
        // 清除token和用户信息
        localStorage.removeItem('token');
        localStorage.removeItem('username');
        
        // 显示错误消息
        ElMessage.error('登录已过期，请重新登录');
        
        // 跳转到登录页
        router.push('/login');
      } else {
        // 处理其他错误
        ElMessage.error(error.response.data.message || '请求失败');
      }
    } else {
      // 处理网络错误或其他错误
      ElMessage.error('网络错误，请稍后重试');
    }
    return Promise.reject(error);
  }
);

export { user, link, group, stats, recycle };

export default {
  request,
  user,
  group,
  link,
  stats,
  recycle
}; 