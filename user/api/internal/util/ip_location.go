package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// IPLocationResponse 高德地图IP定位API响应结构
type IPLocationResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	InfoCode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	Adcode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

// GetIPLocation 获取IP地址对应的地理位置信息
func GetIPLocation(ip string) (string, error) {
	// 如果IP为空，返回空字符串
	if ip == "" {
		return "", nil
	}

	// 如果IP是本地IP，返回"本地"
	if ip == "127.0.0.1" || ip == "localhost" || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.16.") {
		return "本地", nil
	}

	// 构建API请求URL
	apiKey := "9891e494403818e3fc79fb61fcf06b84"
	url := fmt.Sprintf("https://restapi.amap.com/v3/ip?ip=%s&key=%s", ip, apiKey)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 发送请求
	resp, err := client.Get(url)
	if err != nil {
		logx.Errorf("请求高德地图API失败: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("读取高德地图API响应失败: %v", err)
		return "", err
	}

	// 解析JSON响应
	var result IPLocationResponse
	if err := json.Unmarshal(body, &result); err != nil {
		logx.Errorf("解析高德地图API响应失败: %v", err)
		return "", err
	}

	// 检查响应状态
	if result.Status != "1" {
		logx.Errorf("高德地图API返回错误: %s", result.Info)
		return "", fmt.Errorf("API返回错误: %s", result.Info)
	}

	// 构建地理位置字符串
	location := ""
	if result.Province != "" {
		location = result.Province
	}
	if result.City != "" && result.City != result.Province {
		if location != "" {
			location += "-"
		}
		location += result.City
	}

	// 如果地理位置为空，返回"未知地区"
	if location == "" {
		location = "未知地区"
	}

	return location, nil
}
