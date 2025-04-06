import axios from 'axios';
import type { SuccessResp } from './user';

// 统计请求参数接口
export interface ShortLinkStatsReq {
  fullShortUrl: string;
  gid: string;
  enableStatus?: number;
  startDate: string;
  endDate: string;
}

// 分组统计请求参数接口
export interface ShortLinkGroupStatsReq {
  gid: string;
  startDate: string;
  endDate: string;
}

// 访问记录查询请求接口
export interface ShortLinkAccessRecordReq {
  fullShortUrl: string;
  gid: string;
  startDate: string;
  endDate: string;
  enableStatus?: number;
  current: number;
  size: number;
}

// 分组访问记录查询请求接口
export interface ShortLinkGroupAccessRecordReq {
  gid: string;
  startDate: string;
  endDate: string;
  current: number;
  size: number;
}

// PV/UV/UIP统计接口
export interface PvUvUipStats {
  date: string;
  pv: number;
  uv: number;
  uip: number;
}

// 访问记录分页响应接口
export interface AccessRecordPageResp {
  records: AccessRecord[];
  total: number;
  size: number;
  current: number;
}

// 访问记录接口
export interface AccessRecord {
  ip: string;
  browser: string;
  os: string;
  network: string;
  device: string;
  locale: string;
  accessTime: string;
}

// 地域统计接口
export interface LocaleCnStat {
  locale: string;
  cnt: number;
  ratio: number;
}

// 高频访问IP统计接口
export interface TopIpStat {
  ip: string;
  cnt: number;
  ratio: number;
}

// 浏览器统计接口
export interface BrowserStat {
  browser: string;
  cnt: number;
  ratio: number;
}

// 操作系统统计接口
export interface OsStat {
  os: string;
  cnt: number;
  ratio: number;
}

// 访客类型统计接口
export interface UvTypeStat {
  uvType: string;
  cnt: number;
  ratio: number;
}

// 设备统计接口
export interface DeviceStat {
  device: string;
  cnt: number;
  ratio: number;
}

// 网络统计接口
export interface NetworkStat {
  network: string;
  cnt: number;
  ratio: number;
}

// 统计响应接口
export interface ShortLinkStatsRespDTO {
  pvUvUipStatsList: PvUvUipStats[];
  overallPvUvUipStats: PvUvUipStats;
  localeCnStats: LocaleCnStat[];
  hourStats: number[];
  topIpStats: TopIpStat[];
  weekdayStats: number[];
  browserStats: BrowserStat[];
  osStats: OsStat[];
  uvTypeStats: UvTypeStat[];
  deviceStats: DeviceStat[];
  networkStats: NetworkStat[];
}

// 统计API
export default {
  // 获取单个短链接监控数据
  getShortLinkStats(params: ShortLinkStatsReq) {
    return axios.get<ShortLinkStatsRespDTO>('/api/short-link/admin/v1/stats', { params });
  },
  
  // 获取分组短链接监控数据
  getShortLinkGroupStats(params: ShortLinkGroupStatsReq) {
    return axios.get<ShortLinkStatsRespDTO>('/api/short-link/admin/v1/stats/group', { params });
  },
  
  // 短链接访问记录查询
  getShortLinkAccessRecord(params: ShortLinkAccessRecordReq) {
    return axios.get<AccessRecordPageResp>('/api/short-link/admin/v1/stats/access-record', { params });
  },
  
  // 分组短链接访问记录查询
  getShortLinkGroupAccessRecord(params: ShortLinkGroupAccessRecordReq) {
    return axios.get<AccessRecordPageResp>('/api/short-link/admin/v1/stats/access-record/group', { params });
  }
}; 