package repo

import (
	"context"
	"shorterurl/link/rpc/internal/model" // 确保引入 model 包
	"time"

	"gorm.io/gorm"
)

// LinkOsStatsDO 链接操作系统统计数据对象
type LinkOsStatsDO struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	FullShortUrl string    `gorm:"column:full_short_url"`
	Date         time.Time `gorm:"column:date"`
	Os           string    `gorm:"column:os"`
	Cnt          int32     `gorm:"column:cnt"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

// TableName 表名
func (LinkOsStatsDO) TableName() string {
	return "t_link_os_stats"
}

// LinkOsStatsRepo 链接操作系统统计仓库接口
type LinkOsStatsRepo interface {
	// ListOsStatsByShortLink 获取短链接操作系统访问详情
	ListOsStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]map[string]interface{}, error)

	// ListOsStatsByGroup 获取分组操作系统访问详情
	ListOsStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]map[string]interface{}, error)
}

// linkOsStatsRepo 链接操作系统统计仓库实现
type linkOsStatsRepo struct {
	db     *gorm.DB // common db
	linkDB *gorm.DB // link db for querying t_link
}

// NewLinkOsStatsRepo 创建链接操作系统统计仓库
func NewLinkOsStatsRepo(db *gorm.DB, linkDB *gorm.DB) LinkOsStatsRepo {
	return &linkOsStatsRepo{
		db:     db,
		linkDB: linkDB,
	}
}

// ListOsStatsByShortLink 获取短链接操作系统访问详情
func (r *linkOsStatsRepo) ListOsStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := r.db.WithContext(ctx).Table("t_link_os_stats")
	query = query.Select("os, IFNULL(SUM(cnt), 0) as count")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("os").Order("count DESC")

	err := query.Find(&results).Error
	return results, err
}

// ListOsStatsByGroup 获取分组操作系统访问详情
func (r *linkOsStatsRepo) ListOsStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]map[string]interface{}, error) {
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
		return []map[string]interface{}{}, nil
	}

	// 2. 使用 fullShortUrls 列表在 CommonDB 查询统计数据
	var results []map[string]interface{}
	query := r.db.WithContext(ctx).Table(LinkOsStatsDO{}.TableName())
	query = query.Select("os, IFNULL(SUM(cnt), 0) as count")
	query = query.Where("full_short_url IN (?)", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("os").Order("count DESC")

	err = query.Find(&results).Error
	return results, err
}
