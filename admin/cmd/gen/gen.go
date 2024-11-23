// cmd/gen/gen.go
package main

import (
	"flag"
	"fmt"
	"go-zero-shorterurl/admin/internal/config"
	"time"

	"github.com/zeromicro/go-zero/core/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// 定义逻辑表结构
type User struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	Username     string    `gorm:"column:username;not null"`
	Password     string    `gorm:"column:password;not null"`
	RealName     string    `gorm:"column:real_name"`
	Phone        string    `gorm:"column:phone"`
	Mail         string    `gorm:"column:mail"`
	DeletionTime int64     `gorm:"column:deletion_time"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
	DelFlag      bool      `gorm:"column:del_flag"`
}

func (User) TableName() string {
	return "t_user"
}

type Group struct {
	ID         int64     `gorm:"column:id;primaryKey"`
	Username   string    `gorm:"column:username;not null"`
	Name       string    `gorm:"column:name;not null"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (Group) TableName() string {
	return "t_group"
}

func main() {
	configFile := flag.String("f", "../../etc/admin-api.yaml", "配置文件路径")
	flag.Parse()

	// 加载配置
	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 创建数据库连接
	db, err := setupDB(c)
	if err != nil {
		panic(err)
	}

	// 生成代码
	g := gen.NewGenerator(gen.Config{
		OutPath:      "../../internal/dal/query",
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		WithUnitTest: true,
	})

	g.UseDB(db)

	// all tables
	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}

func setupDB(c config.Config) (*gorm.DB, error) {
	// 1. 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/shorterurl?charset=utf8mb4&parseTime=True&loc=Local",
		c.DB.User,
		c.DB.Password,
		c.DB.Host,
		c.DB.Port,
	)
	// 2. 打开数据库连接
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database failed: %v", err)
	}

	// 3. 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB failed: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
