// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameTLinkOsStat = "t_link_os_stats"

// TLinkOsStat mapped from table <t_link_os_stats>
type TLinkOsStat struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement:true;comment:ID" json:"id"` // ID
	FullShortURL string    `gorm:"column:full_short_url;comment:完整短链接" json:"full_short_url"`    // 完整短链接
	Date         time.Time `gorm:"column:date;comment:日期" json:"date"`                           // 日期
	Cnt          int32     `gorm:"column:cnt;comment:访问量" json:"cnt"`                            // 访问量
	Os           string    `gorm:"column:os;comment:操作系统" json:"os"`                             // 操作系统
	CreateTime   time.Time `gorm:"column:create_time;comment:创建时间" json:"create_time"`           // 创建时间
	UpdateTime   time.Time `gorm:"column:update_time;comment:修改时间" json:"update_time"`           // 修改时间
	DelFlag      bool      `gorm:"column:del_flag;comment:删除标识 0：未删除 1：已删除" json:"del_flag"`     // 删除标识 0：未删除 1：已删除
}

// TableName TLinkOsStat's table name
func (*TLinkOsStat) TableName() string {
	return TableNameTLinkOsStat
}
