// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameTLinkGoto = "t_link_goto"

// TLinkGoto mapped from table <t_link_goto>
type TLinkGoto struct {
	ID           int64  `gorm:"column:id;primaryKey;autoIncrement:true;comment:ID" json:"id"` // ID
	Gid          string `gorm:"column:gid;default:default;comment:分组标识" json:"gid"`           // 分组标识
	FullShortURL string `gorm:"column:full_short_url;comment:完整短链接" json:"full_short_url"`    // 完整短链接
}

// TableName TLinkGoto's table name
func (*TLinkGoto) TableName() string {
	return TableNameTLinkGoto
}
