import axios from 'axios';

// 用户注册请求接口
export interface UserRegisterReq {
  username: string;
  password: string;
  realname: string;
  phone: string;
  mail: string;
}

// 用户登录请求接口
export interface UserLoginReq {
  username: string;
  password: string;
}

// 用户信息响应接口
export interface UserInfoResp {
  id: number;
  username: string;
  realname: string;
  phone: string;
  mail: string;
  createTime: string;
  updateTime: string;
}

// 用户登录响应接口
export interface UserLoginResp {
  token: string;
  username: string;
  realname: string;
  createTime: string;
}

// 用户更新请求接口
export interface UserUpdateReq {
  username: string;
  password?: string;
  realName?: string;
  phone?: string;
  mail?: string;
}

// 成功响应接口
export interface SuccessResp {
  code: string;
  success: boolean;
  message?: string;
}

// 用户API
export default {
  // 用户注册
  register(data: UserRegisterReq) {
    return axios.post<any>('/api/short-link/admin/v1/user', data);
  },
  
  // 用户登录
  login(data: UserLoginReq) {
    return axios.post<UserLoginResp>('/api/short-link/admin/v1/user/login', data);
  },
  
  // 获取用户信息
  getUserInfo(username: string) {
    return axios.get<UserInfoResp>(`/api/short-link/admin/v1/user/${username}`);
  },
  
  // 获取无脱敏用户信息
  getActualUserInfo(username: string) {
    return axios.get<UserInfoResp>(`/api/short-link/admin/v1/actual/user/${username}`);
  },
  
  // 检查用户名是否存在
  checkUsername(username: string) {
    return axios.get<any>('/api/short-link/admin/v1/user/has-username', {
      params: { username }
    });
  },
  
  // 更新用户信息
  updateUserInfo(data: UserUpdateReq) {
    return axios.put<any>('/api/short-link/admin/v1/user', data);
  },
  
  // 检查用户是否登录
  checkLogin(username: string, token: string) {
    return axios.get<SuccessResp>('/api/short-link/admin/v1/user/check-login', {
      params: { username, token }
    });
  },
  
  // 用户退出登录
  logout(username: string, token: string) {
    return axios.delete<SuccessResp>('/api/short-link/admin/v1/user/logout', {
      params: { username, token }
    });
  }
}; 