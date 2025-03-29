package repo

import (
	"context"
	"errors"
	"shorterurl/link/rpc/internal/model"
	"time"

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

	// 删除分组 (根据主键 ID)
	Delete(ctx context.Context, id int64) error

	// 根据 GID 删除分组 (硬删除，用于测试清理)
	DeleteByGid(ctx context.Context, gid string) error

	// 根据 GID 和 Username 删除分组 (硬删除，用于测试清理)
	DeleteByGidAndUsername(ctx context.Context, gid, username string) error

	// 检查分组是否属于用户
	CheckGroupBelongToUser(ctx context.Context, gid, username string) (bool, error)
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

// Delete 删除分组 (软删除)
func (r *groupRepo) Delete(ctx context.Context, id int64) error {
	// 先查询出完整的分组信息，以获取 username 作为分片键
	group, err := r.FindByID(ctx, id)
	if err != nil {
		// 如果记录未找到，根据业务逻辑可能返回 nil
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}

	// 使用软删除，并包含分片键 username
	return r.db.WithContext(ctx).
		Model(&model.Group{}).                                // 指定模型
		Where("id = ? AND username = ?", id, group.Username). // 明确指定 id 和分片键 username
		Updates(map[string]interface{}{
			"del_flag":      1,
			"deletion_time": time.Now().Unix(), // 假设 Group 模型也有 deletion_time
		}).Error
}

// DeleteByGid 根据 GID 删除分组 (硬删除)
// 警告: 此方法缺少分片键，在生产环境中可能不安全或无效，仅用于测试
func (r *groupRepo) DeleteByGid(ctx context.Context, gid string) error {
	// 尝试先查找 gid 对应的分组获取 username
	var group model.Group
	err := r.db.WithContext(ctx).Where("gid = ?", gid).First(&group).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在，可以视为"删除成功"
			return nil
		}
		return err
	}

	// 有了 username 后，改用 DeleteByGidAndUsername
	return r.DeleteByGidAndUsername(ctx, gid, group.Username)
}

// DeleteByGidAndUsername 根据 GID 和 Username 删除分组 (硬删除)
// 包含分片键，可以正确路由到特定分片
func (r *groupRepo) DeleteByGidAndUsername(ctx context.Context, gid, username string) error {
	// 使用 username 作为分片键进行路由
	return r.db.WithContext(ctx).Unscoped().
		Where("gid = ? AND username = ?", gid, username).
		Delete(&model.Group{}).Error
}

// CheckGroupBelongToUser 检查分组是否属于用户
func (r *groupRepo) CheckGroupBelongToUser(ctx context.Context, gid, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Group{}).
		Where("gid = ? AND username = ?", gid, username).
		Where("del_flag = ?", 0).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
