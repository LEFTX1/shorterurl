import axios from 'axios';
import type { SuccessResp } from './user';
import type { ShortLinkPageResp } from './link';

// 回收站分页请求接口
export interface RecycleBinPageReq {
  gid?: string; // 分组ID参数
  current?: number;
  size?: number;
  // 不需要显式传递enable_status和del_flag，由后端处理
}

// 回收站操作请求接口
export interface RecycleBinOperateReq {
  gid: string;
  fullShortUrl: string;
}

// 回收站API接口
export default {
  // 分页查询回收站 - 查询enable_status=1(未启用/回收站)且del_flag=0(未永久删除)的链接
  pageRecycleBin(params: RecycleBinPageReq) {
    return axios.get<ShortLinkPageResp>('/api/short-link/admin/v1/recycle-bin/page', { params });
  },
  
  // 保存到回收站 - 将enable_status设为1(未启用/回收站)
  saveToRecycleBin(data: RecycleBinOperateReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/recycle-bin/save', data);
  },
  
  // 从回收站恢复 - 将enable_status设为0(启用/正常状态)
  recoverFromRecycleBin(data: RecycleBinOperateReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/recycle-bin/recover', data);
  },
  
  // 从回收站永久删除 - 将del_flag设为1(永久删除状态)
  removeFromRecycleBin(data: RecycleBinOperateReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/recycle-bin/remove', data);
  }
}; 