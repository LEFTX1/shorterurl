import axios from 'axios';
import type { AxiosRequestConfig, AxiosResponse, AxiosError } from 'axios';
import { ElMessage } from 'element-plus';
import { useUserStore } from '@/store/user';
import router from '@/router';

// 创建带认证头的 axios 实例
const createAuthInstance = () => {
  const instance = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    timeout: 10000,
  });

  // 请求拦截器
  instance.interceptors.request.use(
    (config) => {
      const userStore = useUserStore();
      const { token, username } = userStore;

      if (token && username) {
        config.headers = config.headers || {};
        config.headers['token'] = token;
        config.headers['username'] = username;
      }

      // 调试日志
      console.log('请求配置:', {
        url: config.url,
        method: config.method,
        params: config.params,
        data: config.data,
        headers: {
          token: config.headers?.['token'],
          username: config.headers?.['username'],
        },
      });

      return config;
    },
    (error) => {
      console.error('请求错误:', error);
      return Promise.reject(error);
    }
  );

  // 响应拦截器
  instance.interceptors.response.use(
    (response: AxiosResponse) => {
      const { data } = response;
      if (data.code === 0) {
        return data.data;
      }
      ElMessage.error(data.msg || '请求失败');
      return Promise.reject(new Error(data.msg || '请求失败'));
    },
    (error: AxiosError) => {
      console.error('响应错误:', {
        config: error.config,
        response: error.response,
        message: error.message,
      });

      if (error.response?.status === 401) {
        const userStore = useUserStore();
        userStore.logout();
        router.push('/login');
        ElMessage.error('登录已过期，请重新登录');
      } else {
        ElMessage.error(error.message || '请求失败');
      }
      return Promise.reject(error);
    }
  );

  return instance;
};

// 导出带认证头的实例
export const authRequest = createAuthInstance();

// 导出普通实例
export const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
});

// 默认导出带认证头的实例
export default authRequest; 