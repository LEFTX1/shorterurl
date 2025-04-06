package group

import (
	"context"
	"io"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/api/internal/types/errorx"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type ListGroupsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListGroupsLogic {
	return &ListGroupsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListGroupsLogic) ListGroups() (resp []types.ShortLinkGroupResp, err error) {
	// 详细记录开始处理请求
	logx.Infof("[ListGroups] 开始处理分组列表请求")

	// 记录所有上下文键
	logx.Infof("[ListGroups] 上下文中的值:")
	if l.ctx == nil {
		logx.Errorf("[ListGroups] 上下文为空")
	} else {
		// 检查UserContextKey是否在上下文中
		userCtxValue := l.ctx.Value(types.UserContextKey)
		if userCtxValue == nil {
			logx.Errorf("[ListGroups] 上下文中未找到用户信息键: %v", types.UserContextKey)
		} else {
			logx.Infof("[ListGroups] 上下文中找到用户信息键: %v, 值类型: %T", types.UserContextKey, userCtxValue)
		}
	}

	// 从上下文中获取用户信息
	userCtxValue := l.ctx.Value(types.UserContextKey)
	if userCtxValue == nil {
		logx.Errorf("[ListGroups] 用户上下文信息为空")
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}

	// 尝试类型转换
	logx.Infof("[ListGroups] 尝试将上下文值转换为UserInfo类型，值类型: %T", userCtxValue)
	userInfo, ok := userCtxValue.(*types.UserInfo)
	if !ok {
		logx.Errorf("[ListGroups] 用户信息类型转换失败，期望类型: *types.UserInfo, 实际类型: %T", userCtxValue)
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "用户信息类型错误")
	}

	if userInfo == nil {
		logx.Errorf("[ListGroups] 转换成功但用户信息为nil")
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}

	// 记录用户信息详情
	logx.Infof("[ListGroups] 用户信息: ID=%s, Username=%s, RealName=%s",
		userInfo.ID, userInfo.Username, userInfo.RealName)

	// 使用推荐的方式添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)
	logx.Infof("[ListGroups] 添加用户名到RPC元数据: %s", userInfo.Username)

	// 调用RPC服务获取分组列表
	logx.Infof("[ListGroups] 开始调用RPC服务获取分组列表")
	stream, err := l.svcCtx.UserRpc.GroupList(ctx, &userservice.CommonRequest{})
	if err != nil {
		logx.Errorf("[ListGroups] 获取分组列表失败 username: %s, error: %v", userInfo.Username, err)
		return nil, err
	}

	// 处理流式响应
	var groups []types.ShortLinkGroupResp
	for {
		group, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				logx.Infof("[ListGroups] 分组数据接收完成，共%d个分组", len(groups))
				break
			}
			logx.Errorf("[ListGroups] 接收分组数据失败 username: %s, error: %v", userInfo.Username, err)
			return nil, err
		}

		logx.Infof("[ListGroups] 接收到分组: ID=%s, Name=%s, SortOrder=%d, Count=%d",
			group.Gid, group.Name, group.SortOrder, group.ShortLinkCount)

		groups = append(groups, types.ShortLinkGroupResp{
			Gid:            group.Gid,
			Name:           group.Name,
			SortOrder:      int(group.SortOrder),
			ShortLinkCount: int(group.ShortLinkCount),
		})
	}

	logx.Infof("[ListGroups] 成功获取分组列表，返回%d个分组", len(groups))
	return groups, nil
}
