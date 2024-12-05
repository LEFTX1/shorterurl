package logic

import (
	"context"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"

	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/logx"
)

type RecycleBinPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRecycleBinPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinPageLogic {
	return &RecycleBinPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询回收站短链接
func (l *RecycleBinPageLogic) RecycleBinPage(in *__.RecycleBinPageRequest) (*__.RecycleBinPageResponse, error) {
	// 从metadata中获取用户名
	md, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}
	usernames := md.Get("username")
	if len(usernames) == 0 {
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}
	username := usernames[0]

	// 如果没有提供gid_list，则查询用户所有未删除的分组
	var gidList []string
	if len(in.GidList) == 0 {
		groups, err := l.svcCtx.Query.TGroup.WithContext(l.ctx).
			Where(l.svcCtx.Query.TGroup.Username.Eq(username)).
			Where(l.svcCtx.Query.TGroup.DelFlag.Is(false)).
			Find()
		if err != nil {
			return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "查询分组失败")
		}
		if len(groups) == 0 {
			return &__.RecycleBinPageResponse{
				Records: []*__.ShortLinkPageRecord{},
				Total:   0,
				Size:    int32(in.PageSize),
				Current: int32(in.PageNum),
			}, nil
		}
		for _, group := range groups {
			gidList = append(gidList, group.Gid)
		}
	} else {
		// 验证提供的gid是否属于当前用户
		groups, err := l.svcCtx.Query.TGroup.WithContext(l.ctx).
			Where(l.svcCtx.Query.TGroup.Username.Eq(username)).
			Where(l.svcCtx.Query.TGroup.Gid.In(in.GidList...)).
			Where(l.svcCtx.Query.TGroup.DelFlag.Is(false)).
			Find()
		if err != nil {
			return nil, errorx.New(errorx.SystemError, errorx.ErrInternalServer, "查询分组失败")
		}
		if len(groups) == 0 {
			return &__.RecycleBinPageResponse{
				Records: []*__.ShortLinkPageRecord{},
				Total:   0,
				Size:    int32(in.PageSize),
				Current: int32(in.PageNum),
			}, nil
		}
		for _, group := range groups {
			gidList = append(gidList, group.Gid)
		}
	}

	// TODO: 调用短链接服务的RPC接口查询回收站中的短链接
	// 这里需要实现调用短链接服务的逻辑
	// 示例响应：
	return &__.RecycleBinPageResponse{
		Records: []*__.ShortLinkPageRecord{},
		Total:   0,
		Size:    int32(in.PageSize),
		Current: int32(in.PageNum),
	}, nil
}
