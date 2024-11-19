package main

import (
	"fmt"
	"go-zero-shorterurl/admin/internal/dal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	// 连接数据库
	dsn := "root:123456@tcp(localhost:3306)/link?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("cannot connect db: %w", err))
	}

	// 初始化生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:           "../../internal/dal/query", // 生成代码的输出目录
		ModelPkgPath:      "../../internal/dal/model", // 生成模型的输出目录
		Mode:              gen.WithoutContext,         // 不使用context模式
		FieldNullable:     true,                       // 字段可为空
		FieldCoverable:    true,                       // 字段可覆盖
		FieldSignable:     true,                       // 字段可添加符号
		FieldWithIndexTag: true,                       // 生成字段索引tag
		FieldWithTypeTag:  true,                       // 生成字段类型tag
	})

	g.UseDB(db)

	// 生成基本模型结构体
	models := []interface{}{
		model.User{},  // 使用导入的模型
		model.Group{}, // 使用导入的模型
	}

	// 生成模型代码
	g.ApplyBasic(models...)

	// 执行生成
	g.Execute()
}
