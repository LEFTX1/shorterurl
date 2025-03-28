package model

import (
	"time"
)

// Link 短链接表模型
type Link struct {
	ID            int64     `gorm:"primaryKey;column:id;comment:ID"`
	Domain        string    `gorm:"column:domain;comment:域名"`
	ShortUri      string    `gorm:"column:short_uri;comment:短链接"`
	FullShortUrl  string    `gorm:"column:full_short_url;comment:完整短链接;index"`
	OriginUrl     string    `gorm:"column:origin_url;comment:原始链接"`
	ClickNum      int       `gorm:"column:click_num;default:0;comment:点击量"`
	Gid           string    `gorm:"column:gid;default:default;comment:分组标识;index"`
	Favicon       string    `gorm:"column:favicon;comment:网站图标"`
	EnableStatus  int       `gorm:"column:enable_status;comment:启用标识 0：启用 1：未启用"`
	CreatedType   int       `gorm:"column:created_type;comment:创建类型 0：接口创建 1：控制台创建"`
	ValidDateType int       `gorm:"column:valid_date_type;comment:有效期类型 0：永久有效 1：自定义"`
	ValidDate     time.Time `gorm:"column:valid_date;comment:有效期"`
	Describe      string    `gorm:"column:describe;comment:描述"`
	TotalPv       int       `gorm:"column:total_pv;comment:历史PV"`
	TotalUv       int       `gorm:"column:total_uv;comment:历史UV"`
	TotalUip      int       `gorm:"column:total_uip;comment:历史UIP"`
	CreateTime    time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime    time.Time `gorm:"column:update_time;comment:更新时间"`
	DelTime       int64     `gorm:"column:del_time;default:0;comment:删除时间戳"`
	DelFlag       int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (Link) TableName() string {
	return "t_link"
}

// LinkGoto 短链接跳转表模型
type LinkGoto struct {
	ID           int64  `gorm:"primaryKey;column:id;comment:ID"`
	Gid          string `gorm:"column:gid;default:default;comment:分组标识"`
	FullShortUrl string `gorm:"column:full_short_url;comment:完整短链接;index"`
}

// TableName 表名
func (LinkGoto) TableName() string {
	return "t_link_goto"
}

// Group 分组表模型
type Group struct {
	ID         int64     `gorm:"primaryKey;column:id;comment:ID"`
	Gid        string    `gorm:"column:gid;comment:分组标识"`
	Name       string    `gorm:"column:name;comment:分组名称"`
	Username   string    `gorm:"column:username;comment:创建分组用户名;index"`
	SortOrder  int       `gorm:"column:sort_order;comment:分组排序"`
	CreateTime time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime time.Time `gorm:"column:update_time;comment:更新时间"`
	DelFlag    int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (Group) TableName() string {
	return "t_group"
}

// User 用户表模型
type User struct {
	ID           int64     `gorm:"primaryKey;column:id;comment:ID"`
	Username     string    `gorm:"column:username;comment:用户名;index"`
	Password     string    `gorm:"column:password;comment:密码"`
	RealName     string    `gorm:"column:real_name;comment:真实姓名"`
	Phone        string    `gorm:"column:phone;comment:手机号"`
	Mail         string    `gorm:"column:mail;comment:邮箱"`
	DeletionTime int64     `gorm:"column:deletion_time;comment:注销时间戳"`
	CreateTime   time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime   time.Time `gorm:"column:update_time;comment:更新时间"`
	DelFlag      int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (User) TableName() string {
	return "t_user"
}

// LinkAccessLog 链接访问日志表模型
type LinkAccessLog struct {
	ID           int64     `gorm:"primaryKey;column:id;comment:ID"`
	FullShortUrl string    `gorm:"column:full_short_url;comment:完整短链接;index"`
	User         string    `gorm:"column:user;comment:用户信息"`
	IP           string    `gorm:"column:ip;comment:IP"`
	Browser      string    `gorm:"column:browser;comment:浏览器"`
	Os           string    `gorm:"column:os;comment:操作系统"`
	Network      string    `gorm:"column:network;comment:访问网络"`
	Device       string    `gorm:"column:device;comment:访问设备"`
	Locale       string    `gorm:"column:locale;comment:地区"`
	CreateTime   time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime   time.Time `gorm:"column:update_time;comment:更新时间"`
	DelFlag      int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (LinkAccessLog) TableName() string {
	return "t_link_access_logs"
}

// LinkAccessStats 链接访问统计表模型
type LinkAccessStats struct {
	ID           int64     `gorm:"primaryKey;column:id;comment:ID"`
	FullShortUrl string    `gorm:"column:full_short_url;comment:完整短链接"`
	Date         time.Time `gorm:"column:date;comment:日期"`
	PV           int       `gorm:"column:pv;comment:访问量"`
	UV           int       `gorm:"column:uv;comment:独立访客数"`
	UIP          int       `gorm:"column:uip;comment:独立IP数"`
	Hour         int       `gorm:"column:hour;comment:小时"`
	Weekday      int       `gorm:"column:weekday;comment:星期"`
	CreateTime   time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime   time.Time `gorm:"column:update_time;comment:更新时间"`
	DelFlag      int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (LinkAccessStats) TableName() string {
	return "t_link_access_stats"
}

// LinkBrowserStats 链接浏览器统计表模型
type LinkBrowserStats struct {
	ID           int64     `gorm:"primaryKey;column:id;comment:ID"`
	FullShortUrl string    `gorm:"column:full_short_url;comment:完整短链接"`
	Date         time.Time `gorm:"column:date;comment:日期"`
	Cnt          int       `gorm:"column:cnt;comment:访问量"`
	Browser      string    `gorm:"column:browser;comment:浏览器"`
	CreateTime   time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime   time.Time `gorm:"column:update_time;comment:更新时间"`
	DelFlag      int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (LinkBrowserStats) TableName() string {
	return "t_link_browser_stats"
}

// LinkDeviceStats 链接设备统计表模型
type LinkDeviceStats struct {
	ID           int64     `gorm:"primaryKey;column:id;comment:ID"`
	FullShortUrl string    `gorm:"column:full_short_url;comment:完整短链接"`
	Date         time.Time `gorm:"column:date;comment:日期"`
	Cnt          int       `gorm:"column:cnt;comment:访问量"`
	Device       string    `gorm:"column:device;comment:访问设备"`
	CreateTime   time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime   time.Time `gorm:"column:update_time;comment:更新时间"`
	DelFlag      int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (LinkDeviceStats) TableName() string {
	return "t_link_device_stats"
}

// LinkLocaleStats 链接地区统计表模型
type LinkLocaleStats struct {
	ID           int64     `gorm:"primaryKey;column:id;comment:ID"`
	FullShortUrl string    `gorm:"column:full_short_url;comment:完整短链接"`
	Date         time.Time `gorm:"column:date;comment:日期"`
	Cnt          int       `gorm:"column:cnt;comment:访问量"`
	Province     string    `gorm:"column:province;comment:省份名称"`
	City         string    `gorm:"column:city;comment:市名称"`
	Adcode       string    `gorm:"column:adcode;comment:城市编码"`
	Country      string    `gorm:"column:country;comment:国家标识"`
	CreateTime   time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime   time.Time `gorm:"column:update_time;comment:更新时间"`
	DelFlag      int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (LinkLocaleStats) TableName() string {
	return "t_link_locale_stats"
}

// LinkNetworkStats 链接网络统计表模型
type LinkNetworkStats struct {
	ID           int64     `gorm:"primaryKey;column:id;comment:ID"`
	FullShortUrl string    `gorm:"column:full_short_url;comment:完整短链接"`
	Date         time.Time `gorm:"column:date;comment:日期"`
	Cnt          int       `gorm:"column:cnt;comment:访问量"`
	Network      string    `gorm:"column:network;comment:访问网络"`
	CreateTime   time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime   time.Time `gorm:"column:update_time;comment:更新时间"`
	DelFlag      int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (LinkNetworkStats) TableName() string {
	return "t_link_network_stats"
}

// LinkOsStats 链接操作系统统计表模型
type LinkOsStats struct {
	ID           int64     `gorm:"primaryKey;column:id;comment:ID"`
	FullShortUrl string    `gorm:"column:full_short_url;comment:完整短链接"`
	Date         time.Time `gorm:"column:date;comment:日期"`
	Cnt          int       `gorm:"column:cnt;comment:访问量"`
	Os           string    `gorm:"column:os;comment:操作系统"`
	CreateTime   time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime   time.Time `gorm:"column:update_time;comment:更新时间"`
	DelFlag      int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (LinkOsStats) TableName() string {
	return "t_link_os_stats"
}

// LinkStatsToday 链接当日统计表模型
type LinkStatsToday struct {
	ID           int64     `gorm:"primaryKey;column:id;comment:ID"`
	FullShortUrl string    `gorm:"column:full_short_url;comment:短链接"`
	Date         time.Time `gorm:"column:date;comment:日期"`
	TodayPV      int       `gorm:"column:today_pv;default:0;comment:今日PV"`
	TodayUV      int       `gorm:"column:today_uv;default:0;comment:今日UV"`
	TodayUIP     int       `gorm:"column:today_uip;default:0;comment:今日IP数"`
	CreateTime   time.Time `gorm:"column:create_time;comment:创建时间"`
	UpdateTime   time.Time `gorm:"column:update_time;comment:更新时间"`
	DelFlag      int       `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除"`
}

// TableName 表名
func (LinkStatsToday) TableName() string {
	return "t_link_stats_today"
}

// GroupUnique 分组唯一标识表
type GroupUnique struct {
	ID  int64  `gorm:"primaryKey;column:id;comment:ID"`
	Gid string `gorm:"column:gid;comment:分组标识;index:idx_unique_gid,unique"`
}

// TableName 表名
func (GroupUnique) TableName() string {
	return "t_group_unique"
}
