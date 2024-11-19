package model

import "time"

// 通用字段
type Model struct {
	ID         int64     `gorm:"primarykey;autoIncrement:true" json:"id"`
	CreateTime time.Time `gorm:"column:create_time;not null" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;not null" json:"updateTime"`
	DelFlag    int8      `gorm:"column:del_flag;not null;default:0" json:"delFlag"`
}

// User 用户表基础结构
type User struct {
	Model
	Username     string `gorm:"column:username;type:varchar(256);uniqueIndex:idx_unique_username" json:"username"`
	Password     string `gorm:"column:password;type:varchar(512)" json:"password"`
	RealName     string `gorm:"column:real_name;type:varchar(256)" json:"realName"`
	Phone        string `gorm:"column:phone;type:varchar(128)" json:"phone"`
	Mail         string `gorm:"column:mail;type:varchar(512)" json:"mail"`
	DeletionTime int64  `gorm:"column:deletion_time" json:"deletionTime"`
}

// Group 分组表基础结构
type Group struct {
	Model
	Gid       string `gorm:"column:gid;type:varchar(32)" json:"gid"`
	Name      string `gorm:"column:name;type:varchar(64)" json:"name"`
	Username  string `gorm:"column:username;type:varchar(256);index:idx_username" json:"username"`
	SortOrder int32  `gorm:"column:sort_order;default:0" json:"sortOrder"`
}

// TableName 表名
func (User) TableName() string {
	return "t_user"
}

// TableName 表名
func (Group) TableName() string {
	return "t_group"
}
