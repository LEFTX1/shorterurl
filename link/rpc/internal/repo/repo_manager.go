package repo

import (
	"context"
	"errors"

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

// 定义错误
var (
	ErrUserNotLogin = errors.New("用户未登录")
)

// RepoManager 仓库管理器
type RepoManager struct {
	dbs *DBs

	// 所有仓库
	Link             LinkRepo
	LinkGoto         LinkGotoRepo
	Group            GroupRepo
	User             UserRepo
	LinkAccessStats  LinkAccessStatsRepo
	LinkLocaleStats  LinkLocaleStatsRepo
	LinkAccessLogs   LinkAccessLogsRepo
	LinkBrowserStats LinkBrowserStatsRepo
	LinkOsStats      LinkOsStatsRepo
	LinkDeviceStats  LinkDeviceStatsRepo
	LinkNetworkStats LinkNetworkStatsRepo

	// 添加对 LinkDB 的引用，以便传递给需要的 Repo
	linkDB *gorm.DB
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
		// 保存 LinkDB 引用
		linkDB: linkDB,

		// 初始化各个仓库
		Link:             NewLinkRepo(dbs.LinkDB),
		LinkGoto:         NewLinkGotoRepo(dbs.GotoLinkDB),
		Group:            NewGroupRepo(dbs.GroupDB),
		User:             NewUserRepo(dbs.UserDB),
		LinkAccessStats:  NewLinkAccessStatsRepo(dbs.Common, dbs.LinkDB),  // 传递 LinkDB
		LinkLocaleStats:  NewLinkLocaleStatsRepo(dbs.Common, dbs.LinkDB),  // 传递 LinkDB
		LinkAccessLogs:   NewLinkAccessLogsRepo(dbs.Common, dbs.LinkDB),   // 传递 LinkDB
		LinkBrowserStats: NewLinkBrowserStatsRepo(dbs.Common, dbs.LinkDB), // 传递 LinkDB
		LinkOsStats:      NewLinkOsStatsRepo(dbs.Common, dbs.LinkDB),      // 传递 LinkDB
		LinkDeviceStats:  NewLinkDeviceStatsRepo(dbs.Common, dbs.LinkDB),  // 传递 LinkDB
		LinkNetworkStats: NewLinkNetworkStatsRepo(dbs.Common, dbs.LinkDB), // 传递 LinkDB
	}
}

// GetCurrentUsername 获取当前登录用户名
func (m *RepoManager) GetCurrentUsername(ctx context.Context) (string, error) {
	// 根据实际认证机制，从上下文中获取用户名
	// 如果使用了认证中间件，可以从ctx中获取
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		// 为了测试方便，可以返回一个默认用户名
		// return "test_user", nil
		return "", ErrUserNotLogin
	}
	return username, nil
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
