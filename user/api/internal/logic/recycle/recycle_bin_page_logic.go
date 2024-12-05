package recycle

import (
	"context"
	"shorterurl/user/api/internal/svc"
	"shorterurl/user/api/internal/types"
	"shorterurl/user/api/internal/types/errorx"
	"shorterurl/user/rpc/userservice"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type RecycleBinPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRecycleBinPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RecycleBinPageLogic {
	return &RecycleBinPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RecycleBinPageLogic) RecycleBinPage(req *types.RecycleBinPageReq) (resp *types.RecycleBinPageResp, err error) {
	// 从上下文中获取用户信息
	userInfo := l.ctx.Value(types.UserContextKey).(*types.UserInfo)
	if userInfo == nil {
		return nil, errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}

	// 使用推荐的方式添加元数据
	ctx := metadata.AppendToOutgoingContext(l.ctx, "username", userInfo.Username)

	// 调用RPC服务查询回收站分页
	rpcResp, err := l.svcCtx.UserRpc.RecycleBinPage(ctx, &userservice.RecycleBinPageRequest{
		GidList:  req.GidList,
		PageNum:  int32(req.Current),
		PageSize: int32(req.Size),
	})
	if err != nil {
		logx.Errorf("查询回收站分页失败 username: %s, current: %d, size: %d, error: %v",
			userInfo.Username, req.Current, req.Size, err)
		return nil, err
	}

	// 转换响应
	var records []types.ShortLinkPageRecordDTO
	for _, record := range rpcResp.Records {
		records = append(records, types.ShortLinkPageRecordDTO{
			Id:            record.Id,
			Domain:        record.Domain,
			ShortUri:      record.ShortUri,
			FullShortUrl:  record.FullShortUrl,
			OriginUrl:     record.OriginUrl,
			Gid:           record.Gid,
			ValidDateType: int(record.ValidDateType),
			ValidDate:     record.ValidDate,
			CreateTime:    record.CreateTime,
			Describe:      record.Describe,
			Favicon:       record.Favicon,
			EnableStatus:  int(record.EnableStatus),
			TotalPv:       record.TotalPv,
			TodayPv:       record.TodayPv,
			TotalUv:       record.TotalUv,
			TodayUv:       record.TodayUv,
			TotalUip:      record.TotalUip,
			TodayUip:      record.TodayUip,
		})
	}

	return &types.RecycleBinPageResp{
		Records: records,
		Total:   rpcResp.Total,
		Size:    int(rpcResp.Size),
		Current: int(rpcResp.Current),
	}, nil
}
