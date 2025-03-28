package repo

import (
	"context"
	"shorterurl/link/rpc/internal/model"

	"gorm.io/gorm"
)

// GroupRepo 分组仓库接口
type GroupRepo interface {
	// 创建分组
	Create(ctx context.Context, group *model.Group) error

	// 根据ID查询分组
	FindByID(ctx context.Context, id int64) (*model.Group, error)

	// 根据用户名和分组名查询分组
	FindByUsernameAndName(ctx context.Context, username, name string) (*model.Group, error)

	// 根据用户名查询分组列表
	FindByUsername(ctx context.Context, username string, page, pageSize int) ([]*model.Group, int64, error)

	// 更新分组
	Update(ctx context.Context, group *model.Group) error

	// 删除分组
	Delete(ctx context.Context, id int64) error
}

// groupRepo 分组仓库实现
type groupRepo struct {
	db *gorm.DB
}

// NewGroupRepo 创建分组仓库
func NewGroupRepo(db *gorm.DB) GroupRepo {
	return &groupRepo{
		db: db,
	}
}

// Create 创建分组
func (r *groupRepo) Create(ctx context.Context, group *model.Group) error {
	return r.db.WithContext(ctx).Create(group).Error
}

// FindByID 根据ID查询分组
func (r *groupRepo) FindByID(ctx context.Context, id int64) (*model.Group, error) {
	var group model.Group
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Where("del_flag = ?", 0).
		First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// FindByUsernameAndName 根据用户名和分组名查询分组
func (r *groupRepo) FindByUsernameAndName(ctx context.Context, username, name string) (*model.Group, error) {
	var group model.Group
	err := r.db.WithContext(ctx).
		Where("username = ? AND name = ?", username, name).
		Where("del_flag = ?", 0).
		First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// FindByUsername 根据用户名查询分组列表
func (r *groupRepo) FindByUsername(ctx context.Context, username string, page, pageSize int) ([]*model.Group, int64, error) {
	var groups []*model.Group
	var count int64

	// 查询总数
	err := r.db.WithContext(ctx).
		Model(&model.Group{}).
		Where("username = ?", username).
		Where("del_flag = ?", 0).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	offset := (page - 1) * pageSize
	err = r.db.WithContext(ctx).
		Where("username = ?", username).
		Where("del_flag = ?", 0).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&groups).Error
	if err != nil {
		return nil, 0, err
	}

	return groups, count, nil
}

// Update 更新分组
func (r *groupRepo) Update(ctx context.Context, group *model.Group) error {
	return r.db.WithContext(ctx).Save(group).Error
}

// Delete 删除分组
func (r *groupRepo) Delete(ctx context.Context, id int64) error {
	// 这里使用软删除
	return r.db.WithContext(ctx).
		Model(&model.Group{}).
		Where("id = ?", id).
		Update("del_flag", 1).Error
}
