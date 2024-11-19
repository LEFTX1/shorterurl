package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/shorterurl?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{

		OutPath: "../../internal/dal/query",

		// 生成模式
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
	})

	g.UseDB(db)

	// 生成基础类型安全的 DAO API
	g.ApplyBasic(g.GenerateAllTable()...)

	// 执行生成
	g.Execute()
}
