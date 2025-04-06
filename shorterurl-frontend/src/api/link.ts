import axios from 'axios';
import type { SuccessResp } from './user';

// 短链接记录
export interface ShortLinkRecord {
  id: number;
  domain: string;
  shortUri: string;
  fullShortUrl: string;
  originUrl: string;
  gid: string;
  createTime: string;
  validDate: string;
  describe: string;
  validDateType?: number;
  enableStatus?: number;
  totalPv: number;
  totalUv: number;
  totalUip: number;
  todayPv?: number;
  todayUv?: number;
  todayUip?: number;
}

// 短链分页响应
export interface ShortLinkPageResp {
  records: ShortLinkRecord[];
  total: number;
  size: number;
  current: number;
}

// 创建短链请求
export interface CreateShortLinkReq {
  originUrl: string;
  gid: string;
  validDateType?: number;
  validDate?: string;
  describe?: string;
}

// 批量创建短链请求
export interface BatchCreateShortLinkReq {
  originUrls: string;
  gid: string;
  validDateType?: number;
  validDate?: string;
  describe?: string;
}

// 创建短链接响应接口
export interface CreateShortLinkResp {
  fullShortUrl: string;
  originUrl: string;
  gid: string;
}

// 批量创建短链接响应接口
export interface BatchCreateShortLinkResp {
  total: number;
  baseLinkInfos: ShortLinkBaseInfo[];
}

// 更新短链接请求接口
export interface UpdateShortLinkReq {
  fullShortUrl: string;
  originUrl: string;
  gid: string;
  validDateType: number;
  validDate?: string;
  describe?: string;
}

// 分页查询短链接请求接口
export interface PageShortLinkReq {
  gid: string;
  current: number;
  size: number;
}

// 短链接基本信息接口
export interface ShortLinkBaseInfo {
  fullShortUrl: string;
  originUrl: string;
  describe: string;
}

// 回收站操作请求接口
export interface RecycleBinOperateReq {
  gid: string;
  fullShortUrl: string;
}

// 短链接API
export default {
  // 分页查询短链
  pageShortLink(params: any) {
    return axios.get<ShortLinkPageResp>('/api/short-link/admin/v1/link', { params });
  },
  
  // 创建短链
  createShortLink(data: CreateShortLinkReq) {
    return axios.post<CreateShortLinkResp>('/api/short-link/admin/v1/link', data);
  },
  
  // 更新短链
  updateShortLink(data: UpdateShortLinkReq) {
    return axios.put<SuccessResp>('/api/short-link/admin/v1/link', data);
  },
  
  // 批量创建短链
  batchCreateShortLink(data: BatchCreateShortLinkReq) {
    return axios.post<BatchCreateShortLinkResp>('/api/short-link/admin/v1/link/batch', data);
  },
  
  // 保存到回收站
  saveToRecycleBin(data: RecycleBinOperateReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/recycle-bin/save', data);
  },
  
  // 从回收站恢复
  recoverFromRecycleBin(data: RecycleBinOperateReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/recycle-bin/recover', data);
  },
  
  // 从回收站删除
  removeFromRecycleBin(data: RecycleBinOperateReq) {
    return axios.post<SuccessResp>('/api/short-link/admin/v1/recycle-bin/remove', data);
  },
  
  // 获取网站标题
  getUrlTitle(url: string) {
    return axios.get<{data: string}>('/api/short-link/admin/v1/title', {
      params: { url }
    });
  }
}; 