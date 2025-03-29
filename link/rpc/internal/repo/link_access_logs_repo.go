package repo

import (
	"context"
	"fmt"
	"shorterurl/link/rpc/internal/model" // 确保引入 model 包
	"time"

	"gorm.io/gorm"
)

// LinkAccessLogDO 链接访问日志数据对象
type LinkAccessLogDO struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	FullShortUrl string    `gorm:"column:full_short_url"`
	User         string    `gorm:"column:user"`
	Ip           string    `gorm:"column:ip"`
	Browser      string    `gorm:"column:browser"`
	Os           string    `gorm:"column:os"`
	Device       string    `gorm:"column:device"`
	Network      string    `gorm:"column:network"`
	Locale       string    `gorm:"column:locale"`
	CreateTime   time.Time `gorm:"column:create_time"`
}

// TableName 表名
func (LinkAccessLogDO) TableName() string {
	return "t_link_access_logs"
}

// LinkAccessLogsRepo 链接访问日志仓库接口
type LinkAccessLogsRepo interface {
	// ListTopIpByShortLink 获取短链接高频访问IP详情
	ListTopIpByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]map[string]interface{}, error)

	// ListTopIpByGroup 获取分组高频访问IP详情
	ListTopIpByGroup(ctx context.Context, gid, startDate, endDate string) ([]map[string]interface{}, error)

	// FindUvTypeCntByShortLink 获取短链接访客类型统计
	FindUvTypeCntByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]map[string]interface{}, error)

	// FindUvTypeCntByGroup 获取分组访客类型统计
	FindUvTypeCntByGroup(ctx context.Context, gid, startDate, endDate string) ([]map[string]interface{}, error)

	// PageLinkAccessLogs 分页查询短链接访问记录
	PageLinkAccessLogs(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32, current, size int64) ([]*LinkAccessLogDO, int64, error)

	// PageGroupAccessLogs 分页查询分组访问记录
	PageGroupAccessLogs(ctx context.Context, gid, startDate, endDate string, current, size int64) ([]*LinkAccessLogDO, int64, error)

	// SelectUvTypeByUsers 查询用户的访客类型
	SelectUvTypeByUsers(ctx context.Context, gid, fullShortUrl string, enableStatus int32, startDate, endDate string, userList []string) ([]map[string]interface{}, error)

	// SelectGroupUvTypeByUsers 查询分组用户的访客类型
	SelectGroupUvTypeByUsers(ctx context.Context, gid, startDate, endDate string, userList []string) ([]map[string]interface{}, error)
}

// linkAccessLogsRepo 链接访问日志仓库实现
type linkAccessLogsRepo struct {
	db     *gorm.DB // common db
	linkDB *gorm.DB // link db for querying t_link
}

// NewLinkAccessLogsRepo 创建链接访问日志仓库
func NewLinkAccessLogsRepo(db *gorm.DB, linkDB *gorm.DB) LinkAccessLogsRepo {
	return &linkAccessLogsRepo{
		db:     db,
		linkDB: linkDB,
	}
}

// ListTopIpByShortLink 获取短链接高频访问IP详情
func (r *linkAccessLogsRepo) ListTopIpByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := r.db.WithContext(ctx).Table("t_link_access_logs")
	query = query.Select("ip, COUNT(*) as count")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		endTime, _ := time.Parse("2006-01-02", endDate)
		endTime = endTime.Add(24 * time.Hour)
		query = query.Where("create_time >= ? AND create_time < ?", startTime, endTime)
	}

	query = query.Group("ip").Order("count DESC").Limit(5)

	err := query.Find(&results).Error
	return results, err
}

// ListTopIpByGroup 获取分组高频访问IP详情
func (r *linkAccessLogsRepo) ListTopIpByGroup(ctx context.Context, gid, startDate, endDate string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

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
	query := r.db.WithContext(ctx).Table(LinkAccessLogDO{}.TableName())
	query = query.Select("ip, COUNT(*) as count")
	query = query.Where("full_short_url IN (?)", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		endTime, _ := time.Parse("2006-01-02", endDate)
		endTime = endTime.Add(24 * time.Hour)
		query = query.Where("create_time >= ? AND create_time < ?", startTime, endTime)
	}

	query = query.Group("ip").Order("count DESC").Limit(5)

	err = query.Find(&results).Error
	return results, err
}

// FindUvTypeCntByShortLink 获取短链接访客类型统计
func (r *linkAccessLogsRepo) FindUvTypeCntByShortLink(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// 构建子查询，查找每个用户首次访问的时间
	subQuery := r.db.Table("t_link_access_logs").
		Select("user, MIN(create_time) as first_time").
		Where("full_short_url = ?", fullShortUrl).
		Group("user")

	// 日期过滤
	if startDate != "" && endDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		endTime, _ := time.Parse("2006-01-02", endDate)
		endTime = endTime.Add(24 * time.Hour)
		subQuery = subQuery.Where("create_time >= ? AND create_time < ?", startTime, endTime)
	}

	// 构建主查询，统计各类型访客数量
	query := r.db.WithContext(ctx).Table("t_link_access_logs as a").
		Select("CASE WHEN DATE(a.create_time) = DATE(first.first_time) THEN '新访客' ELSE '旧访客' END as uv_type, COUNT(DISTINCT a.user) as count").
		Joins("JOIN (?) as first ON a.user = first.user", subQuery).
		Where("a.full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		endTime, _ := time.Parse("2006-01-02", endDate)
		endTime = endTime.Add(24 * time.Hour)
		query = query.Where("a.create_time >= ? AND a.create_time < ?", startTime, endTime)
	}

	query = query.Group("uv_type")

	err := query.Find(&results).Error
	return results, err
}

// FindUvTypeCntByGroup 获取分组访客类型统计
func (r *linkAccessLogsRepo) FindUvTypeCntByGroup(ctx context.Context, gid, startDate, endDate string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

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

	// 2. 构建子查询，查找每个用户首次访问的时间
	subQuery := r.db.Table("t_link_access_logs").
		Select("user, MIN(create_time) as first_time").
		Where("full_short_url IN (?)", fullShortUrls).
		Group("user")

	// 日期过滤
	if startDate != "" && endDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		endTime, _ := time.Parse("2006-01-02", endDate)
		endTime = endTime.Add(24 * time.Hour)
		subQuery = subQuery.Where("create_time >= ? AND create_time < ?", startTime, endTime)
	}

	// 3. 构建主查询，统计各类型访客数量
	query := r.db.WithContext(ctx).Table("t_link_access_logs as a").
		Select("CASE WHEN DATE(a.create_time) = DATE(first.first_time) THEN '新访客' ELSE '旧访客' END as uv_type, COUNT(DISTINCT a.user) as count").
		Joins("JOIN (?) as first ON a.user = first.user", subQuery).
		Where("a.full_short_url IN (?)", fullShortUrls)

	// 日期过滤
	if startDate != "" && endDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		endTime, _ := time.Parse("2006-01-02", endDate)
		endTime = endTime.Add(24 * time.Hour)
		query = query.Where("a.create_time >= ? AND a.create_time < ?", startTime, endTime)
	}

	query = query.Group("uv_type")

	err = query.Find(&results).Error
	return results, err
}

// PageLinkAccessLogs 分页查询短链接访问记录
func (r *linkAccessLogsRepo) PageLinkAccessLogs(ctx context.Context, fullShortUrl, gid, startDate, endDate string, enableStatus int32, current, size int64) ([]*LinkAccessLogDO, int64, error) {
	var records []*LinkAccessLogDO
	var count int64

	query := r.db.WithContext(ctx).Table("t_link_access_logs")
	query = query.Where("full_short_url = ?", fullShortUrl)

	// 日期过滤
	if startDate != "" && endDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		endTime, _ := time.Parse("2006-01-02", endDate)
		endTime = endTime.Add(24 * time.Hour)
		query = query.Where("create_time >= ? AND create_time < ?", startTime, endTime)
	}

	// 查询总数
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	offset := (current - 1) * size
	err = query.Offset(int(offset)).Limit(int(size)).Order("create_time DESC").Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, count, nil
}

// PageGroupAccessLogs 分页查询分组访问记录
func (r *linkAccessLogsRepo) PageGroupAccessLogs(ctx context.Context, gid, startDate, endDate string, current, size int64) ([]*LinkAccessLogDO, int64, error) {
	// 1. 从 LinkDB 查询 gid 对应的 full_short_url 列表
	var fullShortUrls []string
	err := r.linkDB.WithContext(ctx).
		Model(&model.Link{}).
		Where("gid = ? AND del_flag = 0", gid).
		Pluck("full_short_url", &fullShortUrls).Error
	if err != nil {
		return nil, 0, err
	}
	if len(fullShortUrls) == 0 {
		return []*LinkAccessLogDO{}, 0, nil
	}

	// 2. 使用 fullShortUrls 列表在 CommonDB 查询总数和分页数据
	var logs []*LinkAccessLogDO
	var total int64

	tx := r.db.WithContext(ctx).Model(&LinkAccessLogDO{}) // 使用 Model(&LinkAccessLogDO{}) 指定模型
	tx = tx.Where("full_short_url IN (?)", fullShortUrls)

	// 添加时间范围过滤，将日期字符串转换为 'YYYY-MM-DD HH:MM:SS' 格式
	startDateTime := fmt.Sprintf("%s 00:00:00", startDate)
	endDateTime := fmt.Sprintf("%s 23:59:59", endDate) // 包含结束日期当天
	tx = tx.Where("create_time >= ? AND create_time <= ?", startDateTime, endDateTime)

	// 先查询总数
	err = tx.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 再进行分页查询
	offset := (current - 1) * size
	err = tx.Offset(int(offset)).Limit(int(size)).Order("create_time DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// SelectUvTypeByUsers 查询用户的访客类型
func (r *linkAccessLogsRepo) SelectUvTypeByUsers(ctx context.Context, gid, fullShortUrl string, enableStatus int32, startDate, endDate string, userList []string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	if len(userList) == 0 {
		return results, nil
	}

	// 构建子查询，查找每个用户首次访问的时间
	subQuery := r.db.Table("t_link_access_logs").
		Select("user, MIN(create_time) as first_time").
		Where("full_short_url = ?", fullShortUrl).
		Group("user")

	// 用户列表过滤
	subQuery = subQuery.Where("user IN ?", userList)

	// 构建主查询，获取用户类型
	query := r.db.WithContext(ctx).Table("t_link_access_logs as a").
		Select("a.user, CASE WHEN DATE(a.create_time) = DATE(first.first_time) THEN '新访客' ELSE '旧访客' END as uvType").
		Joins("JOIN (?) as first ON a.user = first.user", subQuery).
		Where("a.full_short_url = ?", fullShortUrl).
		Where("a.user IN ?", userList)

	// 日期过滤
	if startDate != "" && endDate != "" {
		startTime, _ := time.Parse("2006-01-02", startDate)
		endTime, _ := time.Parse("2006-01-02", endDate)
		endTime = endTime.Add(24 * time.Hour)
		query = query.Where("a.create_time >= ? AND a.create_time < ?", startTime, endTime)
	}

	query = query.Group("a.user, uvType")

	err := query.Find(&results).Error
	return results, err
}

// SelectGroupUvTypeByUsers 查询分组用户的访客类型
func (r *linkAccessLogsRepo) SelectGroupUvTypeByUsers(ctx context.Context, gid, startDate, endDate string, userList []string) ([]map[string]interface{}, error) {
	// 1. 从 LinkDB 查询 gid 对应的 full_short_url 列表
	var fullShortUrls []string
	err := r.linkDB.WithContext(ctx).
		Model(&model.Link{}).
		Where("gid = ? AND del_flag = 0", gid).
		Pluck("full_short_url", &fullShortUrls).Error
	if err != nil {
		return nil, err
	}
	if len(fullShortUrls) == 0 || len(userList) == 0 {
		return []map[string]interface{}{}, nil
	}

	// 2. 使用 fullShortUrls 和 userList 在 CommonDB 查询访客类型
	var results []map[string]interface{}

	subQuery := r.db.Table(LinkAccessLogDO{}.TableName()).
		Select("MIN(id) as id").
		Where("full_short_url IN (?)", fullShortUrls).
		Where("user IN (?)", userList)

	// 添加时间范围过滤
	startDateTime := fmt.Sprintf("%s 00:00:00", startDate)
	endDateTime := fmt.Sprintf("%s 23:59:59", endDate)
	subQuery = subQuery.Where("create_time >= ? AND create_time <= ?", startDateTime, endDateTime)

	subQuery = subQuery.Group("full_short_url, user")

	query := r.db.WithContext(ctx).Table(fmt.Sprintf("%s as t1", LinkAccessLogDO{}.TableName())).
		Select("t1.user, CASE WHEN COUNT(DISTINCT DATE(t1.create_time)) > 1 THEN '旧访客' ELSE '新访客' END AS uvType")

	// 修正 Join 子句的参数传递方式
	query = query.Joins("JOIN (?) as t2 ON t1.id = t2.id", subQuery)
	query = query.Group("t1.user")

	err = query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
