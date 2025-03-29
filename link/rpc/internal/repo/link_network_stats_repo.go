package repo

import (
	"context"
	"shorterurl/link/rpc/internal/model" // 确保引入 model 包
	"time"

	"gorm.io/gorm"
)

// LinkNetworkStatsDO 链接网络统计数据对象
type LinkNetworkStatsDO struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	FullShortUrl string    `gorm:"column:full_short_url"`
	Date         time.Time `gorm:"column:date"`
	Network      string    `gorm:"column:network"`
	Cnt          int32     `gorm:"column:cnt"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

// TableName 表名
func (LinkNetworkStatsDO) TableName() string {
	return "t_link_network_stats"
}

// LinkNetworkStatsRepo 链接网络统计仓库接口
type LinkNetworkStatsRepo interface {
	// ListNetworkStatsByShortLink 获取短链接网络访问详情
	ListNetworkStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*NetworkStats, error)

	// ListNetworkStatsByGroup 获取分组网络访问详情
	ListNetworkStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]*NetworkStats, error)
}

// NetworkStats 网络统计
type NetworkStats struct {
	Network string
	Cnt     int32
}

// linkNetworkStatsRepo 链接网络统计仓库实现
type linkNetworkStatsRepo struct {
	db     *gorm.DB // common db
	linkDB *gorm.DB // link db for querying t_link
}

// NewLinkNetworkStatsRepo 创建链接网络统计仓库
func NewLinkNetworkStatsRepo(db *gorm.DB, linkDB *gorm.DB) LinkNetworkStatsRepo {
	return &linkNetworkStatsRepo{
		db:     db,
		linkDB: linkDB,
	}
}

// ListNetworkStatsByShortLink 获取短链接网络访问详情
func (r *linkNetworkStatsRepo) ListNetworkStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*NetworkStats, error) {
	var results []*NetworkStats

	query := r.db.WithContext(ctx).Table("t_link_network_stats")
	query = query.Select("network, IFNULL(SUM(cnt), 0) as cnt")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("network").Order("cnt DESC")

	err := query.Find(&results).Error
	return results, err
}

// ListNetworkStatsByGroup 获取分组网络访问详情
func (r *linkNetworkStatsRepo) ListNetworkStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]*NetworkStats, error) {
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
		return []*NetworkStats{}, nil
	}

	// 2. 使用 fullShortUrls 列表在 CommonDB 查询统计数据
	var results []*NetworkStats
	query := r.db.WithContext(ctx).Table(LinkNetworkStatsDO{}.TableName())
	query = query.Select("network, IFNULL(SUM(cnt), 0) as cnt")
	query = query.Where("full_short_url IN (?)", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("network").Order("cnt DESC")

	err = query.Find(&results).Error
	return results, err
}
