package repo

import (
	"context"
	"shorterurl/link/rpc/internal/model" // 确保引入 model 包
	"time"

	"gorm.io/gorm"
)

// LinkBrowserStatsDO 链接浏览器统计数据对象
type LinkBrowserStatsDO struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	FullShortUrl string    `gorm:"column:full_short_url"`
	Date         time.Time `gorm:"column:date"`
	Browser      string    `gorm:"column:browser"`
	Cnt          int32     `gorm:"column:cnt"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

// TableName 表名
func (LinkBrowserStatsDO) TableName() string {
	return "t_link_browser_stats"
}

// LinkBrowserStatsRepo 链接浏览器统计仓库接口
type LinkBrowserStatsRepo interface {
	// ListBrowserStatsByShortLink 获取短链接浏览器访问详情
	ListBrowserStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]map[string]interface{}, error)

	// ListBrowserStatsByGroup 获取分组浏览器访问详情
	ListBrowserStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]map[string]interface{}, error)
}

// linkBrowserStatsRepo 链接浏览器统计仓库实现
type linkBrowserStatsRepo struct {
	db     *gorm.DB // common db
	linkDB *gorm.DB // link db for querying t_link
}

// NewLinkBrowserStatsRepo 创建链接浏览器统计仓库
func NewLinkBrowserStatsRepo(db *gorm.DB, linkDB *gorm.DB) LinkBrowserStatsRepo {
	return &linkBrowserStatsRepo{
		db:     db,
		linkDB: linkDB,
	}
}

// ListBrowserStatsByShortLink 获取短链接浏览器访问详情
func (r *linkBrowserStatsRepo) ListBrowserStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := r.db.WithContext(ctx).Table("t_link_browser_stats")
	query = query.Select("browser, IFNULL(SUM(cnt), 0) as count")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("browser").Order("count DESC")

	err := query.Find(&results).Error
	return results, err
}

// ListBrowserStatsByGroup 获取分组浏览器访问详情
func (r *linkBrowserStatsRepo) ListBrowserStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]map[string]interface{}, error) {
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
	query := r.db.WithContext(ctx).Table(LinkBrowserStatsDO{}.TableName())
	query = query.Select("browser, IFNULL(SUM(cnt), 0) as count")
	query = query.Where("full_short_url IN (?)", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("browser").Order("count DESC")

	err = query.Find(&results).Error
	return results, err
}
