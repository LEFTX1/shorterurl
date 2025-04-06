import axios from 'axios';
import type { SuccessResp } from './user';

// 创建分组请求接口
export interface ShortLinkGroupSaveReq {
  name: string;
}

// 更新分组请求接口
export interface ShortLinkGroupUpdateReq {
  gid: string;
  name: string;
}

// 分组排序请求接口
export interface ShortLinkGroupSortReq {
  groups: SortGroup[];
}

// 排序分组项接口
export interface SortGroup {
  gid: string;
  sortOrder: number;
}

// 分组响应接口
export interface ShortLinkGroupResp {
  gid: string;
  name: string;
  sortOrder: number;
  shortLinkCount: number;
}

// 分组API
export default {
  // 创建分组
  createGroup(data: ShortLinkGroupSaveReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/group', data);
  },
  
  // 获取分组列表
  listGroups() {
    // 添加日志输出存储的认证信息
    const token = localStorage.getItem('token');
    const username = localStorage.getItem('username');
    
    console.log('获取分组前的认证信息:', { token, username });
    
    // 手动添加认证头
    const headers: Record<string, string> = {};
    if (token && username) {
      headers['token'] = token;
      headers['username'] = username;
    }
    
    console.log('即将发送的请求头:', headers);
    
    return axios.get<ShortLinkGroupResp[]>('/api/short-link/admin/v1/group', { headers });
  },
  
  // 更新分组
  updateGroup(data: ShortLinkGroupUpdateReq) {
    return axios.put<SuccessResp>('/api/short-link/admin/v1/group', data);
  },
  
  // 删除分组
  deleteGroup(gid: string) {
    return axios.delete<SuccessResp>('/api/short-link/admin/v1/group', {
      params: { gid }
    });
  },
  
  // 分组排序
  sortGroups(data: ShortLinkGroupSortReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/group/sort', data);
  }
}; 