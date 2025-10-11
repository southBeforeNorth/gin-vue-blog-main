package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"time"

	global "gin-blog/internal/global"
)

var Geocoding = new(geocodingUtil)

type geocodingUtil struct{}

// 高德地图逆地理编码API响应结构
type AmapGeocodeResponse struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Infocode  string `json:"infocode"`
	Regeocode struct {
		FormattedAddress string `json:"formatted_address"`
		AddressComponent struct {
			Country  string `json:"country"`
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"district"`
		} `json:"addressComponent"`
	} `json:"regeocode"`
}

// 百度地图逆地理编码API响应结构
type BaiduGeocodeResponse struct {
	Status int `json:"status"`
	Result struct {
		FormattedAddress string `json:"formatted_address"`
		AddressComponent struct {
			Country  string `json:"country"`
			Province string `json:"province"`
			City     string `json:"city"`
			District string `json:"district"`
		} `json:"addressComponent"`
	} `json:"result"`
}

// 使用高德地图API进行逆地理编码
func (g *geocodingUtil) GetAddressByAmap(latitude, longitude float64) (string, error) {
	// 从配置文件中获取高德地图的API Key
	config := global.GetConfig()
	apiKey := config.Geocoding.AmapApiKey
	if apiKey == "" {
		return "", fmt.Errorf("高德地图API Key未配置")
	}

	url := fmt.Sprintf("https://restapi.amap.com/v3/geocode/regeo?key=%s&location=%f,%f&poitype=&radius=1000&extensions=base&batch=false&roadlevel=0",
		apiKey, longitude, latitude)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求高德地图API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}
	slog.Debug("高德地图API响应: " + string(body))
	var result AmapGeocodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Status != "1" {
		return "", fmt.Errorf("高德地图API返回错误: %s", result.Info)
	}

	return result.Regeocode.FormattedAddress, nil
}

// 使用百度地图API进行逆地理编码
func (g *geocodingUtil) GetAddressByBaidu(latitude, longitude float64) (string, error) {
	// 从配置文件中获取百度地图的API Key
	config := global.GetConfig()
	apiKey := config.Geocoding.BaiduApiKey
	if apiKey == "" {
		return "", fmt.Errorf("百度地图API Key未配置")
	}

	// 百度地图需要将WGS84坐标转换为BD09坐标
	bdLat, bdLng := g.WGS84ToBD09(latitude, longitude)

	url := fmt.Sprintf("https://api.map.baidu.com/reverse_geocoding/v3/?ak=%s&output=json&coordtype=bd09ll&location=%f,%f",
		apiKey, bdLat, bdLng)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("请求百度地图API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	var result BaiduGeocodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}

	if result.Status != 0 {
		return "", fmt.Errorf("百度地图API返回错误: status=%d", result.Status)
	}

	return result.Result.FormattedAddress, nil
}

// 坐标转换：WGS84转BD09
func (g *geocodingUtil) WGS84ToBD09(lat, lng float64) (float64, float64) {
	// 先转换为GCJ02
	gcjLat, gcjLng := g.WGS84ToGCJ02(lat, lng)
	// 再转换为BD09
	return g.GCJ02ToBD09(gcjLat, gcjLng)
}

// 坐标转换：WGS84转GCJ02
func (g *geocodingUtil) WGS84ToGCJ02(lat, lng float64) (float64, float64) {
	const a = 6378245.0
	const ee = 0.00669342162296594323

	if g.outOfChina(lat, lng) {
		return lat, lng
	}

	dlat := g.transformLat(lng-105.0, lat-35.0)
	dlng := g.transformLng(lng-105.0, lat-35.0)

	radlat := lat / 180.0 * 3.14159265358979323846264338327950288
	magic := math.Sin(radlat)
	magic = 1 - ee*magic*magic
	sqrtmagic := math.Sqrt(magic)

	dlat = (dlat * 180.0) / ((a * (1 - ee)) / (magic * sqrtmagic) * 3.14159265358979323846264338327950288)
	dlng = (dlng * 180.0) / (a / sqrtmagic * math.Cos(radlat) * 3.14159265358979323846264338327950288)

	return lat + dlat, lng + dlng
}

// 坐标转换：GCJ02转BD09
func (g *geocodingUtil) GCJ02ToBD09(lat, lng float64) (float64, float64) {
	z := math.Sqrt(lng*lng+lat*lat) + 0.00002*math.Sin(lat*3.14159265358979323846264338327950288*3000.0/180.0)
	theta := math.Atan2(lat, lng) + 0.000003*math.Cos(lng*3.14159265358979323846264338327950288*3000.0/180.0)

	return z*math.Sin(theta) + 0.006, z*math.Cos(theta) + 0.0065
}

// 判断是否在中国境内
func (g *geocodingUtil) outOfChina(lat, lng float64) bool {
	return lng < 72.004 || lng > 137.8347 || lat < 0.8293 || lat > 55.8271
}

// 纬度转换
func (g *geocodingUtil) transformLat(lng, lat float64) float64 {
	ret := -100.0 + 2.0*lng + 3.0*lat + 0.2*lat*lat + 0.1*lng*lat + 0.2*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*3.14159265358979323846264338327950288) + 20.0*math.Sin(2.0*lng*3.14159265358979323846264338327950288)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*3.14159265358979323846264338327950288) + 40.0*math.Sin(lat/3.0*3.14159265358979323846264338327950288)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*3.14159265358979323846264338327950288) + 320*math.Sin(lat*3.14159265358979323846264338327950288/30.0)) * 2.0 / 3.0
	return ret
}

// 经度转换
func (g *geocodingUtil) transformLng(lng, lat float64) float64 {
	ret := 300.0 + lng + 2.0*lat + 0.1*lng*lng + 0.1*lng*lat + 0.1*math.Sqrt(math.Abs(lng))
	ret += (20.0*math.Sin(6.0*lng*3.14159265358979323846264338327950288) + 20.0*math.Sin(2.0*lng*3.14159265358979323846264338327950288)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lng*3.14159265358979323846264338327950288) + 40.0*math.Sin(lng/3.0*3.14159265358979323846264338327950288)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lng/12.0*3.14159265358979323846264338327950288) + 300.0*math.Sin(lng/30.0*3.14159265358979323846264338327950288)) * 2.0 / 3.0
	return ret
}

// 获取地址信息（优先使用高德，失败则使用百度）
func (g *geocodingUtil) GetAddress(latitude, longitude float64) string {
	// 首先尝试使用高德地图
	address, err := g.GetAddressByAmap(latitude, longitude)
	if err != nil {
		slog.Error(fmt.Sprintf("高德地图逆地理编码失败: %v", err))
		// 如果高德失败，尝试百度地图
		address, err = g.GetAddressByBaidu(latitude, longitude)
		if err != nil {
			slog.Error(fmt.Sprintf("百度地图逆地理编码也失败: %v", err))
			return fmt.Sprintf("经纬度: %.6f,%.6f", latitude, longitude)
		}
	}

	if address == "" {
		return fmt.Sprintf("经纬度: %.6f,%.6f", latitude, longitude)
	}

	return address
}
