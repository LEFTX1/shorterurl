package logic

import (
	"context"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StatsGroupAccessRecordQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewStatsGroupAccessRecordQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsGroupAccessRecordQueryLogic {
	return &StatsGroupAccessRecordQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// StatsGroupAccessRecordQuery 查询分组短链接指定时间内访问记录
func (l *StatsGroupAccessRecordQueryLogic) StatsGroupAccessRecordQuery(in *pb.GroupAccessRecordQueryRequest) (*pb.GroupAccessRecordQueryResponse, error) {
	// 参数验证
	if in.Gid == "" {
		return nil, status.Error(codes.InvalidArgument, "分组标识不能为空")
	}
	if in.StartDate == "" || in.EndDate == "" {
		return nil, status.Error(codes.InvalidArgument, "开始日期和结束日期不能为空")
	}
	// 校验分页参数
	if in.Current <= 0 {
		return nil, status.Error(codes.InvalidArgument, "当前页码必须大于0")
	}
	if in.Size <= 0 {
		return nil, status.Error(codes.InvalidArgument, "每页数量必须大于0")
	}

	// 验证分组是否属于当前用户
	if err := l.checkGroupBelongToUser(in.Gid); err != nil {
		return nil, err
	}

	// 查询访问记录
	accessLogs, total, err := l.svcCtx.RepoManager.LinkAccessLogs.PageGroupAccessLogs(
		l.ctx,
		in.Gid,
		in.StartDate,
		in.EndDate,
		in.Current,
		in.Size,
	)
	if err != nil {
		l.Logger.Errorf("查询分组访问记录失败: %v", err)
		return nil, status.Error(codes.Internal, "查询分组访问记录失败")
	}

	// 如果没有记录，返回空结果
	if len(accessLogs) == 0 {
		return &pb.GroupAccessRecordQueryResponse{
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
	uvTypes, err := l.svcCtx.RepoManager.LinkAccessLogs.SelectGroupUvTypeByUsers(
		l.ctx,
		in.Gid,
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
	return &pb.GroupAccessRecordQueryResponse{
		Records: records,
		Total:   total,
		Size:    in.Size,
		Current: in.Current,
	}, nil
}

// checkGroupBelongToUser 检查分组是否属于当前用户
func (l *StatsGroupAccessRecordQueryLogic) checkGroupBelongToUser(gid string) error {
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
