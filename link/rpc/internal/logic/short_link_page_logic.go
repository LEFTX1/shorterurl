package logic

import (
	"context"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShortLinkPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortLinkPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkPageLogic {
	return &ShortLinkPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShortLinkPageLogic) ShortLinkPage(in *pb.PageShortLinkRequest) (*pb.PageShortLinkResponse, error) {
	// 参数校验
	if in.Gid == "" {
		return nil, status.Error(codes.InvalidArgument, "分组标识不能为空")
	}

	// 设置默认分页参数
	page := int(in.Current)
	if page <= 0 {
		page = 1
	}

	pageSize := int(in.Size)
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	// 查询数据
	links, total, err := l.svcCtx.RepoManager.Link.FindByGid(l.ctx, in.Gid, page, pageSize)
	if err != nil {
		l.Logger.Errorf("查询短链接列表失败: %v", err)
		return nil, status.Error(codes.Internal, "查询短链接列表失败")
	}

	// 构建响应
	records := make([]*pb.ShortLinkRecord, 0, len(links))
	for _, link := range links {
		record := &pb.ShortLinkRecord{
			FullShortUrl: link.FullShortUrl,
			OriginUrl:    link.OriginUrl,
			Domain:       link.Domain,
			Gid:          link.Gid,
			CreateTime:   link.CreateTime.Format(time.RFC3339),
			ValidDate:    link.ValidDate.Format(time.RFC3339),
			Describe:     link.Describe,
			TotalPv:      int32(link.TotalPv),
			TotalUv:      int32(link.TotalUv),
			TotalUip:     int32(link.TotalUip),
		}
		records = append(records, record)
	}

	return &pb.PageShortLinkResponse{
		Total:   int32(total),
		Size:    int32(pageSize),
		Current: int32(page),
		Records: records,
	}, nil
}
