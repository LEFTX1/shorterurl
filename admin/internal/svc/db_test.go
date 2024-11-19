package svc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-shorterurl/admin/internal/config"
	"go-zero-shorterurl/admin/internal/dal/model"
	"go-zero-shorterurl/pkg/snowflake"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
	"time"
)

// 测试配置
func getTestConfig() config.Config {
	return config.Config{
		DB: struct {
			Host     string
			Port     int
			User     string
			Password string
			Database string
		}{
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "root",
			Password: "123456",
			Database: "link",
		},
	}
}

// getTestDB 获取测试用的DB实例
func getTestDB() (*gorm.DB, error) {
	// 初始化雪花算法
	if err := snowflake.InitSnowflake(); err != nil {
		return nil, fmt.Errorf("init snowflake failed: %v", err)
	}

	// 获取ID生成器
	idGen, err := snowflake.GetSnowflakeGenerator()
	if err != nil {
		return nil, fmt.Errorf("get snowflake generator failed: %v", err)
	}

	// 初始化数据库
	return NewDB(getTestConfig(), idGen)
}

// 清理测试数据
func cleanTestData(db *gorm.DB) error {
	for i := 0; i < 16; i++ {
		tableName := fmt.Sprintf("t_user_%02d", i)
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE username LIKE 'test_%%'", tableName)).Error; err != nil {
			log.Printf("Warning: Failed to clean table %s: %v", tableName, err)
		}
	}
	return nil
}

// TestMain 用于测试的初始化和清理
func TestMain(m *testing.M) {
	db, err := getTestDB()
	if err != nil {
		log.Fatalf("Failed to init test db: %v", err)
	}

	// 清理旧的测试数据
	err = cleanTestData(db)
	if err != nil {
		logx.Errorf("Failed to clean test data: %v", err)
		return
	}

	// 运行测试
	code := m.Run()

	// 清理测试数据
	err = cleanTestData(db)
	if err != nil {
		logx.Errorf("Failed to clean test data: %v", err)
		return
	}

	os.Exit(code)
}

// TestUserSharding 测试用户分片
func TestUserSharding(t *testing.T) {
	db, err := getTestDB()
	assert.Nil(t, err)

	// 显示现有的表
	var tables []string
	err = db.Raw("SHOW TABLES LIKE 't_user_%'").Scan(&tables).Error
	assert.Nil(t, err)
	t.Logf("Existing tables: %v", tables)

	// 测试用户创建和分片
	testCases := []struct {
		username string
	}{
		{"test_user001"}, // 应该分到某个分片
		{"test_user002"}, // 应该分到另一个分片
		{"test_admin"},   // 测试特殊用户名
		{"test_query"},   // 再测试一个
	}

	// 创建用户并观察分片情况
	for _, tc := range testCases {
		t.Run(tc.username, func(t *testing.T) {
			user := &model.TUser{
				Username:   tc.username,
				Password:   "test123",
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}

			// 开启Debug模式查看SQL
			err := db.Debug().Create(user).Error
			assert.Nil(t, err)
			assert.NotZero(t, user.ID)

			// 验证是否可以查询到
			var result model.TUser
			err = db.Debug().Where("username = ?", tc.username).First(&result).Error
			assert.Nil(t, err)
			assert.Equal(t, tc.username, result.Username)

			t.Logf("User %s created with ID %d", tc.username, user.ID)
		})
	}

	// 验证分片查询
	t.Run("验证分片查询", func(t *testing.T) {
		var users []model.TUser
		// 这个查询应该能正确路由到各个分片
		err := db.Debug().Where("username LIKE ?", "test_%").Find(&users).Error
		assert.Nil(t, err)
		t.Logf("Found %d users", len(users))

		// 打印每个用户信息
		for _, user := range users {
			t.Logf("User: %s, ID: %d", user.Username, user.ID)
		}
	})
}
