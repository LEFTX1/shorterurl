package testutil

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go-zero-shorterurl/admin/internal/dal/model"
	"testing"
	"time"
)

func TestUserSharding(t *testing.T) {
	testDB, err := GetTestDB()
	assert.Nil(t, err, "获取测试DB失败")

	// 1. 清理所有测试数据
	for i := 0; i < 16; i++ {
		tableName := fmt.Sprintf("t_user_%02d", i)
		err := testDB.Exec(fmt.Sprintf("DELETE FROM %s WHERE username LIKE 'test_%%'", tableName)).Error
		if err != nil {
			t.Logf("清理表 %s 失败: %v", tableName, err)
		}
	}

	// 2. 检查分片表
	var tables []string
	err = testDB.Raw("SHOW TABLES LIKE 't_user_%'").Scan(&tables).Error
	assert.Nil(t, err)
	t.Logf("现有分片表: %v", tables)

	// 3. 测试数据
	testUsers := []string{
		"test_user001",
		"test_user002",
		"test_admin",
	}

	for _, username := range testUsers {
		// 创建用户
		user := &model.TUser{
			Username:   username,
			Password:   "test123",
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}

		// 插入数据
		err = testDB.Debug().Create(user).Error
		assert.Nil(t, err, fmt.Sprintf("创建用户 %s 失败", username))
		t.Logf("插入用户: username=%s, id=%d", username, user.ID)

		// 查询验证
		var found model.TUser
		err = testDB.Debug().Where("username = ?", username).First(&found).Error
		assert.Nil(t, err, fmt.Sprintf("查询用户 %s 失败", username))
		t.Logf("查询用户: username=%s, id=%d", found.Username, found.ID)

		// 查找数据所在的表
		for _, table := range tables {
			var count int64
			err = testDB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE username = ?", table), username).Count(&count).Error
			assert.Nil(t, err)
			if count > 0 {
				t.Logf("用户 %s 在表 %s 中找到", username, table)
			}
		}
	}
}
