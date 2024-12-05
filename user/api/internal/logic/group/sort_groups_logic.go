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

type SortGroupsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSortGroupsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SortGroupsLogic {
	return &SortGroupsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SortGroupsLogic) SortGroups(req *types.ShortLinkGroupSortReq) (resp *types.SuccessResp, err error) {
	// 从上下文中获取用户信息
	userInfo := l.ctx.Value(types.UserContextKey).(*types.UserInfo)
	if userInfo == nil {
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}

	// 使用推荐的方式添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)

	// 调用RPC服务排序分组
	stream, err := l.svcCtx.UserRpc.GroupSort(ctx)
	if err != nil {
		logx.Errorf("创建排序流失败 username: %s, error: %v", userInfo.Username, err)
		return nil, err
	}

	// 发送排序请求
	for _, sort := range req.Groups {
		err = stream.Send(&userservice.GroupSortRequest{
			Gid:       sort.Gid,
			SortOrder: int32(sort.SortOrder),
		})
		if err != nil {
			logx.Errorf("发送排序数据失败 username: %s, gid: %s, error: %v", userInfo.Username, sort.Gid, err)
			return nil, err
		}
	}

	// 关闭发送并接收响应
	_, err = stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		logx.Errorf("关闭排序流失败 username: %s, error: %v", userInfo.Username, err)
		return nil, err
	}

	return &types.SuccessResp{
		Code:    "0",
		Success: true,
	}, nil
}
