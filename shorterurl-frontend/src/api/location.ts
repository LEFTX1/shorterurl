import axios from 'axios';

// 高德地图API密钥
const AMAP_KEY = '9891e494403818e3fc79fb61fcf06b84';

// IP定位响应接口
export interface IPLocationResponse {
  status: string;       // 返回结果状态值 值为0或1,0表示失败；1表示成功
  info: string;         // 返回状态说明，status为0时返回错误原因，否则返回"OK"
  infocode: string;     // 状态码，10000代表正确
  province: string;     // 省份名称，直辖市显示直辖市名称；国外显示空
  city: string;         // 城市名称，直辖市显示直辖市名称；国外显示空
  adcode: string;       // 城市的adcode编码
  rectangle: string;    // 所在城市矩形区域范围
}

// 后端IP定位响应接口
export interface ServerIPLocationResponse {
  status: string;       // 返回结果状态值 值为0或1,0表示失败；1表示成功
  info: string;         // 返回状态说明，status为0时返回错误原因，否则返回"OK"
  infocode: string;     // 状态码，10000代表正确，与高德地图API保持一致
  province: string;     // 省份名称
  city: string;         // 城市名称
  adcode: string;       // 城市的adcode编码
  rectangle: string;    // 所在城市矩形区域范围
}

// 位置服务API
export default {
  /**
   * IP定位，获取IP所在地理位置（直接调用高德API）
   * @param ip 可选，IP地址。若不传则使用发起请求的客户端IP
   * @returns 返回IP定位结果
   */
  getIpLocation(ip?: string) {
    const params: Record<string, string> = {
      key: AMAP_KEY
    };
    
    if (ip) {
      params.ip = ip;
    }
    
    return axios.get<IPLocationResponse>('https://restapi.amap.com/v3/ip', { 
      params,
      // 添加超时设置
      timeout: 10000
    });
  },

  /**
   * 获取IP地理位置信息（调用后端接口，后端会缓存结果）
   * @param ip IP地址
   * @returns 返回IP地理位置信息
   */
  getServerIpLocation(ip: string) {
    return axios.get<ServerIPLocationResponse>('/api/short-link/admin/v1/ip-location', {
      params: { ip },
      timeout: 10000
    });
  }
}; 