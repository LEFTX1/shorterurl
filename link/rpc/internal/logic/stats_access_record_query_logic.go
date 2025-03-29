package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatsAccessRecordQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatsAccessRecordQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsAccessRecordQueryLogic {
	return &StatsAccessRecordQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// StatsAccessRecordQuery 查询单个短链接指定时间内访问记录
func (l *StatsAccessRecordQueryLogic) StatsAccessRecordQuery(in *pb.AccessRecordQueryRequest) (*pb.AccessRecordQueryResponse, error) {
	// 参数验证
	if in.FullShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "短链接不能为空")
	}
	if in.Gid == "" {
		return nil, status.Error(codes.InvalidArgument, "分组标识不能为空")
	}
	if in.StartDate == "" || in.EndDate == "" {
		return nil, status.Error(codes.InvalidArgument, "开始日期和结束日期不能为空")
	}

	// 验证分组是否属于当前用户
	if err := l.checkGroupBelongToUser(in.Gid); err != nil {
		return nil, err
	}

	// 查询访问记录
	accessLogs, total, err := l.svcCtx.RepoManager.LinkAccessLogs.PageLinkAccessLogs(
		l.ctx,
		in.FullShortUrl,
		in.Gid,
		in.StartDate,
		in.EndDate,
		in.EnableStatus,
		in.Current,
		in.Size,
	)
	if err != nil {
		l.Logger.Errorf("查询短链接访问记录失败: %v", err)
		return nil, status.Error(codes.Internal, "查询短链接访问记录失败")
	}

	// 如果没有记录，返回空结果
	if len(accessLogs) == 0 {
		return &pb.AccessRecordQueryResponse{
			Records: []*pb.AccessRecord{},
			Total:   0,
			Size:    in.Size,
			Current: in.Current,
		}, nil
	}

	// 提取用户标识列表
	userList := make([]string, 0, len(accessLogs))
	for _, log := range accessLogs {
		userList = append(userList, log.User)
	}

	// 获取用户的访客类型
	uvTypes, err := l.svcCtx.RepoManager.LinkAccessLogs.SelectUvTypeByUsers(
		l.ctx,
		in.Gid,
		in.FullShortUrl,
		in.EnableStatus,
		in.StartDate,
		in.EndDate,
		userList,
	)
	if err != nil {
		l.Logger.Errorf("获取用户访客类型失败: %v", err)
		return nil, status.Error(codes.Internal, "获取用户访客类型失败")
	}

	// 创建用户访客类型映射
	uvTypeMap := make(map[string]string)
	for _, uvType := range uvTypes {
		user := uvType["user"].(string)
		uvTypeValue := uvType["uvType"].(string)
		uvTypeMap[user] = uvTypeValue
	}

	// 将数据转换为响应格式
	records := make([]*pb.AccessRecord, 0, len(accessLogs))
	for _, log := range accessLogs {
		// 获取访客类型，默认为"旧访客"
		uvType, ok := uvTypeMap[log.User]
		if !ok {
			uvType = "旧访客"
		}

		// 转换时间格式
		record := &pb.AccessRecord{
			UvType:     uvType,
			Browser:    log.Browser,
			Os:         log.Os,
			Ip:         log.Ip,
			Network:    log.Network,
			Device:     log.Device,
			Locale:     log.Locale,
			User:       log.User,
			CreateTime: log.CreateTime.Format("2006-01-02 15:04:05"),
		}
		records = append(records, record)
	}

	// 返回结果
	return &pb.AccessRecordQueryResponse{
		Records: records,
		Total:   total,
		Size:    in.Size,
		Current: in.Current,
	}, nil
}

// checkGroupBelongToUser 检查分组是否属于当前用户
func (l *StatsAccessRecordQueryLogic) checkGroupBelongToUser(gid string) error {
	// 获取当前登录用户
	username, err := l.svcCtx.RepoManager.GetCurrentUsername(l.ctx)
	if err != nil {
		return status.Error(codes.Unauthenticated, "用户未登录")
	}

	// 检查分组是否属于该用户
	exist, err := l.svcCtx.RepoManager.Group.CheckGroupBelongToUser(l.ctx, gid, username)
	if err != nil {
		l.Logger.Errorf("检查分组归属失败: %v", err)
		return status.Error(codes.Internal, "检查分组归属失败")
	}

	if !exist {
		return status.Error(codes.PermissionDenied, "用户信息与分组标识不匹配")
	}

	return nil
}
