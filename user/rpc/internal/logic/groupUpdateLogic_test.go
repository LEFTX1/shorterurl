package logic

import (
	"errors"
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

func TestGroupUpdate(t *testing.T) {
	svcCtx, ctx := setupTest(t)
	logic := NewGroupUpdateLogic(ctx, svcCtx)
	createLogic := NewGroupCreateLogic(ctx, svcCtx)
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

	// 创建测试分组
	createReq := &__.GroupSaveRequest{
		Username:  username,
		GroupName: "测试分组",
	}
	_, err = createLogic.GroupCreate(createReq)
	require.NoError(t, err, "创建分组失败")

	// 获取创建的分组ID
	group, err := svcCtx.Query.TGroup.WithContext(ctx).
		Where(svcCtx.Query.TGroup.Username.Eq(username)).
		Where(svcCtx.Query.TGroup.Name.Eq(createReq.GroupName)).
		First()
	require.NoError(t, err, "查询分组失败")
	require.NotNil(t, group, "分组不应为空")
	t.Logf("创建的分组: %+v", group.Gid)

	// 清理测试数据
	defer func() {
		q := query.Use(svcCtx.DB)
		_, _ = q.TGroup.WithContext(ctx).Where(q.TGroup.Username.Eq(username)).Delete()
		_, _ = svcCtx.Redis.DelCtx(ctx, constant.LockGroupUpdateKey+group.Gid)
	}()

	t.Run("成功更新分组", func(t *testing.T) {
		req := &__.GroupUpdateRequest{
			Gid:  group.Gid,
			Name: "更新后的分组名",
		}

		resp, err := logic.GroupUpdate(req)
		require.NoError(t, err, "更新分组应该成功")
		require.NotNil(t, resp, "响应不应为空")
		assert.True(t, resp.Success, "更新应该成功")
		assert.Equal(t, "更新成功", resp.Message)

		// 验证分组是否真实更新
		updatedGroup, err := svcCtx.Query.TGroup.WithContext(ctx).
			Where(svcCtx.Query.TGroup.Gid.Eq(group.Gid)).
			First()
		require.NoError(t, err, "查询分组失败")
		require.NotNil(t, updatedGroup, "分组应该存在")
		assert.Equal(t, req.Name, updatedGroup.Name)
		assert.True(t, updatedGroup.UpdateTime.After(group.UpdateTime))
	})

	t.Run("分组不存在", func(t *testing.T) {
		req := &__.GroupUpdateRequest{
			Gid:  "nonexistent_gid",
			Name: "新分组名",
		}

		resp, err := logic.GroupUpdate(req)
		require.Error(t, err, "应该返回错误")
		assert.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		var appErr *errorx.AppError
		ok := errors.As(err, &appErr)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type)
		assert.Equal(t, errorx.ErrGroupNotFound, appErr.Code)
	})

	t.Run("并发更新", func(t *testing.T) {
		// 清理之前的测试数据
		_, _ = svcCtx.Redis.DelCtx(ctx, constant.LockGroupUpdateKey+group.Gid)
		time.Sleep(time.Second)

		// 创建一个新的分布式锁
		lockKey := constant.LockGroupUpdateKey + group.Gid
		lock := redis.NewRedisLock(svcCtx.Redis, lockKey)
		lock.SetExpire(30)
		acquired, err := lock.AcquireCtx(ctx)
		require.NoError(t, err, "获取锁失败")
		require.True(t, acquired, "应该能获取到锁")

		// 在锁被占用的情况下尝试更新分组
		req := &__.GroupUpdateRequest{
			Gid:  group.Gid,
			Name: "并发更新测试",
		}
		resp, err := logic.GroupUpdate(req)
		require.Error(t, err, "应该返回错误")
		assert.Nil(t, resp, "响应应该为空")

		// 验证错误类型
		var appErr *errorx.AppError
		ok := errors.As(err, &appErr)
		assert.True(t, ok, "应该返回 AppError")
		assert.Equal(t, errorx.ClientError, appErr.Type)
		assert.Equal(t, errorx.ErrTooManyRequests, appErr.Code)

		// 释放锁
		released, err := lock.Release()
		require.NoError(t, err, "释放锁失败")
		require.True(t, released, "锁应该被释放")
	})
}
