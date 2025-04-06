import axios from 'axios';
import type { SuccessResp } from './user';
import type { ShortLinkPageResp } from './link';

// 回收站分页请求接口
export interface RecycleBinPageReq {
  gidList?: string[];
  current?: number;
  size?: number;
}

// 回收站操作请求接口
export interface RecycleBinOperateReq {
  gid: string;
  fullShortUrl: string;
}

// 回收站API接口
export default {
  // 分页查询回收站
  pageRecycleBin(params: RecycleBinPageReq) {
    return axios.get<ShortLinkPageResp>('/api/short-link/admin/v1/recycle-bin/page', { params });
  },
  
  // 保存到回收站
  saveToRecycleBin(data: RecycleBinOperateReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/recycle-bin/save', data);
  },
  
  // 从回收站恢复
  recoverFromRecycleBin(data: RecycleBinOperateReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/recycle-bin/recover', data);
  },
  
  // 从回收站删除(永久删除)
  removeFromRecycleBin(data: RecycleBinOperateReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/recycle-bin/remove', data);
  }
}; 