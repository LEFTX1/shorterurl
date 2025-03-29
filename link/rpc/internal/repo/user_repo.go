package repo

import (
	"context"
	"shorterurl/link/rpc/internal/model"
	"time"

	"gorm.io/gorm"
)

// UserRepo 用户仓库接口
type UserRepo interface {
	// 创建用户
	Create(ctx context.Context, user *model.User) error

	// 根据ID查询用户
	FindByID(ctx context.Context, id int64) (*model.User, error)

	// 根据用户名查询用户
	FindByUsername(ctx context.Context, username string) (*model.User, error)

	// 根据邮箱查询用户
	FindByMail(ctx context.Context, mail string) (*model.User, error)

	// 根据手机号查询用户
	FindByPhone(ctx context.Context, phone string) (*model.User, error)

	// 更新用户
	Update(ctx context.Context, user *model.User) error

	// 删除用户 (根据主键 ID，软删除)
	Delete(ctx context.Context, id int64) error

	// 根据用户名删除用户 (硬删除，用于测试清理)
	DeleteByUsername(ctx context.Context, username string) error

	// 更新密码
	UpdatePassword(ctx context.Context, id int64, password string) error

	// 分页查询用户列表
	FindPage(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
}

// userRepo 用户仓库实现
type userRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建用户仓库
func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

// Create 创建用户
func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据ID查询用户
func (r *userRepo) FindByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Where("del_flag = ?", 0).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查询用户
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		Where("del_flag = ?", 0).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByMail 根据邮箱查询用户
func (r *userRepo) FindByMail(ctx context.Context, mail string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("mail = ?", mail).
		Where("del_flag = ?", 0).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone 根据手机号查询用户
func (r *userRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		Where("del_flag = ?", 0).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户 (软删除)
func (r *userRepo) Delete(ctx context.Context, id int64) error {
	// 这里使用软删除
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"del_flag":      1,
			"deletion_time": time.Now().Unix(),
		}).Error
}

// DeleteByUsername 根据用户名删除用户 (硬删除)
func (r *userRepo) DeleteByUsername(ctx context.Context, username string) error {
	// 硬删除，同时使用 username 作为分片键进行路由
	return r.db.WithContext(ctx).Unscoped().Where("username = ?", username).Delete(&model.User{}).Error
}

// UpdatePassword 更新密码
func (r *userRepo) UpdatePassword(ctx context.Context, id int64, password string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("password", password).Error
}

// FindPage 分页查询用户列表
func (r *userRepo) FindPage(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var count int64

	// 查询总数
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("del_flag = ?", 0).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询数据
	offset := (page - 1) * pageSize
	err = r.db.WithContext(ctx).
		Where("del_flag = ?", 0).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}
