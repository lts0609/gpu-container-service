package logic

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"gpu-container-service/internal/svc"
	"gpu-container-service/internal/types"
	"net/http"
)

type GetInstanceLogic struct {
	Logger logx.Logger
	ctx    context.Context
	svcCtx *svc.GpuContainerServiceContext
}

type GetAllInstanceLogic struct {
	Logger logx.Logger
	ctx    context.Context
	svcCtx *svc.GpuContainerServiceContext
}

func NewGetAllInstanceLogic(ctx context.Context, svcCtx *svc.GpuContainerServiceContext) *GetAllInstanceLogic {
	return &GetAllInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func NewGetInstanceLogic(ctx context.Context, svcCtx *svc.GpuContainerServiceContext) *GetInstanceLogic {
	return &GetInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllInstanceLogic) Service(ctx context.Context) (resp *types.Response, err error) {
	total := len(l.svcCtx.Instances)
	data := types.GetAllInstanceResponseData{
		ClusterName: l.svcCtx.Config.ClusterName,
		Total:       total,
		Items:       make([]types.Instance, total),
	}
	if total == 0 {
		return &types.Response{
			Code:    http.StatusOK,
			Message: "get instance success",
			Data:    data,
		}, nil
	}

	for _, item := range l.svcCtx.Instances {
		data.Items = append(data.Items, item)
	}

	return &types.Response{
		Code:    http.StatusOK,
		Message: "get instance success",
		Data:    data,
	}, nil
}

func (l *GetInstanceLogic) Service(ctx context.Context, req *types.GetInstanceRequest) (resp *types.Response, err error) {
	uuid := req.Uuid
	item, ok := l.svcCtx.Instances[uuid]
	if !ok {
		return &types.Response{
			Code:    int(types.ERR_RESOURCE_NOT_FOUND),
			Message: "instance not found",
		}, nil
	}

	data := types.GetInstanceResponseData{
		ClusterName: "cluster",
		Item:        item,
	}

	return &types.Response{
		Code:    http.StatusOK,
		Message: "get instance success",
		Data:    data,
	}, nil
}
