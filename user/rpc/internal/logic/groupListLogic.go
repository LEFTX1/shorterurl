package logic

import (
	"context"
	"shorterurl/user/rpc/internal/svc"
	"shorterurl/user/rpc/internal/types/errorx"
	__ "shorterurl/user/rpc/pb"

	"google.golang.org/grpc/metadata"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupListLogic {
	return &GroupListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupList 获取分组列表
func (l *GroupListLogic) GroupList(in *__.CommonRequest, stream __.UserService_GroupListServer) error {
	// 从metadata中获取用户名
	md, ok := metadata.FromIncomingContext(l.ctx)
	if !ok {
		return errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}
	usernames := md.Get("username")
	if len(usernames) == 0 {
		return errorx.New(errorx.ClientError, errorx.ErrInternalServer, "未找到用户信息")
	}
	username := usernames[0]

	// 添加详细日志
	logx.Infof("[GroupList] 开始查询用户 '%s' 的分组列表", username)

	// 查询用户的所有分组，按排序号和更新时间降序排列
	groups, err := l.svcCtx.Query.TGroup.WithContext(l.ctx).
		Where(l.svcCtx.Query.TGroup.Username.Eq(username)).
		Where(l.svcCtx.Query.TGroup.DelFlag.Is(false)).
		Order(l.svcCtx.Query.TGroup.SortOrder.Desc(), l.svcCtx.Query.TGroup.UpdateTime.Desc()).
		Find()

	// 记录SQL查询结果
	if err != nil {
		logx.Errorf("[GroupList] 查询分组失败: %v", err)
		return errorx.New(errorx.SystemError, errorx.ErrInternalServer, "获取分组列表失败")
	}

	logx.Infof("[GroupList] 查询到 %d 个分组", len(groups))
	for i, group := range groups {
		logx.Infof("[GroupList] 分组 #%d: ID=%s, 名称=%s, 用户名=%s, GID=%s, 排序=%d, 创建时间=%v",
			i+1, group.ID, group.Name, group.Username, group.Gid, group.SortOrder, group.CreateTime)
	}

	// 收集所有分组的 GID
	var gids []string
	for _, group := range groups {
		gids = append(gids, group.Gid)
	}

	// 检查是否找到分组
	if len(groups) == 0 {
		logx.Errorf("[GroupList] 用户 '%s' 没有任何分组，这可能是一个问题，因为正常注册后应该有默认分组", username)
		// 尝试查询是否有被标记为删除的分组
		deletedGroups, err := l.svcCtx.Query.TGroup.WithContext(l.ctx).
			Where(l.svcCtx.Query.TGroup.Username.Eq(username)).
			Where(l.svcCtx.Query.TGroup.DelFlag.Is(true)).
			Find()
		if err != nil {
			logx.Errorf("[GroupList] 查询已删除分组失败: %v", err)
		} else {
			logx.Infof("[GroupList] 找到 %d 个已删除的分组", len(deletedGroups))
			for i, group := range deletedGroups {
				logx.Infof("[GroupList] 已删除分组 #%d: ID=%s, 名称=%s, GID=%s",
					i+1, group.ID, group.Name, group.Gid)
			}
		}
	}

	// 发送分组数据到客户端
	for _, group := range groups {
		// 构建响应对象
		resp := &__.GroupResponse{
			Gid:            group.Gid,
			Name:           group.Name,
			SortOrder:      int32(group.SortOrder),
			ShortLinkCount: 0, // 默认值，未来会从链接服务获取
		}

		// 发送响应
		if err := stream.Send(resp); err != nil {
			logx.Errorf("[GroupList] 发送分组数据失败: %v", err)
			return err
		}
	}

	logx.Infof("[GroupList] 成功发送所有分组数据到客户端")
	return nil
}
