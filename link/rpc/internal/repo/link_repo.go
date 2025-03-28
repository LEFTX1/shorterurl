package repo

import (
	"context"
	"shorterurl/link/rpc/internal/model"
	"time"

	"gorm.io/gorm"
)

// LinkRepo 短链接仓库接口
type LinkRepo interface {
	// 创建短链接
	Create(ctx context.Context, link *model.Link) error

	// 根据ID查询短链接
	FindByID(ctx context.Context, id int64) (*model.Link, error)

	// 根据短链接查询
	FindByShortUri(ctx context.Context, shortUri string) (*model.Link, error)

	// 根据完整短链接查询
	FindByFullShortUrl(ctx context.Context, fullShortUrl string) (*model.Link, error)

	// 根据短链接和分组ID查询
	FindByShortUriAndGid(ctx context.Context, shortUri, gid string) (*model.Link, error)

	// 根据完整短链接和分组ID查询
	FindByFullShortUrlAndGid(ctx context.Context, fullShortUrl, gid string) (*model.Link, error)

	// 根据分组ID查询短链接列表
	FindByGid(ctx context.Context, gid string, page, pageSize int) ([]*model.Link, int64, error)

	// 查询回收站中的短链接列表
	FindRecycleBin(ctx context.Context, gid string, page, pageSize int) ([]*model.Link, int64, error)

	// 统计分组中的短链接数量
	CountByGid(ctx context.Context, gid string) (int64, error)

	// 更新短链接
	Update(ctx context.Context, link *model.Link) error

	// 删除短链接
	Delete(ctx context.Context, id int64) error

	// 批量创建短链接
	BatchCreate(ctx context.Context, links []*model.Link) error

	// 根据条件查询短链接
	FindByCondition(ctx context.Context, conditions map[string]interface{}, page, pageSize int) ([]*model.Link, error)

	// 根据条件统计短链接数量
	CountByCondition(ctx context.Context, conditions map[string]interface{}) (int64, error)
}

// linkRepo 短链接仓库实现
type linkRepo struct {
	db *gorm.DB
}

// NewLinkRepo 创建短链接仓库
func NewLinkRepo(db *gorm.DB) *linkRepo {
	return &linkRepo{
		db: db,
	}
}

// Create 创建短链接
func (r *linkRepo) Create(ctx context.Context, link *model.Link) error {
	return r.db.WithContext(ctx).Create(link).Error
}

// FindByID 根据ID查询短链接
func (r *linkRepo) FindByID(ctx context.Context, id int64) (*model.Link, error) {
	var link model.Link
	err := r.db.WithContext(ctx).Where("id = ? AND del_flag = 0", id).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// FindByShortUri 根据短链接查询
func (r *linkRepo) FindByShortUri(ctx context.Context, shortUri string) (*model.Link, error) {
	var link model.Link
	// 由于分片键是 gid，我们需要先通过 short_uri 找到对应的记录
	err := r.db.WithContext(ctx).Where("short_uri = ?", shortUri).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// FindByFullShortUrl 根据完整短链接查询
func (r *linkRepo) FindByFullShortUrl(ctx context.Context, fullShortUrl string) (*model.Link, error) {
	var link model.Link
	// 由于分片键是 gid，我们需要先通过 full_short_url 找到对应的记录
	err := r.db.WithContext(ctx).Where("full_short_url = ?", fullShortUrl).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// FindByShortUriAndGid 根据短链接和分组ID查询
func (r *linkRepo) FindByShortUriAndGid(ctx context.Context, shortUri, gid string) (*model.Link, error) {
	var link model.Link
	err := r.db.WithContext(ctx).Where("short_uri = ? AND gid = ?", shortUri, gid).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// FindByFullShortUrlAndGid 根据完整短链接和分组ID查询
func (r *linkRepo) FindByFullShortUrlAndGid(ctx context.Context, fullShortUrl, gid string) (*model.Link, error) {
	var link model.Link
	err := r.db.Where("full_short_url = ? AND gid = ? AND del_flag = 0", fullShortUrl, gid).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// FindByGid 根据分组ID查询短链接列表
func (r *linkRepo) FindByGid(ctx context.Context, gid string, page, pageSize int) ([]*model.Link, int64, error) {
	var links []*model.Link
	var count int64

	// 查询总数
	err := r.db.WithContext(ctx).Model(&model.Link{}).
		Where("gid = ?", gid).
		Where("del_flag = ?", 0).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	offset := (page - 1) * pageSize
	err = r.db.WithContext(ctx).
		Where("gid = ?", gid).
		Where("del_flag = ?", 0).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&links).Error
	if err != nil {
		return nil, 0, err
	}

	return links, count, nil
}

// FindRecycleBin 查询回收站中的短链接列表
func (r *linkRepo) FindRecycleBin(ctx context.Context, gid string, page, pageSize int) ([]*model.Link, int64, error) {
	var links []*model.Link
	var count int64

	// 查询总数
	err := r.db.WithContext(ctx).Model(&model.Link{}).
		Where("gid = ?", gid).
		Where("del_flag = ?", 1).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	offset := (page - 1) * pageSize
	err = r.db.WithContext(ctx).
		Where("gid = ?", gid).
		Where("del_flag = ?", 1).
		Offset(offset).
		Limit(pageSize).
		Order("del_time DESC").
		Find(&links).Error
	if err != nil {
		return nil, 0, err
	}

	return links, count, nil
}

// CountByGid 统计分组中的短链接数量
func (r *linkRepo) CountByGid(ctx context.Context, gid string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Link{}).
		Where("gid = ?", gid).
		Where("del_flag = ?", 0).
		Where("enable_status = ?", 0).
		Count(&count).Error
	return count, err
}

// Update 更新短链接
func (r *linkRepo) Update(ctx context.Context, link *model.Link) error {
	// 注意: GORM的Save会尝试更新所有字段，这里需要显式指定WHERE条件
	// 并且必须传入分片键的值，否则会出现分片错误
	return r.db.WithContext(ctx).
		Table(link.TableName()).
		Where("id = ? AND gid = ?", link.ID, link.Gid). // 使用ID和分片键gid作为条件
		Updates(map[string]interface{}{
			"domain":          link.Domain,
			"short_uri":       link.ShortUri,
			"full_short_url":  link.FullShortUrl,
			"origin_url":      link.OriginUrl,
			"click_num":       link.ClickNum,
			"gid":             link.Gid, // 包含分片键
			"favicon":         link.Favicon,
			"enable_status":   link.EnableStatus,
			"created_type":    link.CreatedType,
			"valid_date_type": link.ValidDateType,
			"valid_date":      link.ValidDate,
			"describe":        link.Describe,
			"total_pv":        link.TotalPv,
			"total_uv":        link.TotalUv,
			"total_uip":       link.TotalUip,
			"create_time":     link.CreateTime,
			"update_time":     link.UpdateTime,
			"del_time":        link.DelTime,
			"del_flag":        link.DelFlag,
		}).Error
}

// Delete 删除短链接
func (r *linkRepo) Delete(ctx context.Context, id int64) error {
	// 这里使用软删除
	return r.db.WithContext(ctx).Model(&model.Link{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"del_flag": 1,
			"del_time": time.Now().Unix(),
		}).Error
}

// BatchCreate 批量创建短链接
func (r *linkRepo) BatchCreate(ctx context.Context, links []*model.Link) error {
	return r.db.WithContext(ctx).CreateInBatches(links, 100).Error
}

// FindByCondition 根据条件查询短链接
func (r *linkRepo) FindByCondition(ctx context.Context, conditions map[string]interface{}, page, pageSize int) ([]*model.Link, error) {
	var links []*model.Link
	err := r.db.WithContext(ctx).Where(conditions).Offset((page - 1) * pageSize).Limit(pageSize).Find(&links).Error
	return links, err
}

// CountByCondition 根据条件统计短链接数量
func (r *linkRepo) CountByCondition(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Link{}).Where(conditions).Count(&count).Error
	return count, err
}
