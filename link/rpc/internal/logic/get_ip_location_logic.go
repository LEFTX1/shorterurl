package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 高德地图API常量
const (
	// 高德地图API
	AmapIPLocationAPI = "https://restapi.amap.com/v3/ip"
	AmapAPIKey        = "9891e494403818e3fc79fb61fcf06b84"

	// IP位置缓存Key前缀
	IPLocationCacheKey = "ip:location:%s"
	// IP位置缓存时间 (秒)
	IPLocationCacheExpire = 30 * 24 * 60 * 60 // 30天
)

// IP位置信息响应
type IPLocationResponse struct {
	Status    string `json:"status"`    // 返回结果状态值：1成功，0失败
	Info      string `json:"info"`      // 返回状态说明
	InfoCode  string `json:"infocode"`  // 状态码
	Province  string `json:"province"`  // 省份名称
	City      string `json:"city"`      // 城市名称
	AdCode    string `json:"adcode"`    // 城市的adcode编码
	Rectangle string `json:"rectangle"` // 所在城市矩形区域范围
}

type GetIpLocationLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetIpLocationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetIpLocationLogic {
	return &GetIpLocationLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetIpLocation 获取IP地理位置信息接口
func (l *GetIpLocationLogic) GetIpLocation(in *pb.GetIPLocationRequest) (*pb.GetIPLocationResponse, error) {
	// 日志记录
	l.Logger.Infof("获取IP地址 [%s] 的地理位置信息", in.Ip)

	// 参数校验
	if in.Ip == "" {
		return nil, status.Error(codes.InvalidArgument, "IP不能为空")
	}

	// 检查缓存
	cacheKey := fmt.Sprintf(IPLocationCacheKey, in.Ip)
	cachedResult, err := l.svcCtx.BizRedis.Get(cacheKey)
	if err == nil && cachedResult != "" {
		// 从缓存中获取成功
		var response pb.GetIPLocationResponse
		if err := json.Unmarshal([]byte(cachedResult), &response); err == nil {
			return &response, nil
		}
	}

	// 调用高德地图 API 获取 IP 地理位置
	response, err := l.fetchIPLocationFromAmap(in.Ip)
	if err != nil {
		l.Logger.Errorf("调用高德地图 API 获取 IP 地理位置失败: %v", err)
		return nil, status.Error(codes.Internal, "获取 IP 地理位置失败")
	}

	// 缓存结果（1 天有效期）
	if response.Status == "1" {
		jsonBytes, _ := json.Marshal(response)
		l.svcCtx.BizRedis.Setex(cacheKey, string(jsonBytes), 24*60*60)
	}

	return response, nil
}

// fetchIPLocationFromAPI 从高德地图API获取IP位置信息
func (l *GetIpLocationLogic) fetchIPLocationFromAPI(ip string) (*IPLocationResponse, error) {
	// 创建带超时的context
	ctx, cancel := context.WithTimeout(l.ctx, 5*time.Second)
	defer cancel()

	// 构建请求URL
	reqUrl := fmt.Sprintf("%s?key=%s&ip=%s", AmapIPLocationAPI, AmapAPIKey, ip)

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求高德地图API失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求错误, 状态码: %d", resp.StatusCode)
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取API响应失败: %v", err)
	}

	// 解析JSON响应
	var result IPLocationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析API响应失败: %v", err)
	}

	return &result, nil
}

// GetFormattedLocation 获取格式化的地理位置信息
func (l *GetIpLocationLogic) GetFormattedLocation(ip string) (string, error) {
	location, err := l.GetIpLocation(&pb.GetIPLocationRequest{Ip: ip})
	if err != nil {
		return "", err
	}

	if location.Status != "1" {
		return "", fmt.Errorf("获取地理位置失败: %s", location.Info)
	}

	// 如果是境外地址，可能只有省份没有城市
	if location.City == "" {
		return location.Province, nil
	}

	// 如果省份和城市相同（如直辖市），只返回城市名
	if location.Province == location.City {
		return location.City, nil
	}

	// 返回 "省份,城市" 格式
	return fmt.Sprintf("%s,%s", location.Province, location.City), nil
}

// 调用高德地图 API 获取 IP 地理位置
func (l *GetIpLocationLogic) fetchIPLocationFromAmap(ip string) (*pb.GetIPLocationResponse, error) {
	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 构建请求 URL
	// 使用配置的高德地图 API Key（这里使用了用户提供的 Key）
	url := fmt.Sprintf("%s?key=%s&ip=%s", AmapIPLocationAPI, "9891e494403818e3fc79fb61fcf06b84", ip)

	// 发送请求
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求高德地图 API 失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %v", err)
	}

	// 解析 JSON 响应
	var amapResponse struct {
		Status    string `json:"status"`
		Info      string `json:"info"`
		Infocode  string `json:"infocode"`
		Province  string `json:"province"`
		City      string `json:"city"`
		Adcode    string `json:"adcode"`
		Rectangle string `json:"rectangle"`
	}

	if err := json.Unmarshal(body, &amapResponse); err != nil {
		return nil, fmt.Errorf("解析响应内容失败: %v", err)
	}

	// 转换为 gRPC 响应格式
	return &pb.GetIPLocationResponse{
		Status:    amapResponse.Status,
		Info:      amapResponse.Info,
		Infocode:  amapResponse.Infocode,
		Province:  amapResponse.Province,
		City:      amapResponse.City,
		Adcode:    amapResponse.Adcode,
		Rectangle: amapResponse.Rectangle,
	}, nil
}
