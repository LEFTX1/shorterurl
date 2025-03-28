package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ShortLinkListGroupCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortLinkListGroupCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkListGroupCountLogic {
	return &ShortLinkListGroupCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ShortLinkListGroupCount 查询短链接分组内数量
func (l *ShortLinkListGroupCountLogic) ShortLinkListGroupCount(in *pb.GroupShortLinkCountRequest) (*pb.GroupShortLinkCountResponse, error) {
	// 参数校验
	if len(in.Gids) == 0 {
		return nil, status.Error(codes.InvalidArgument, "分组标识列表不能为空")
	}

	l.Logger.Infof("处理查询短链接分组内数量请求，分组标识列表: %v", in.Gids)

	// 构建查询
	// 相当于Java代码中的：
	// QueryWrapper<ShortLinkDO> queryWrapper = Wrappers.query(new ShortLinkDO())
	//     .select("gid as gid, count(*) as shortLinkCount")
	//     .in("gid", requestParam)
	//     .eq("enable_status", 0)
	//     .eq("del_flag", 0)
	//     .eq("del_time", 0L)
	//     .groupBy("gid");

	// 结果集
	result := &pb.GroupShortLinkCountResponse{
		GroupCounts: make([]*pb.ShortLinkGroupCountItem, 0, len(in.Gids)),
	}

	// 使用已有的Repo方法查询分组下的短链接数量
	// 在实际项目中，可以考虑在repo层添加批量查询方法，此处为简化实现
	for _, gid := range in.Gids {
		count, err := l.svcCtx.RepoManager.Link.CountByGid(l.ctx, gid)
		if err != nil {
			l.Logger.Errorf("查询分组 %s 的短链接数量失败: %v", gid, err)
			continue // 继续处理其他分组，不因一个分组查询失败而中断整个请求
		}

		result.GroupCounts = append(result.GroupCounts, &pb.ShortLinkGroupCountItem{
			Gid:            gid,
			ShortLinkCount: count,
		})

		l.Logger.Infof("分组 %s 的短链接数量: %d", gid, count)
	}

	l.Logger.Infof("查询短链接分组内数量成功，共 %d 个分组", len(result.GroupCounts))
	return result, nil
}
