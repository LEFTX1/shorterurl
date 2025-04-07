package logic

import (
	"context"
	"fmt"
	"time"

	"shorterurl/link/rpc/internal/model"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

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

// RecycleBinPage 分页查询回收站短链接
func (l *RecycleBinPageLogic) RecycleBinPage(in *pb.PageRecycleBinShortLinkRequest) (*pb.PageRecycleBinShortLinkResponse, error) {
	// 参数校验
	if in.Current <= 0 {
		in.Current = 1
	}
	if in.Size <= 0 {
		in.Size = 10
	}

	// 检查是否提供了gid参数
	if in.Gid == "" {
		l.Logger.Error("查询回收站短链接失败: 必须提供分组标识(gid)")
		return &pb.PageRecycleBinShortLinkResponse{
			Records: []*pb.ShortLinkRecord{},
			Total:   0,
			Size:    in.Size,
			Current: in.Current,
		}, nil
	}

	l.Logger.Infof("分页查询回收站短链接, 分组: %s, 页码: %d, 每页数量: %d", in.Gid, in.Current, in.Size)

	// 查询特定分组下的回收站短链接
	linksList, err := l.queryRecycleBinLinksByGid(in.Gid, int(in.Current), int(in.Size))
	if err != nil {
		l.Logger.Errorf("查询回收站短链接失败, 分组: %s, 错误: %v", in.Gid, err)
		return &pb.PageRecycleBinShortLinkResponse{
			Records: []*pb.ShortLinkRecord{},
			Total:   0,
			Size:    in.Size,
			Current: in.Current,
		}, nil
	}

	// 获取该分组下回收站短链接的总数
	count, err := l.countRecycleBinLinksByGid(in.Gid)
	if err != nil {
		l.Logger.Errorf("查询回收站短链接数量失败, 分组: %s, 错误: %v", in.Gid, err)
		count = 0
	}

	// 转换查询结果为ShortLinkRecord列表
	var links []*pb.ShortLinkRecord
	for _, link := range linksList {
		record := &pb.ShortLinkRecord{
			FullShortUrl: link.FullShortUrl,
			OriginUrl:    link.OriginUrl,
			Domain:       "http://" + link.Domain,
			Gid:          link.Gid,
			CreateTime:   link.CreateTime.Format(time.RFC3339),
			Describe:     link.Describe,
			TotalPv:      int32(link.TotalPv),
			TotalUv:      int32(link.TotalUv),
			TotalUip:     int32(link.TotalUip),
		}

		// 设置有效期
		if link.ValidDate.Unix() > 0 {
			record.ValidDate = link.ValidDate.Format(time.RFC3339)
		}

		links = append(links, record)
	}

	return &pb.PageRecycleBinShortLinkResponse{
		Records: links,
		Total:   int32(count),
		Size:    in.Size,
		Current: in.Current,
	}, nil
}

// queryRecycleBinLinksByGid 查询特定分组下的回收站短链接
func (l *RecycleBinPageLogic) queryRecycleBinLinksByGid(gid string, page, pageSize int) ([]*model.Link, error) {
	links, _, err := l.svcCtx.RepoManager.Link.FindRecycleBin(l.ctx, gid, page, pageSize)
	if err != nil {
		l.Logger.Errorf("查询回收站短链接失败: %v", err)
		return nil, fmt.Errorf("查询回收站短链接失败: %w", err)
	}
	return links, nil
}

// countRecycleBinLinksByGid 统计特定分组下回收站短链接的数量
func (l *RecycleBinPageLogic) countRecycleBinLinksByGid(gid string) (int64, error) {
	_, count, err := l.svcCtx.RepoManager.Link.FindRecycleBin(l.ctx, gid, 1, 0)
	if err != nil {
		l.Logger.Errorf("统计回收站短链接数量失败: %v", err)
		return 0, fmt.Errorf("统计回收站短链接数量失败: %w", err)
	}
	return count, nil
}
