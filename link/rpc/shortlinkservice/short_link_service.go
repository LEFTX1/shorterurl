// Code generated by goctl. DO NOT EDIT.
// Source: link.proto

package shortlinkservice

import (
	"context"

	"shorterurl/link/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	AccessRecord                    = pb.AccessRecord
	AccessRecordQueryRequest        = pb.AccessRecordQueryRequest
	AccessRecordQueryResponse       = pb.AccessRecordQueryResponse
	BatchCreateResult               = pb.BatchCreateResult
	BatchCreateShortLinkRequest     = pb.BatchCreateShortLinkRequest
	BatchCreateShortLinkResponse    = pb.BatchCreateShortLinkResponse
	BrowserStat                     = pb.BrowserStat
	CreateShortLinkRequest          = pb.CreateShortLinkRequest
	CreateShortLinkResponse         = pb.CreateShortLinkResponse
	DailyStat                       = pb.DailyStat
	DeviceStat                      = pb.DeviceStat
	GetGroupStatsRequest            = pb.GetGroupStatsRequest
	GetGroupStatsResponse           = pb.GetGroupStatsResponse
	GetShortLinkCountRequest        = pb.GetShortLinkCountRequest
	GetShortLinkCountResponse       = pb.GetShortLinkCountResponse
	GetSingleStatsRequest           = pb.GetSingleStatsRequest
	GetSingleStatsResponse          = pb.GetSingleStatsResponse
	GetUrlTitleRequest              = pb.GetUrlTitleRequest
	GetUrlTitleResponse             = pb.GetUrlTitleResponse
	GroupAccessRecordQueryRequest   = pb.GroupAccessRecordQueryRequest
	GroupAccessRecordQueryResponse  = pb.GroupAccessRecordQueryResponse
	GroupCount                      = pb.GroupCount
	LocaleCnStat                    = pb.LocaleCnStat
	NetworkStat                     = pb.NetworkStat
	OSStat                          = pb.OSStat
	PageRecycleBinShortLinkRequest  = pb.PageRecycleBinShortLinkRequest
	PageRecycleBinShortLinkResponse = pb.PageRecycleBinShortLinkResponse
	PageShortLinkRequest            = pb.PageShortLinkRequest
	PageShortLinkResponse           = pb.PageShortLinkResponse
	RecoverFromRecycleBinRequest    = pb.RecoverFromRecycleBinRequest
	RecoverFromRecycleBinResponse   = pb.RecoverFromRecycleBinResponse
	RemoveFromRecycleBinRequest     = pb.RemoveFromRecycleBinRequest
	RemoveFromRecycleBinResponse    = pb.RemoveFromRecycleBinResponse
	SaveToRecycleBinRequest         = pb.SaveToRecycleBinRequest
	SaveToRecycleBinResponse        = pb.SaveToRecycleBinResponse
	ShortLinkRecord                 = pb.ShortLinkRecord
	UpdateShortLinkRequest          = pb.UpdateShortLinkRequest
	UpdateShortLinkResponse         = pb.UpdateShortLinkResponse

	ShortLinkService interface {
		// --------------------- 短链接管理接口 ---------------------
		ShortLinkCreate(ctx context.Context, in *CreateShortLinkRequest, opts ...grpc.CallOption) (*CreateShortLinkResponse, error)
		ShortLinkBatchCreate(ctx context.Context, in *BatchCreateShortLinkRequest, opts ...grpc.CallOption) (*BatchCreateShortLinkResponse, error)
		ShortLinkUpdate(ctx context.Context, in *UpdateShortLinkRequest, opts ...grpc.CallOption) (*UpdateShortLinkResponse, error)
		ShortLinkPage(ctx context.Context, in *PageShortLinkRequest, opts ...grpc.CallOption) (*PageShortLinkResponse, error)
		// --------------------- 回收站管理接口 ---------------------
		RecycleBinSave(ctx context.Context, in *SaveToRecycleBinRequest, opts ...grpc.CallOption) (*SaveToRecycleBinResponse, error)
		RecycleBinRecover(ctx context.Context, in *RecoverFromRecycleBinRequest, opts ...grpc.CallOption) (*RecoverFromRecycleBinResponse, error)
		RecycleBinRemove(ctx context.Context, in *RemoveFromRecycleBinRequest, opts ...grpc.CallOption) (*RemoveFromRecycleBinResponse, error)
		RecycleBinPage(ctx context.Context, in *PageRecycleBinShortLinkRequest, opts ...grpc.CallOption) (*PageRecycleBinShortLinkResponse, error)
		// --------------------- 短链接统计接口 ---------------------
		StatsGetSingle(ctx context.Context, in *GetSingleStatsRequest, opts ...grpc.CallOption) (*GetSingleStatsResponse, error)
		StatsGetGroup(ctx context.Context, in *GetGroupStatsRequest, opts ...grpc.CallOption) (*GetGroupStatsResponse, error)
		StatsGetShortLinkCount(ctx context.Context, in *GetShortLinkCountRequest, opts ...grpc.CallOption) (*GetShortLinkCountResponse, error)
		StatsAccessRecordQuery(ctx context.Context, in *AccessRecordQueryRequest, opts ...grpc.CallOption) (*AccessRecordQueryResponse, error)
		StatsGroupAccessRecordQuery(ctx context.Context, in *GroupAccessRecordQueryRequest, opts ...grpc.CallOption) (*GroupAccessRecordQueryResponse, error)
		// --------------------- URL标题功能接口 ---------------------
		UrlTitleGet(ctx context.Context, in *GetUrlTitleRequest, opts ...grpc.CallOption) (*GetUrlTitleResponse, error)
	}

	defaultShortLinkService struct {
		cli zrpc.Client
	}
)

func NewShortLinkService(cli zrpc.Client) ShortLinkService {
	return &defaultShortLinkService{
		cli: cli,
	}
}

// --------------------- 短链接管理接口 ---------------------
func (m *defaultShortLinkService) ShortLinkCreate(ctx context.Context, in *CreateShortLinkRequest, opts ...grpc.CallOption) (*CreateShortLinkResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.ShortLinkCreate(ctx, in, opts...)
}

func (m *defaultShortLinkService) ShortLinkBatchCreate(ctx context.Context, in *BatchCreateShortLinkRequest, opts ...grpc.CallOption) (*BatchCreateShortLinkResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.ShortLinkBatchCreate(ctx, in, opts...)
}

func (m *defaultShortLinkService) ShortLinkUpdate(ctx context.Context, in *UpdateShortLinkRequest, opts ...grpc.CallOption) (*UpdateShortLinkResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.ShortLinkUpdate(ctx, in, opts...)
}

func (m *defaultShortLinkService) ShortLinkPage(ctx context.Context, in *PageShortLinkRequest, opts ...grpc.CallOption) (*PageShortLinkResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.ShortLinkPage(ctx, in, opts...)
}

// --------------------- 回收站管理接口 ---------------------
func (m *defaultShortLinkService) RecycleBinSave(ctx context.Context, in *SaveToRecycleBinRequest, opts ...grpc.CallOption) (*SaveToRecycleBinResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.RecycleBinSave(ctx, in, opts...)
}

func (m *defaultShortLinkService) RecycleBinRecover(ctx context.Context, in *RecoverFromRecycleBinRequest, opts ...grpc.CallOption) (*RecoverFromRecycleBinResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.RecycleBinRecover(ctx, in, opts...)
}

func (m *defaultShortLinkService) RecycleBinRemove(ctx context.Context, in *RemoveFromRecycleBinRequest, opts ...grpc.CallOption) (*RemoveFromRecycleBinResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.RecycleBinRemove(ctx, in, opts...)
}

func (m *defaultShortLinkService) RecycleBinPage(ctx context.Context, in *PageRecycleBinShortLinkRequest, opts ...grpc.CallOption) (*PageRecycleBinShortLinkResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.RecycleBinPage(ctx, in, opts...)
}

// --------------------- 短链接统计接口 ---------------------
func (m *defaultShortLinkService) StatsGetSingle(ctx context.Context, in *GetSingleStatsRequest, opts ...grpc.CallOption) (*GetSingleStatsResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.StatsGetSingle(ctx, in, opts...)
}

func (m *defaultShortLinkService) StatsGetGroup(ctx context.Context, in *GetGroupStatsRequest, opts ...grpc.CallOption) (*GetGroupStatsResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.StatsGetGroup(ctx, in, opts...)
}

func (m *defaultShortLinkService) StatsGetShortLinkCount(ctx context.Context, in *GetShortLinkCountRequest, opts ...grpc.CallOption) (*GetShortLinkCountResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.StatsGetShortLinkCount(ctx, in, opts...)
}

func (m *defaultShortLinkService) StatsAccessRecordQuery(ctx context.Context, in *AccessRecordQueryRequest, opts ...grpc.CallOption) (*AccessRecordQueryResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.StatsAccessRecordQuery(ctx, in, opts...)
}

func (m *defaultShortLinkService) StatsGroupAccessRecordQuery(ctx context.Context, in *GroupAccessRecordQueryRequest, opts ...grpc.CallOption) (*GroupAccessRecordQueryResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.StatsGroupAccessRecordQuery(ctx, in, opts...)
}

// --------------------- URL标题功能接口 ---------------------
func (m *defaultShortLinkService) UrlTitleGet(ctx context.Context, in *GetUrlTitleRequest, opts ...grpc.CallOption) (*GetUrlTitleResponse, error) {
	client := pb.NewShortLinkServiceClient(m.cli.Conn())
	return client.UrlTitleGet(ctx, in, opts...)
}