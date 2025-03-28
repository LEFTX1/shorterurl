package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/threading"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Redis统计队列键
const (
	ShortLinkStatsQueueKey = "queue:stats:shortlink"
)

type ShortLinkStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShortLinkStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShortLinkStatsLogic {
	return &ShortLinkStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ShortLinkStats 短链接统计
func (l *ShortLinkStatsLogic) ShortLinkStats(in *pb.ShortLinkStatsRequest) (*pb.EmptyResponse, error) {
	// 参数校验
	if in.FullShortUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "完整短链接不能为空")
	}

	if in.Gid == "" {
		return nil, status.Error(codes.InvalidArgument, "分组标识不能为空")
	}

	l.Logger.Infof("处理短链接统计请求，短链接: %s, 分组: %s", in.FullShortUrl, in.Gid)

	// 创建访问记录
	// 参考Java代码，这里我们使用消息队列的方式异步处理统计数据
	// 为了简化实现，我们将统计记录发送到Redis的List中

	// 构建完整统计记录，包括不存储在数据库中的字段
	type StatsRecordDTO struct {
		FullShortUrl string    `json:"fullShortUrl"`
		Gid          string    `json:"gid"`
		User         string    `json:"user"`
		IP           string    `json:"ip"`
		Browser      string    `json:"browser"`
		Os           string    `json:"os"`
		Device       string    `json:"device"`
		Network      string    `json:"network"`
		Locale       string    `json:"locale"`
		UvType       string    `json:"uvType"`
		CreateTime   time.Time `json:"createTime"`
	}

	statsRecord := StatsRecordDTO{
		FullShortUrl: in.FullShortUrl,
		Gid:          in.Gid,
		User:         in.User,
		IP:           in.Ip,
		Browser:      in.Browser,
		Os:           in.Os,
		Device:       in.Device,
		Network:      in.Network,
		Locale:       in.Locale,
		UvType:       in.UvType,
		CreateTime:   time.Now(),
	}

	// 将记录序列化为JSON
	jsonData, err := json.Marshal(statsRecord)
	if err != nil {
		l.Logger.Errorf("序列化统计记录失败: %v", err)
		return nil, status.Error(codes.Internal, "序列化统计记录失败")
	}

	// 将JSON数据发送到Redis队列
	_, err = l.svcCtx.BizRedis.Lpush(ShortLinkStatsQueueKey, string(jsonData))
	if err != nil {
		l.Logger.Errorf("发送统计记录到队列失败: %v", err)
		return nil, status.Error(codes.Internal, "发送统计记录到队列失败")
	}

	// 异步更新链接访问量（PV）
	l.asyncUpdateLinkStats(in.FullShortUrl, in.Gid)

	l.Logger.Infof("短链接统计成功，记录已发送到队列")
	return &pb.EmptyResponse{}, nil
}

// asyncUpdateLinkStats 异步更新链接的访问统计数据
func (l *ShortLinkStatsLogic) asyncUpdateLinkStats(fullShortUrl, gid string) {
	threading.GoSafe(func() {
		ctx := context.Background()

		// 1. 查询链接信息
		link, err := l.svcCtx.RepoManager.Link.FindByFullShortUrlAndGid(ctx, fullShortUrl, gid)
		if err != nil {
			logx.Errorf("异步查询链接信息失败: %v", err)
			return
		}

		// 2. 更新PV（访问量）
		link.TotalPv = link.TotalPv + 1

		// 3. 更新链接信息
		err = l.svcCtx.RepoManager.Link.Update(ctx, link)
		if err != nil {
			logx.Errorf("异步更新链接统计数据失败: %v", err)
			return
		}

		// 4. 更新Redis缓存中的统计数据
		cacheKey := fmt.Sprintf("link:stats:%s", fullShortUrl)
		statsData := map[string]string{
			"pv":  fmt.Sprintf("%d", link.TotalPv),
			"uv":  fmt.Sprintf("%d", link.TotalUv),
			"uip": fmt.Sprintf("%d", link.TotalUip),
		}
		for k, v := range statsData {
			err = l.svcCtx.BizRedis.Hset(cacheKey, k, v)
			if err != nil {
				logx.Errorf("更新统计缓存失败 [%s=%s]: %v", k, v, err)
			}
		}

		logx.Infof("异步更新链接统计数据成功: %s", fullShortUrl)
	})
}
