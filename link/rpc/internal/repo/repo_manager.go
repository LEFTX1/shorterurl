package repo

import (
	"gorm.io/gorm"
)

// DBs 包含所有数据库连接的引用
type DBs struct {
	Common     *gorm.DB
	LinkDB     *gorm.DB
	GotoLinkDB *gorm.DB
	GroupDB    *gorm.DB
	UserDB     *gorm.DB
}

// RepoManager 仓库管理器
type RepoManager struct {
	dbs *DBs

	// 所有仓库
	Link     LinkRepo
	LinkGoto LinkGotoRepo
	Group    GroupRepo
	User     UserRepo
}

// NewRepoManager 创建仓库管理器
func NewRepoManager(common, linkDB, gotoLinkDB, groupDB, userDB *gorm.DB) *RepoManager {
	dbs := &DBs{
		Common:     common,
		LinkDB:     linkDB,
		GotoLinkDB: gotoLinkDB,
		GroupDB:    groupDB,
		UserDB:     userDB,
	}

	return &RepoManager{
		dbs: dbs,

		// 初始化各个仓库
		Link:     NewLinkRepo(dbs.LinkDB),
		LinkGoto: NewLinkGotoRepo(dbs.GotoLinkDB),
		Group:    NewGroupRepo(dbs.GroupDB),
		User:     NewUserRepo(dbs.UserDB),
	}
}

// GetCommonDB 获取通用数据库连接
func (m *RepoManager) GetCommonDB() *gorm.DB {
	return m.dbs.Common
}

// GetLinkDB 获取短链接数据库连接
func (m *RepoManager) GetLinkDB() *gorm.DB {
	return m.dbs.LinkDB
}

// GetLinkGotoDB 获取短链接跳转数据库连接
func (m *RepoManager) GetLinkGotoDB() *gorm.DB {
	return m.dbs.GotoLinkDB
}

// GetGroupDB 获取分组数据库连接
func (m *RepoManager) GetGroupDB() *gorm.DB {
	return m.dbs.GroupDB
}

// GetUserDB 获取用户数据库连接
func (m *RepoManager) GetUserDB() *gorm.DB {
	return m.dbs.UserDB
}
