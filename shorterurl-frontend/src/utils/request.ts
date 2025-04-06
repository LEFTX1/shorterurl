import axios from 'axios';
import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { ElMessage } from 'element-plus';

// axios实例
const service: AxiosInstance = axios.create({
  baseURL: '',
  timeout: 15000
});

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    // 从本地存储获取token
    const token = localStorage.getItem('token');
    const username = localStorage.getItem('username');
    
    // 如果有token，添加到请求头
    if (token && username) {
      config.headers['Authorization'] = `Bearer ${token}`;
      config.headers['Username'] = username;
    }
    
    return config;
  },
  (error) => {
    console.error('请求错误:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse) => {
    const data = response.data;
    
    // 如果响应码不是200，则认为有错误
    if (response.status !== 200) {
      ElMessage.error('网络错误，请稍后重试');
      return Promise.reject(new Error('网络错误'));
    }
    
    // 处理业务错误码
    if (data.code && data.code !== '0000' && data.code !== 'SUCCESS') {
      ElMessage.error(data.message || '操作失败');
      return Promise.reject(new Error(data.message || '未知错误'));
    }
    
    return data;
  },
  (error) => {
    console.error('响应错误:', error);
    
    // 处理401未授权错误
    if (error.response && error.response.status === 401) {
      ElMessage.error('登录已过期，请重新登录');
      // 清除本地存储
      localStorage.removeItem('token');
      localStorage.removeItem('username');
      
      // 跳转到登录页面
      window.location.href = '/login';
      return Promise.reject(error);
    }
    
    ElMessage.error(error.message || '网络错误，请稍后重试');
    return Promise.reject(error);
  }
);

// 封装请求方法
export const request = {
  get<T>(url: string, params?: any): Promise<T> {
    return service.get(url, { params });
  },
  
  post<T>(url: string, data?: any): Promise<T> {
    return service.post(url, data);
  },
  
  put<T>(url: string, data?: any): Promise<T> {
    return service.put(url, data);
  },
  
  delete<T>(url: string, params?: any): Promise<T> {
    return service.delete(url, { params });
  }
};

export default service; 