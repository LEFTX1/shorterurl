package user

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-shorterurl/admin/internal/constant"
	"go-zero-shorterurl/admin/internal/dal/model"
	"go-zero-shorterurl/admin/internal/types/errorx"
	"gorm.io/gorm"
	"time"

	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// UserRegisterLogic 处理用户注册逻辑的结构体
type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserRegisterLogic 创建一个新的 UserRegisterLogic 实例
func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserRegister 处理用户注册请求
func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq) (*types.UserRegisterResp, error) {
	// 1. 检查布隆过滤器中是否存在该用户名
	exists, err := l.svcCtx.BloomFilters.UserExists(l.ctx, req.Username)
	if err != nil {
		// 如果检查布隆过滤器时出错，返回系统错误
		return nil, errorx.NewSystemError(errorx.ServiceError)
	}
	if exists {
		// 如果用户名已存在，返回用户错误
		return nil, errorx.NewUserError(errorx.UserNameExistError)
	}

	// 2. 创建分布式锁，防止并发注册相同用户名
	lockKey := constant.Lock_User_Register + req.Username
	lock := redis.NewRedisLock(l.svcCtx.Redis, lockKey)
	lock.SetExpire(30) // 设置锁的过期时间为30秒

	// 3. 获取锁
	acquired, err := lock.AcquireCtx(l.ctx)
	if err != nil {
		// 如果获取锁时出错，返回系统错误
		return nil, errorx.NewSystemError(errorx.ServiceError)
	}
	if !acquired {
		// 如果未能获取锁，返回用户错误
		return nil, errorx.NewUserError(errorx.UserNameExistError)
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
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		createTime = now

		// 4.1 创建用户
		user := &model.TUser{
			Username:   req.Username,
			Password:   req.Password,
			RealName:   req.RealName,
			Phone:      req.Phone,
			Mail:       req.Mail,
			CreateTime: now,
			UpdateTime: now,
		}

		// 4.2 尝试将用户信息插入数据库
		if err = tx.Create(user).Error; err != nil {
			if IsDuplicateError(err) {
				// 如果是重复键错误，返回用户错误
				return errorx.NewUserError(errorx.UserNameExistError)
			}
			// 其他错误返回系统错误
			return errorx.NewSystemError(errorx.ServiceError)
		}

		// 4.3 创建默认分组
		group := &model.TGroup{
			Username:   req.Username,
			Name:       "默认分组",
			CreateTime: now,
			UpdateTime: now,
		}

		// 4.4 尝试将分组信息插入数据库
		if err = tx.Create(group).Error; err != nil {
			// 如果插入分组信息时出错，返回系统错误
			return errorx.NewSystemError(errorx.ServiceError)
		}

		return nil
	})

	if err != nil {
		// 如果事务处理时出错，返回相应的错误
		return nil, err
	}

	// 5. 将用户名添加到布隆过滤器
	if err = l.svcCtx.BloomFilters.AddUser(l.ctx, req.Username); err != nil {
		// 如果添加到布隆过滤器时出错，记录错误日志
		logx.Errorf("添加布隆过滤器失败: %v", err)
	}

	// 6. 返回注册成功的响应
	return &types.UserRegisterResp{
		Username:   req.Username,
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
