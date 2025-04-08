import axios from 'axios'
import type { SuccessResp } from './user';
import type { AxiosResponse } from 'axios'

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

// 统计请求参数接口
export interface StatsQueryParams {
  fullShortUrl: string;
  gid: string;
  startDate?: string;
  endDate?: string;
  current?: number;
  size?: number;
}

// 统计数据接口
export interface StatsData {
  date: string;
  cnt: number;
  [key: string]: any;
}

// 设备统计数据
export interface DeviceStats extends StatsData {
  device: string;
}

// 浏览器统计数据
export interface BrowserStats extends StatsData {
  browser: string;
}

// 操作系统统计数据
export interface OsStats extends StatsData {
  os: string;
}

// 网络类型统计数据
export interface NetworkStats extends StatsData {
  network: string;
}

// 地区统计数据
export interface LocaleStats extends StatsData {
  country: string;
  province: string;
  city: string;
  adcode: string;
}

// 访问日志数据
export interface AccessLog {
  ip: string;
  browser: string;
  os: string;
  device: string;
  network: string;
  locale: string;
  create_time: string;
}

// 获取认证头
const getAuthHeaders = () => {
  const token = localStorage.getItem('token')
  const username = localStorage.getItem('username')
  
  console.log('获取认证信息:', { token, username })
  
  const headers: Record<string, string> = {}
  if (token && username) {
    headers['token'] = token
    headers['username'] = username
  }
  
  console.log('生成的请求头:', headers)
  return headers
}

// API响应类型
interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

// 统计API
export default {
  // 获取单个短链接监控数据
  getShortLinkStats(params: ShortLinkStatsReq) {
    return axios.get('/api/short-link/admin/v1/stats', {
      params,
      headers: getAuthHeaders()
    });
  },
  
  // 获取分组短链接监控数据
  getShortLinkGroupStats(params: ShortLinkGroupStatsReq) {
    return axios.get('/api/short-link/admin/v1/stats/group', {
      params,
      headers: getAuthHeaders()
    }).then(response => response.data);
  },
  
  // 短链接访问记录查询
  getShortLinkAccessRecord(params: ShortLinkAccessRecordReq) {
    return axios.get('/api/short-link/admin/v1/stats/access-record', {
      params,
      headers: getAuthHeaders()
    });
  },
  
  // 分组短链接访问记录查询
  getShortLinkGroupAccessRecord(params: ShortLinkGroupAccessRecordReq) {
    return axios.get('/api/short-link/admin/v1/stats/access-record/group', {
      params,
      headers: getAuthHeaders()
    });
  },

  // 获取设备统计
  getDeviceStats(params: StatsQueryParams) {
    return axios.get('/api/short-link/admin/v1/stats', {
      params: {
        ...params,
        startDate: new Date().toISOString().split('T')[0],
        endDate: new Date().toISOString().split('T')[0]
      },
      headers: getAuthHeaders()
    });
  },

  // 获取浏览器统计
  getBrowserStats(params: StatsQueryParams) {
    return axios.get('/api/short-link/admin/v1/stats', {
      params: {
        ...params,
        startDate: new Date().toISOString().split('T')[0],
        endDate: new Date().toISOString().split('T')[0]
      },
      headers: getAuthHeaders()
    });
  },

  // 获取操作系统统计
  getOsStats(params: StatsQueryParams) {
    return axios.get('/api/short-link/admin/v1/stats', {
      params: {
        ...params,
        startDate: new Date().toISOString().split('T')[0],
        endDate: new Date().toISOString().split('T')[0]
      },
      headers: getAuthHeaders()
    });
  },

  // 获取网络类型统计
  getNetworkStats(params: StatsQueryParams) {
    return axios.get('/api/short-link/admin/v1/stats', {
      params: {
        ...params,
        startDate: new Date().toISOString().split('T')[0],
        endDate: new Date().toISOString().split('T')[0]
      },
      headers: getAuthHeaders()
    });
  },

  // 获取地区统计
  getLocaleStats(params: StatsQueryParams) {
    return axios.get('/api/short-link/admin/v1/stats', {
      params: {
        ...params,
        startDate: new Date().toISOString().split('T')[0],
        endDate: new Date().toISOString().split('T')[0]
      },
      headers: getAuthHeaders()
    });
  },

  // 获取访问日志
  getAccessLogs(params: StatsQueryParams) {
    return axios.get('/api/short-link/admin/v1/stats/access-record', {
      params: {
        ...params,
        startDate: new Date().toISOString().split('T')[0],
        endDate: new Date().toISOString().split('T')[0],
        current: 1,
        size: 100
      },
      headers: getAuthHeaders()
    });
  }
}; 