// internal/logic/user/user_crud_test.go
package user

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"go-zero-shorterurl/admin/internal/dal/model"
	"go-zero-shorterurl/admin/internal/dal/query"
	"io"
	"os"
	"strings"
	"testing"
)

func TestUserCRUD(t *testing.T) {
	// 捕获标准输出和标准错误
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rStdout, wStdout, err := os.Pipe()
	require.NoError(t, err)
	rStderr, wStderr, err := os.Pipe()
	require.NoError(t, err)

	// 替换标准输出和标准错误
	os.Stdout = wStdout
	os.Stderr = wStderr

	outputChan := make(chan string)

	go func() {
		var output strings.Builder
		mr := io.MultiReader(rStdout, rStderr)
		_, err := io.Copy(&output, mr)
		if err != nil {
			t.Logf("捕获输出时发生错误: %v", err)
		}
		outputChan <- output.String()
	}()

	defer func() {
		os.Stdout = oldStdout
		os.Stderr = oldStderr
		wStdout.Close()
		wStderr.Close()
		output := <-outputChan
		t.Logf("\n捕获的所有输出信息:\n%s", output)
	}()

	// 获取服务上下文
	svcCtx, ctx := setupTest(t)

	// GORM 测试组
	t.Run("GORM操作测试", func(t *testing.T) {
		// INSERT
		t.Run("插入操作", func(t *testing.T) {
			user := &model.TUser{
				Username: generateTestUsername(),
				DelFlag:  false,
			}

			result := svcCtx.DB.Create(user)
			require.NoError(t, result.Error)
			t.Logf("GORM单条插入成功，插入的用户信息: %+v", user)
		})

		// SELECT
		t.Run("查询操作", func(t *testing.T) {
			var users []model.TUser

			result := svcCtx.DB.Where("username = ? AND del_flag = ?", "test_user", false).Find(&users)
			require.NoError(t, result.Error)
			t.Logf("GORM查询成功，查询到的用户列表: %+v", users)
		})

		// UPDATE
		t.Run("更新操作", func(t *testing.T) {
			result := svcCtx.DB.Model(&model.TUser{}).
				Where("username = ?", "test_user").
				Update("del_flag", true)

			require.NoError(t, result.Error)
			t.Logf("GORM更新成功，影响的行数: %d", result.RowsAffected)
		})

		// DELETE
		t.Run("删除操作", func(t *testing.T) {
			result := svcCtx.DB.Where("username = ?", "test_user").Delete(&model.TUser{})
			require.NoError(t, result.Error)
			t.Logf("GORM删除成功，删除的行数: %d", result.RowsAffected)
		})
	})

	// Gen 测试组
	t.Run("Gen操作测试", func(t *testing.T) {
		// INSERT
		t.Run("插入操作", func(t *testing.T) {
			user := &model.TUser{
				Username: generateTestUsername(),
				DelFlag:  false,
			}

			err := svcCtx.Query.TUser.WithContext(ctx).Create(user)
			require.NoError(t, err)
			t.Logf("Gen单条插入成功，插入的用户信息: %+v", user)
		})

		// SELECT
		t.Run("查询操作", func(t *testing.T) {
			users, err := svcCtx.Query.TUser.WithContext(ctx).
				Where(svcCtx.Query.TUser.Username.Eq("test_user_gen")).
				Where(svcCtx.Query.TUser.DelFlag.Is(false)).
				Find()

			require.NoError(t, err)
			fmt.Printf("ending sql : %v\n", svcCtx.Sharding.LastQuery())
			t.Logf("Gen查询成功，查询到的用户列表: %+v", users)
		})

		// UPDATE
		t.Run("更新操作", func(t *testing.T) {
			_, err := svcCtx.Query.TUser.WithContext(ctx).
				Where(svcCtx.Query.TUser.Username.Eq("test_user_gen")).
				Update(svcCtx.Query.TUser.DelFlag, true)

			require.NoError(t, err)
			t.Logf("Gen更新操作执行成功")
		})

		// DELETE
		t.Run("删除操作", func(t *testing.T) {
			_, err := svcCtx.Query.TUser.WithContext(ctx).
				Where(svcCtx.Query.TUser.Username.Eq("test_user_gen")).
				Delete()

			require.NoError(t, err)
			t.Logf("Gen删除操作执行成功")
		})
	})

	// 事务测试
	t.Run("事务操作测试", func(t *testing.T) {
		// GORM 事务
		t.Run("GORM事务", func(t *testing.T) {
			tx := svcCtx.DB.Begin()
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			user := &model.TUser{Username: generateTestUsername(), DelFlag: false}
			if err := tx.Create(user).Error; err != nil {
				tx.Rollback()
				t.Fatalf("GORM事务中创建用户失败: %v", err)
			}

			if err := tx.Commit().Error; err != nil {
				t.Fatalf("GORM事务提交失败: %v", err)
			}

			t.Log("GORM事务执行成功完成")
		})

		// Gen 事务
		t.Run("Gen事务", func(t *testing.T) {
			err := svcCtx.Query.Transaction(func(tx *query.Query) error {
				user := &model.TUser{Username: generateTestUsername(), DelFlag: false}
				return tx.TUser.WithContext(ctx).Create(user)
			})
			require.NoError(t, err)
			t.Log("Gen事务执行成功完成")
		})
	})

}
