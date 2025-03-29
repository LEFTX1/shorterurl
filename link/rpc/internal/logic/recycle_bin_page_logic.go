package logic

import (
	"context"
	"fmt"
	"time"

	"shorterurl/link/rpc/internal/model"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if len(in.Gids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "分组标识列表不能为空")
	}
	if in.Current <= 0 {
		in.Current = 1
	}
	if in.Size <= 0 {
		in.Size = 10
	}

	l.Logger.Infof("分页查询回收站短链接, 分组: %v, 页码: %d, 每页数量: %d", in.Gids, in.Current, in.Size)

	// 查询回收站中的短链接列表
	var links []*pb.ShortLinkRecord
	var total int64 = 0

	// 考虑到这里需要查询多个分组的数据，我们需要进行多次查询并合并结果
	for _, gid := range in.Gids {
		// 查询回收站中的短链接
		// 注意: 这里不使用FindRecycleBin，因为该方法是通过DelFlag=1来查询，
		// 而在回收站中的链接是EnableStatus=1但DelFlag=0
		linksList, err := l.queryRecycleBinLinksByGid(gid, int(in.Current), int(in.Size))
		if err != nil {
			l.Logger.Errorf("查询回收站短链接失败, 分组: %s, 错误: %v", gid, err)
			continue
		}

		// 转换查询结果为ShortLinkRecord列表
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

		// 获取该分组下回收站短链接的总数
		count, err := l.countRecycleBinLinksByGid(gid)
		if err != nil {
			l.Logger.Errorf("查询回收站短链接数量失败, 分组: %s, 错误: %v", gid, err)
			continue
		}
		total += count
	}

	return &pb.PageRecycleBinShortLinkResponse{
		Records: links,
		Total:   int32(total),
		Size:    in.Size,
		Current: in.Current,
	}, nil
}

// 查询指定分组下回收站中的短链接
func (l *RecycleBinPageLogic) queryRecycleBinLinksByGid(gid string, page, pageSize int) ([]*model.Link, error) {
	// 使用 FindByGidWithCondition 代替 FindByCondition
	links, err := l.svcCtx.RepoManager.Link.FindByGidWithCondition(l.ctx, gid, map[string]interface{}{
		"enable_status": 1, // 未启用状态表示在回收站中
		"del_flag":      0, // 未删除
	}, page, pageSize)

	if err != nil {
		return nil, fmt.Errorf("查询回收站短链接失败: %v", err)
	}

	return links, nil
}

// 统计指定分组下回收站中的短链接数量
func (l *RecycleBinPageLogic) countRecycleBinLinksByGid(gid string) (int64, error) {
	// 使用 CountByGidWithCondition 代替 CountByCondition
	count, err := l.svcCtx.RepoManager.Link.CountByGidWithCondition(l.ctx, gid, map[string]interface{}{
		"enable_status": 1, // 未启用状态表示在回收站中
		"del_flag":      0, // 未删除
	})

	if err != nil {
		return 0, fmt.Errorf("统计回收站短链接数量失败: %v", err)
	}

	return count, nil
}
