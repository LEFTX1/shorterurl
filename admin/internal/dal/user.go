package dal

import (
	"go-zero-shorterurl/admin/internal/dal/model"
	"go-zero-shorterurl/admin/internal/dal/query"
)

type UserRepo struct {
	q *query.Query
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		q: GetQuery(),
	}
}

// Create 创建用户
func (r *UserRepo) Create(user *model.User) error {
	return DB.Create(user).Error
}

// FindByUsername 根据用户名查找用户
func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := DB.Where("username = ? AND del_flag = 0", username).First(&user).Error
	return &user, err
}

// FindByUsernameAndPassword 用户登录验证
func (r *UserRepo) FindByUsernameAndPassword(username, password string) (*model.User, error) {
	var user model.User
	err := DB.Where("username = ? AND password = ? AND del_flag = 0",
		username, password).First(&user).Error
	return &user, err
}

// Update 更新用户信息
func (r *UserRepo) Update(user *model.User) error {
	return DB.Where("username = ?", user.Username).Updates(user).Error
}
