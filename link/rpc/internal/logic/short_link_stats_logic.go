package logic

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"shorterurl/link/rpc/internal/consumer"
	"shorterurl/link/rpc/internal/svc"
	"shorterurl/link/rpc/pb"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
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

// 短链接统计 - 接收外部统计请求
func (l *ShortLinkStatsLogic) ShortLinkStats(in *pb.ShortLinkStatsRequest) (*pb.EmptyResponse, error) {
	// 记录请求信息
	l.Logger.Infof("接收到短链接统计请求: %s, IP: %s", in.FullShortUrl, in.Ip)

	// 构建用户标识 - 如果请求中没有用户标识，根据IP和浏览器信息生成
	user := in.User
	if user == "" {
		// 生成临时用户标识
		h := md5.New()
		io.WriteString(h, in.Ip)
		io.WriteString(h, in.Browser)
		io.WriteString(h, in.Os)
		user = fmt.Sprintf("%x", h.Sum(nil))
	}

	// 检查是否是新的 UV 和 UIP
	uvFirstFlag := l.checkFirstUv(in.FullShortUrl, user)
	uipFirstFlag := l.checkFirstUip(in.FullShortUrl, in.Ip)

	// 提交统计记录到消费者队列
	statsRecord := &consumer.StatsRecord{
		FullShortUrl: in.FullShortUrl,
		Gid:          in.Gid,
		User:         user,
		UvFirstFlag:  uvFirstFlag,
		UipFirstFlag: uipFirstFlag,
		Ip:           in.Ip,
		Browser:      in.Browser,
		Os:           in.Os,
		Device:       in.Device,
		Network:      in.Network,
		Locale:       in.Locale,
		CurrentDate:  time.Now(),
	}

	// 如果没有地理位置信息但有IP，尝试获取地理位置
	if statsRecord.Locale == "" && statsRecord.Ip != "" {
		// 创建IP位置查询逻辑
		ipLocationLogic := NewGetIpLocationLogic(l.ctx, l.svcCtx)
		formattedLocation, err := ipLocationLogic.GetFormattedLocation(statsRecord.Ip)

		if err == nil && formattedLocation != "" {
			l.Logger.Infof("IP地理位置解析成功: %s -> %s", statsRecord.Ip, formattedLocation)
			statsRecord.Locale = formattedLocation
		}
	}

	// 提交到统计消费者队列
	l.svcCtx.StatsConsumer.Submit(statsRecord)

	return &pb.EmptyResponse{}, nil
}

// 检查是否是新的 UV
func (l *ShortLinkStatsLogic) checkFirstUv(fullShortUrl, user string) bool {
	key := fmt.Sprintf("short-link:stats:uv:%s", fullShortUrl)
	added, err := l.svcCtx.BizRedis.Sadd(key, user)
	if err != nil {
		l.Logger.Errorf("检查UV失败: %v", err)
		return false
	}
	// 设置过期时间 (90天)
	l.svcCtx.BizRedis.Expire(key, 90*24*60*60)
	return added > 0
}

// 检查是否是新的 UIP
func (l *ShortLinkStatsLogic) checkFirstUip(fullShortUrl, ip string) bool {
	if ip == "" {
		return false
	}
	key := fmt.Sprintf("short-link:stats:uip:%s", fullShortUrl)
	added, err := l.svcCtx.BizRedis.Sadd(key, ip)
	if err != nil {
		l.Logger.Errorf("检查UIP失败: %v", err)
		return false
	}
	// 设置过期时间 (90天)
	l.svcCtx.BizRedis.Expire(key, 90*24*60*60)
	return added > 0
}
