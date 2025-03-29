package repo

import (
	"context"
	"shorterurl/link/rpc/internal/model"
	"time"

	"gorm.io/gorm"
)

// LinkLocaleStatsDO 链接区域统计数据对象
type LinkLocaleStatsDO struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	FullShortUrl string    `gorm:"column:full_short_url"`
	Date         time.Time `gorm:"column:date"`
	Province     string    `gorm:"column:province"`
	Cnt          int32     `gorm:"column:cnt"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

// TableName 表名
func (LinkLocaleStatsDO) TableName() string {
	return "t_link_locale_stats"
}

// LinkLocaleStatsRepo 链接区域统计仓库接口
type LinkLocaleStatsRepo interface {
	// ListLocaleByShortLink 获取短链接地区访问详情
	ListLocaleByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*LocaleStats, error)

	// ListLocaleByGroup 获取分组地区访问详情
	ListLocaleByGroup(ctx context.Context, gid, startDate, endDate string) ([]*LocaleStats, error)
}

// LocaleStats 地区统计
type LocaleStats struct {
	Province string
	Cnt      int32
}

// linkLocaleStatsRepo 链接区域统计仓库实现
type linkLocaleStatsRepo struct {
	db     *gorm.DB
	linkDB *gorm.DB
}

// NewLinkLocaleStatsRepo 创建链接区域统计仓库
func NewLinkLocaleStatsRepo(db *gorm.DB, linkDB *gorm.DB) LinkLocaleStatsRepo {
	return &linkLocaleStatsRepo{
		db:     db,
		linkDB: linkDB,
	}
}

// ListLocaleByShortLink 获取短链接地区访问详情
func (r *linkLocaleStatsRepo) ListLocaleByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*LocaleStats, error) {
	var result []*LocaleStats

	query := r.db.WithContext(ctx).Table("t_link_locale_stats")
	query = query.Select("province, IFNULL(SUM(cnt), 0) as cnt")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("province").Order("cnt DESC")

	err := query.Find(&result).Error
	return result, err
}

// ListLocaleByGroup 获取分组地区访问详情
func (r *linkLocaleStatsRepo) ListLocaleByGroup(ctx context.Context, gid, startDate, endDate string) ([]*LocaleStats, error) {
	var result []*LocaleStats

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
		return []*LocaleStats{}, nil
	}

	// 2. 使用 fullShortUrls 列表在 CommonDB 查询统计数据
	query := r.db.WithContext(ctx).Table(LinkLocaleStatsDO{}.TableName())
	query = query.Select("province, IFNULL(SUM(cnt), 0) as cnt")
	query = query.Where("full_short_url IN (?)", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("province").Order("cnt DESC")

	err = query.Find(&result).Error
	return result, err
}
