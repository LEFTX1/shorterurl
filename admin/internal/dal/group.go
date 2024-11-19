package dal

import (
	"go-zero-shorterurl/admin/internal/dal/model"
	"go-zero-shorterurl/admin/internal/dal/query"
)

type GroupRepo struct {
	q *query.Query
}

func NewGroupRepo() *GroupRepo {
	return &GroupRepo{
		q: GetQuery(),
	}
}

// Create 创建分组
func (r *GroupRepo) Create(group *model.Group) error {
	return DB.Create(group).Error
}

// FindByUsername 查找用户的所有分组
func (r *GroupRepo) FindByUsername(username string) ([]*model.Group, error) {
	var groups []*model.Group
	err := DB.Where("username = ? AND del_flag = 0", username).
		Order("sort_order DESC, update_time DESC").
		Find(&groups).Error
	return groups, err
}

// UpdateSort 更新分组排序
func (r *GroupRepo) UpdateSort(username string, gid string, sortOrder int32) error {
	return DB.Exec("UPDATE t_group SET sort_order = ? WHERE username = ? AND gid = ?",
		sortOrder, username, gid).Error
}

// Delete 删除分组（软删除）
func (r *GroupRepo) Delete(username string, gid string) error {
	return DB.Exec("UPDATE t_group SET del_flag = 1 WHERE username = ? AND gid = ?",
		username, gid).Error
}
