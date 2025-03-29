package repo

import (
	"context"
	"shorterurl/link/rpc/internal/model" // 确保引入 model 包
	"time"

	"gorm.io/gorm"
)

// LinkDeviceStatsDO 链接设备统计数据对象
type LinkDeviceStatsDO struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	FullShortUrl string    `gorm:"column:full_short_url"`
	Date         time.Time `gorm:"column:date"`
	Device       string    `gorm:"column:device"`
	Cnt          int32     `gorm:"column:cnt"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

// TableName 表名
func (LinkDeviceStatsDO) TableName() string {
	return "t_link_device_stats"
}

// LinkDeviceStatsRepo 链接设备统计仓库接口
type LinkDeviceStatsRepo interface {
	// ListDeviceStatsByShortLink 获取短链接设备访问详情
	ListDeviceStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*DeviceStats, error)

	// ListDeviceStatsByGroup 获取分组设备访问详情
	ListDeviceStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]*DeviceStats, error)
}

// DeviceStats 设备统计
type DeviceStats struct {
	Device string
	Cnt    int32
}

// linkDeviceStatsRepo 链接设备统计仓库实现
type linkDeviceStatsRepo struct {
	db     *gorm.DB // common db
	linkDB *gorm.DB // link db for querying t_link
}

// NewLinkDeviceStatsRepo 创建链接设备统计仓库
func NewLinkDeviceStatsRepo(db *gorm.DB, linkDB *gorm.DB) LinkDeviceStatsRepo {
	return &linkDeviceStatsRepo{
		db:     db,
		linkDB: linkDB,
	}
}

// ListDeviceStatsByShortLink 获取短链接设备访问详情
func (r *linkDeviceStatsRepo) ListDeviceStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*DeviceStats, error) {
	var results []*DeviceStats

	query := r.db.WithContext(ctx).Table("t_link_device_stats")
	query = query.Select("device, IFNULL(SUM(cnt), 0) as cnt")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("device").Order("cnt DESC")

	err := query.Find(&results).Error
	return results, err
}

// ListDeviceStatsByGroup 获取分组设备访问详情
func (r *linkDeviceStatsRepo) ListDeviceStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]*DeviceStats, error) {
	// 1. 从 LinkDB 查询 gid 对应的 full_short_url 列表
	var fullShortUrls []string
	err := r.linkDB.WithContext(ctx).
		Model(&model.Link{}).
		Where("gid = ? AND del_flag = 0", gid).
		Pluck("full_short_url", &fullShortUrls).Error
	if err != nil {
		return nil, err
	}
	if len(fullShortUrls) == 0 {
		return []*DeviceStats{}, nil
	}

	// 2. 使用 fullShortUrls 列表在 CommonDB 查询统计数据
	var results []*DeviceStats
	query := r.db.WithContext(ctx).Table(LinkDeviceStatsDO{}.TableName())
	query = query.Select("device, IFNULL(SUM(cnt), 0) as cnt")
	query = query.Where("full_short_url IN (?)", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("device").Order("cnt DESC")

	err = query.Find(&results).Error
	return results, err
}
