package repo

import (
	"context"
	"shorterurl/link/rpc/internal/model" // 确保引入了 model 包
	"time"

	"gorm.io/gorm"
)

// LinkAccessStatsDO 链接访问统计数据对象
type LinkAccessStatsDO struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	FullShortUrl string    `gorm:"column:full_short_url"`
	Date         time.Time `gorm:"column:date"`
	Hour         int32     `gorm:"column:hour"`
	Weekday      int32     `gorm:"column:weekday"`
	Pv           int32     `gorm:"column:pv"`
	Uv           int32     `gorm:"column:uv"`
	Uip          int32     `gorm:"column:uip"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

// TableName 表名
func (LinkAccessStatsDO) TableName() string {
	return "t_link_access_stats"
}

// LinkAccessStatsRepo 链接访问统计仓库接口
type LinkAccessStatsRepo interface {
	// FindPvUvUipStatsByShortLink 查询单个链接的访问统计数据
	FindPvUvUipStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) (*PvUvUipStats, error)

	// FindPvUvUipStatsByGroup 查询分组的访问统计数据
	FindPvUvUipStatsByGroup(ctx context.Context, gid, startDate, endDate string) (*PvUvUipStats, error)

	// ListStatsByShortLink 获取短链接每日访问详情
	ListStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*DailyStats, error)

	// ListStatsByGroup 获取分组每日访问详情
	ListStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]*DailyStats, error)

	// ListHourStatsByShortLink 获取短链接小时访问详情
	ListHourStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*HourStats, error)

	// ListHourStatsByGroup 获取分组小时访问详情
	ListHourStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]*HourStats, error)

	// ListWeekdayStatsByShortLink 获取短链接一周访问详情
	ListWeekdayStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*WeekdayStats, error)

	// ListWeekdayStatsByGroup 获取分组一周访问详情
	ListWeekdayStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]*WeekdayStats, error)
}

// PvUvUipStats PV、UV、UIP统计结果
type PvUvUipStats struct {
	Pv  int32
	Uv  int32
	Uip int32
}

// DailyStats 每日统计
type DailyStats struct {
	Date time.Time
	Pv   int32
	Uv   int32
	Uip  int32
}

// HourStats 小时统计
type HourStats struct {
	Hour int32
	Pv   int32
}

// WeekdayStats 周统计
type WeekdayStats struct {
	Weekday int32
	Pv      int32
}

// linkAccessStatsRepo 链接访问统计仓库实现
type linkAccessStatsRepo struct {
	db     *gorm.DB // common db
	linkDB *gorm.DB // link db for querying t_link
}

// NewLinkAccessStatsRepo 创建链接访问统计仓库
func NewLinkAccessStatsRepo(db *gorm.DB, linkDB *gorm.DB) LinkAccessStatsRepo {
	return &linkAccessStatsRepo{
		db:     db,
		linkDB: linkDB,
	}
}

// FindPvUvUipStatsByShortLink 查询单个链接的访问统计数据
func (r *linkAccessStatsRepo) FindPvUvUipStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) (*PvUvUipStats, error) {
	var stats struct {
		Pv  int32 `gorm:"column:pv"`
		Uv  int32 `gorm:"column:uv"`
		Uip int32 `gorm:"column:uip"`
	}

	query := r.db.WithContext(ctx).Table("t_link_access_stats")
	query = query.Select("IFNULL(SUM(pv), 0) as pv, IFNULL(SUM(uv), 0) as uv, IFNULL(SUM(uip), 0) as uip")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	err := query.First(&stats).Error
	if err != nil {
		return nil, err
	}

	return &PvUvUipStats{
		Pv:  stats.Pv,
		Uv:  stats.Uv,
		Uip: stats.Uip,
	}, nil
}

// FindPvUvUipStatsByGroup 查询分组的访问统计数据
func (r *linkAccessStatsRepo) FindPvUvUipStatsByGroup(ctx context.Context, gid, startDate, endDate string) (*PvUvUipStats, error) {
	// 1. 从 LinkDB 查询 gid 对应的 full_short_url 列表
	var fullShortUrls []string
	err := r.linkDB.WithContext(ctx).
		Model(&model.Link{}). // 使用正确的 model
		Where("gid = ? AND del_flag = 0", gid).
		Pluck("full_short_url", &fullShortUrls).Error
	if err != nil {
		return nil, err
	}
	if len(fullShortUrls) == 0 {
		// 如果该分组下没有链接，直接返回零值
		return &PvUvUipStats{}, nil
	}

	// 2. 使用 fullShortUrls 列表在 CommonDB 查询统计数据
	var stats PvUvUipStats
	query := r.db.WithContext(ctx).Table(LinkAccessStatsDO{}.TableName()) // 使用 TableName()
	query = query.Select("IFNULL(SUM(pv), 0) as pv, IFNULL(SUM(uv), 0) as uv, IFNULL(SUM(uip), 0) as uip")
	query = query.Where("full_short_url IN (?) ", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	err = query.Take(&stats).Error // 使用 Take 获取单条记录
	if err != nil {
		// 如果查询出错（非 gorm.ErrRecordNotFound），返回错误
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
		// 如果是没找到记录，返回零值统计（已被初始化为零值）
		return &PvUvUipStats{}, nil
	}

	return &stats, nil
}

// ListStatsByShortLink 获取短链接每日访问详情
func (r *linkAccessStatsRepo) ListStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*DailyStats, error) {
	var result []*DailyStats

	query := r.db.WithContext(ctx).Table("t_link_access_stats")
	query = query.Select("date, IFNULL(SUM(pv), 0) as pv, IFNULL(SUM(uv), 0) as uv, IFNULL(SUM(uip), 0) as uip")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("date").Order("date ASC")

	err := query.Find(&result).Error
	return result, err
}

// ListStatsByGroup 获取分组每日访问详情
func (r *linkAccessStatsRepo) ListStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]*DailyStats, error) {
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
		return []*DailyStats{}, nil
	}

	// 2. 使用 fullShortUrls 列表在 CommonDB 查询统计数据
	var results []*DailyStats
	query := r.db.WithContext(ctx).Table(LinkAccessStatsDO{}.TableName())
	query = query.Select("date, SUM(pv) as pv, SUM(uv) as uv, SUM(uip) as uip")
	query = query.Where("full_short_url IN (?)", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("date").Order("date ASC")

	err = query.Scan(&results).Error // 使用 Scan 将结果映射到 []*DailyStats
	return results, err
}

// ListHourStatsByShortLink 获取短链接小时访问详情
func (r *linkAccessStatsRepo) ListHourStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*HourStats, error) {
	var result []*HourStats

	query := r.db.WithContext(ctx).Table("t_link_access_stats")
	query = query.Select("hour, IFNULL(SUM(pv), 0) as pv")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("hour").Order("hour ASC")

	err := query.Find(&result).Error
	return result, err
}

// ListHourStatsByGroup 获取分组小时访问详情
func (r *linkAccessStatsRepo) ListHourStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]*HourStats, error) {
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
		return []*HourStats{}, nil
	}

	// 2. 使用 fullShortUrls 列表在 CommonDB 查询统计数据
	var results []*HourStats
	query := r.db.WithContext(ctx).Table(LinkAccessStatsDO{}.TableName())
	query = query.Select("hour, SUM(pv) as pv")
	query = query.Where("full_short_url IN (?)", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("hour").Order("hour ASC")

	err = query.Scan(&results).Error
	return results, err
}

// ListWeekdayStatsByShortLink 获取短链接一周访问详情
func (r *linkAccessStatsRepo) ListWeekdayStatsByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]*WeekdayStats, error) {
	var result []*WeekdayStats

	query := r.db.WithContext(ctx).Table("t_link_access_stats")
	query = query.Select("weekday, IFNULL(SUM(pv), 0) as pv")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("weekday").Order("weekday ASC")

	err := query.Find(&result).Error
	return result, err
}

// ListWeekdayStatsByGroup 获取分组一周访问详情
func (r *linkAccessStatsRepo) ListWeekdayStatsByGroup(ctx context.Context, gid, startDate, endDate string) ([]*WeekdayStats, error) {
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
		return []*WeekdayStats{}, nil
	}

	// 2. 使用 fullShortUrls 列表在 CommonDB 查询统计数据
	var results []*WeekdayStats
	query := r.db.WithContext(ctx).Table(LinkAccessStatsDO{}.TableName())
	query = query.Select("weekday, SUM(pv) as pv")
	query = query.Where("full_short_url IN (?)", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		query = query.Where("date >= ? AND date <= ?", startDate, endDate)
	}

	query = query.Group("weekday").Order("weekday ASC")

	err = query.Scan(&results).Error
	return results, err
}
