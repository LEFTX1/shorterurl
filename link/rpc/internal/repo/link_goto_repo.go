package repo

import (
	"context"
	"shorterurl/link/rpc/internal/model"

	"gorm.io/gorm"
)

// LinkGotoRepo 链接跳转仓库接口
type LinkGotoRepo interface {
	// 根据完整短链接查找链接跳转
	FindByFullShortUrl(ctx context.Context, fullShortUrl string) (*model.LinkGoto, error)
	// 创建链接跳转
	Create(ctx context.Context, linkGoto *model.LinkGoto) error
	// 删除链接跳转
	Delete(ctx context.Context, id interface{}) error
	// 批量创建链接跳转
	BatchCreate(ctx context.Context, linkGotos []*model.LinkGoto) error
	// 根据分组ID和完整短链接删除记录
	DeleteByGidAndFullShortUrl(ctx context.Context, gid string, fullShortUrl string) error
}

// linkGotoRepo 链接跳转仓库实现
type linkGotoRepo struct {
	db *gorm.DB
}

// NewLinkGotoRepo 创建链接跳转仓库
func NewLinkGotoRepo(db *gorm.DB) LinkGotoRepo {
	return &linkGotoRepo{
		db: db,
	}
}

// FindByFullShortUrl 根据完整短链接查找链接跳转
// 注意：full_short_url是分片键，这个查询会被正确路由到对应的分片
func (r *linkGotoRepo) FindByFullShortUrl(ctx context.Context, fullShortUrl string) (*model.LinkGoto, error) {
	var linkGoto model.LinkGoto
	err := r.db.WithContext(ctx).
		Where("full_short_url = ?", fullShortUrl).
		First(&linkGoto).Error
	if err != nil {
		return nil, err
	}
	return &linkGoto, nil
}

// Create 创建链接跳转
func (r *linkGotoRepo) Create(ctx context.Context, linkGoto *model.LinkGoto) error {
	return r.db.WithContext(ctx).Create(linkGoto).Error
}

// Delete 删除链接跳转
func (r *linkGotoRepo) Delete(ctx context.Context, id interface{}) error {
	if fullShortUrl, ok := id.(string); ok {
		return r.db.WithContext(ctx).
			Where("full_short_url = ?", fullShortUrl).
			Delete(&model.LinkGoto{}).Error
	}
	var linkGoto model.LinkGoto
	if err := r.db.WithContext(ctx).First(&linkGoto, id).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Where("id = ? AND full_short_url = ?", linkGoto.ID, linkGoto.FullShortUrl).
		Delete(&model.LinkGoto{}).Error
}

// BatchCreate 批量创建链接跳转
func (r *linkGotoRepo) BatchCreate(ctx context.Context, linkGotos []*model.LinkGoto) error {
	return r.db.WithContext(ctx).CreateInBatches(linkGotos, 100).Error
}

// DeleteByGidAndFullShortUrl 根据分组ID和完整短链接删除记录
// 注意：full_short_url是分片键，这个删除操作会被正确路由到对应的分片
func (r *linkGotoRepo) DeleteByGidAndFullShortUrl(ctx context.Context, gid string, fullShortUrl string) error {
	return r.db.WithContext(ctx).
		Where("gid = ? AND full_short_url = ?", gid, fullShortUrl).
		Delete(&model.LinkGoto{}).Error
}
