package user

import (
	"context"
	"errors"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-sql-driver/mysql"
	"go-zero-shorterurl/admin/internal/constant"
	"go-zero-shorterurl/admin/internal/dal/model"
	"go-zero-shorterurl/admin/internal/types/errorx"
	"gorm.io/gorm"
	"time"

	"go-zero-shorterurl/admin/internal/svc"
	"go-zero-shorterurl/admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// internal/logic/user/user_register_logic.go
func (l *UserRegisterLogic) UserRegister(req *types.UserRegisterReq) (*types.UserRegisterResp, error) {
	// 1. 判断用户名是否存在
	exists, err := l.svcCtx.BloomFilters.UserExists(l.ctx, req.Username)
	if err != nil {
		return nil, errorx.NewSystemError(errorx.ServiceError) // 布隆过滤器操作失败
	}

	// 2.如果已经存在
	if exists {
		return nil, errorx.NewUserError(errorx.UserNameExistError)
	}

	// 3. 获取分布式锁
	lockKey := constant.Lock_User_Register + req.Username
	mutex := l.svcCtx.RS.NewMutex(lockKey,
		redsync.WithExpiry(time.Second*30), // 锁30秒过期
		redsync.WithTries(1),               // 只尝试一次
	)

	// 3.1 尝试获取锁
	if err = mutex.Lock(); err != nil {
		// 获取锁失败
		return nil, errorx.NewUserError(errorx.UserNameExistError)
	}
	defer func(mutex *redsync.Mutex) {
		_, err := mutex.Unlock()
		if err != nil {
			logx.Errorf("unlock failed: %v", err)
		}
	}(mutex)

	var createTime time.Time
	// 4. 开启事务
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		createTime = now // 保存创建时间用于返回

		// 4.1 创建用户
		user := &model.TUser{
			Username:   req.Username,
			Password:   req.Password, // 注意：实际使用时需要加密
			RealName:   req.RealName,
			Phone:      req.Phone,
			Mail:       req.Mail,
			CreateTime: now,
			UpdateTime: now,
		}
		// 4.2 尝试插入记录
		if err = tx.Create(user).Error; err != nil {
			if IsDuplicateError(err) {
				return errorx.NewUserError(errorx.UserNameExistError)
			}
			return errorx.NewSystemError(errorx.ServiceError)
		}

		// 4.3 创建默认分组
		group := &model.TGroup{
			Username:   req.Username,
			Name:       "默认分组",
			CreateTime: now,
			UpdateTime: now,
		}

		if err = tx.Create(group).Error; err != nil {
			return errorx.NewSystemError(errorx.ServiceError)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 5. 添加到布隆过滤器
	if err = l.svcCtx.BloomFilters.AddUser(l.ctx, req.Username); err != nil {
		// 这里的错误只记录日志，不影响主流程
		logx.Errorf("add to bloom filter failed: %v", err)
	}

	// 6. 返回成功响应
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
		return mysqlErr.Number == 1062
	}
	return false
}
