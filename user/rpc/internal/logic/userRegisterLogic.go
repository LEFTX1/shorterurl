package logic

import (
	"context"
	"errors"
	"shorterurl/user/rpc/internal/constant"
	"shorterurl/user/rpc/internal/dal/model"
	"shorterurl/user/rpc/internal/dal/query"
	"shorterurl/user/rpc/internal/types/errorx"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/redis"

	"shorterurl/user/rpc/internal/svc"
	__ "shorterurl/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户注册
func (l *UserRegisterLogic) UserRegister(in *__.RegisterRequest) (*__.RegisterResponse, error) {
	// 1. 检查布隆过滤器中是否存在该用户名
	exists, err := l.svcCtx.BloomFilters.UserExists(l.ctx, in.Username)
	if err != nil {
		// 如果检查布隆过滤器时出错，返回系统错误
		return nil, errorx.New(errorx.SystemError, errorx.ErrBloomFilterCheck, errorx.Message(errorx.ErrBloomFilterCheck))
	}
	if exists {
		// 如果用户名已存在，返回用户错误
		return nil, errorx.New(errorx.ClientError, errorx.ErrUserNameExists, errorx.Message(errorx.ErrUserNameExists))
	}
	// 2. 创建分布式锁，防止并发注册相同用户名
	lockKey := constant.Lock_User_Register + in.Username
	lock := redis.NewRedisLock(l.svcCtx.Redis, lockKey)
	lock.SetExpire(30) // 设置锁的过期时间为30秒

	// 3. 获取锁
	acquired, err := lock.AcquireCtx(l.ctx)
	if err != nil {
		// 如果获取锁时出错，返回系统错误
		return nil, errorx.New(errorx.SystemError, errorx.ErrDistributedLock, errorx.Message(errorx.ErrDistributedLock))
	}
	if !acquired {
		// 如果未能获取锁，返回用户错误
		return nil, errorx.New(errorx.ClientError, errorx.ErrTooManyRequests, errorx.Message(errorx.ErrTooManyRequests))
	}

	// 确保在函数结束时释放锁
	defer func() {
		if released, err := lock.Release(); err != nil {
			// 如果释放锁时出错，记录错误日志
			logx.Errorf("释放锁失败: %v", err)
		} else if !released {
			// 如果锁未被主动释放，记录错误日志
			logx.Error("锁未被主动释放")
		}
	}()
	var createTime time.Time
	// 4. 事务处理，确保用户和默认分组的创建是原子的
	err = l.svcCtx.Query.Transaction(func(tx *query.Query) error {
		createTime = time.Now()
		// 4.1 创建用户
		user := &model.TUser{
			Username:   in.Username,
			Password:   in.Password,
			RealName:   in.RealName,
			Phone:      in.Phone,
			Mail:       in.Mail,
			CreateTime: createTime,
			UpdateTime: createTime,
		}

		// 4.2 尝试将用户信息插入数据库
		if err := tx.TUser.Create(user); err != nil {
			if IsDuplicateError(err) {
				// 如果是重复键错误，返回用户错误
				return errorx.New(errorx.ClientError, errorx.ErrUserNameExists, errorx.Message(errorx.ErrUserNameExists))
			}
			// 其他错误返回系统错误
			return errorx.New(errorx.SystemError, errorx.ErrDatabaseOperation, errorx.Message(errorx.ErrDatabaseOperation))
		}

		createTime = time.Now()

		// 4.3 创建默认分组
		group := &model.TGroup{
			Username:   in.Username,
			Name:       "默认分组",
			CreateTime: createTime,
			UpdateTime: createTime,
		}

		// 4.4 尝试将分组信息插入数据库
		if err := tx.TGroup.Create(group); err != nil {
			// 如果插入分组信息时出错，返回系统错误
			return errorx.New(errorx.SystemError, errorx.ErrDatabaseOperation, errorx.Message(errorx.ErrDatabaseOperation))
		}

		return nil
	})

	if err != nil {
		// 如果事务处理时出错，返回相应的错误
		logx.Errorf("事务处理失败: %v 时间: %v", err, createTime)
		return nil, err
	}

	// 5. 将用户名添加到布隆过滤器
	if err = l.svcCtx.BloomFilters.AddUser(l.ctx, in.Username); err != nil {
		// 如果添加到布隆过滤器时出错，记录错误日志
		logx.Errorf("添加布隆过滤器失败: %v", err)
	}

	return &__.RegisterResponse{
		Username:   in.Username,
		CreateTime: createTime.Format("2006-01-02 15:04:05"),
		Message:    "注册成功",
	}, nil
}

// IsDuplicateError 辅助函数：检查是否是重复键错误
func IsDuplicateError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		// MySQL错误码1062表示重复键错误
		return mysqlErr.Number == 1062
	}
	return false
}
