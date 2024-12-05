package logic

import (
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/dal/query"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func TestGroupCreate(t *testing.T) {
	svcCtx, ctx := setupTest(t)
	logic := NewGroupCreateLogic(ctx, svcCtx)
	registerLogic := NewUserRegisterLogic(ctx, svcCtx)

	// 准备测试数据：用户名和注册请求
	username := generateTestUsername()
	registerReq := &__.RegisterRequest{
		Username: username,
		Password: "password123",
		RealName: "Test User",
		Phone:    "13800138000",
		Mail:     "test@example.com",
	}

	// 注册用户
	_, err := registerLogic.UserRegister(registerReq)
	require.NoError(t, err, "注册用户失败")

	// 清理测试数据
	defer func() {
		q := query.Use(svcCtx.DB)
		_, _ = q.TGroup.WithContext(ctx).Where(q.TGroup.Username.Eq(username)).Delete()
		_, _ = svcCtx.Redis.DelCtx(ctx, constant.LOCK_GROUP_CREATE_KEY+username)
	}()

	t.Run("成功创建分组", func(t *testing.T) {
		req := &__.GroupSaveRequest{
			Username:  username,
			GroupName: "测试分组",
		}

		resp, err := logic.GroupCreate(req)
		require.NoError(t, err, "创建分组应该成功")
		require.NotNil(t, resp, "响应不应为空")
		assert.True(t, resp.Success, "创建应该成功")
		assert.Equal(t, "创建成功", resp.Message)

		// 验证分组是否真实创建
		group, err := svcCtx.Query.TGroup.WithContext(ctx).
			Where(svcCtx.Query.TGroup.Username.Eq(username)).
			Where(svcCtx.Query.TGroup.Name.Eq(req.GroupName)).
			First()
		require.NoError(t, err, "查询分组失败")
		require.NotNil(t, group, "分组应该存在")
		assert.Equal(t, req.GroupName, group.Name)
		assert.Equal(t, int32(0), group.SortOrder)
		assert.Equal(t, username, group.Username)
		assert.Equal(t, false, group.DelFlag)
		assert.Len(t, group.Gid, 8, "Gid应该是8位字符串")
	})

	t.Run("创建第二个分组", func(t *testing.T) {
		req := &__.GroupSaveRequest{
			Username:  username,
			GroupName: "第二个分组",
		}

		resp, err := logic.GroupCreate(req)
		require.NoError(t, err, "创建分组应该成功")
		require.NotNil(t, resp, "响应不应为空")
		assert.True(t, resp.Success, "创建应该成功")

		// 验证分组是否真实创建
		group, err := svcCtx.Query.TGroup.WithContext(ctx).
			Where(svcCtx.Query.TGroup.Username.Eq(username)).
			Where(svcCtx.Query.TGroup.Name.Eq(req.GroupName)).
			First()
		require.NoError(t, err, "查询分组失败")
		require.NotNil(t, group, "分组应该存在")
		assert.Equal(t, req.GroupName, group.Name)
		assert.Equal(t, int32(0), group.SortOrder)
		assert.Equal(t, username, group.Username)
		assert.Equal(t, false, group.DelFlag)
		assert.Len(t, group.Gid, 8, "Gid应该是8位字符串")
	})

	t.Run("超出分组数量限制", func(t *testing.T) {
		// 清理之前的测试数据
		q := query.Use(svcCtx.DB)
		_, err := q.TGroup.WithContext(ctx).Where(q.TGroup.Username.Eq(username)).Delete()
		require.NoError(t, err, "清理测试数据失败")

		// 创建19个分组
		for i := 0; i < 19; i++ {
			req := &__.GroupSaveRequest{
				Username:  username,
				GroupName: generateTestUsername(), // 使用随机名称
			}
			resp, err := logic.GroupCreate(req)
			require.NoError(t, err, "创建分组应该成功")
			require.NotNil(t, resp, "响应不应为空")
		}

		// 尝试创建第20个分组，应该成功
		req := &__.GroupSaveRequest{
			Username:  username,
			GroupName: "第20个分组",
		}
		resp, err := logic.GroupCreate(req)
		require.NoError(t, err, "创建第20个分组应该成功")
		require.NotNil(t, resp, "响应不应为空")

		// 尝试创建第21个分组，应该失败
		req = &__.GroupSaveRequest{
			Username:  username,
			GroupName: "超限分组",
		}
		resp, err = logic.GroupCreate(req)
		require.Error(t, err, "应该返回错误")
		assert.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		appErr, ok := err.(*errorx.AppError)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type)
		assert.Equal(t, errorx.ErrGroupLimit, appErr.Code)
	})

	t.Run("并发创建分组", func(t *testing.T) {
		// 清理之前的测试数据
		q := query.Use(svcCtx.DB)
		_, err := q.TGroup.WithContext(ctx).Where(q.TGroup.Username.Eq(username)).Delete()
		require.NoError(t, err, "清理测试数据失败")

		// 确保Redis锁已释放
		time.Sleep(time.Second)

		// 创建一个新的分布式锁
		lockKey := constant.LOCK_GROUP_CREATE_KEY + username
		lock := redis.NewRedisLock(svcCtx.Redis, lockKey)
		lock.SetExpire(30)
		acquired, err := lock.AcquireCtx(ctx)
		require.NoError(t, err, "获取锁失败")
		require.True(t, acquired, "应该能获取到锁")

		// 在锁被占用的情况下尝试创建分组
		req := &__.GroupSaveRequest{
			Username:  username,
			GroupName: "并发测试分组",
		}
		resp, err := logic.GroupCreate(req)
		require.Error(t, err, "应该返回错误")
		assert.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		appErr, ok := err.(*errorx.AppError)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type)
		assert.Equal(t, errorx.ErrTooManyRequests, appErr.Code)

		// 释放锁
		released, err := lock.Release()
		require.NoError(t, err, "释放锁失败")
		require.True(t, released, "锁应该被释放")
	})
}
