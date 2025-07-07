package logic

import (
	"context"
	"gpu-container-service/internal/svc"
	"gpu-container-service/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type TerminalSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.GpuContainerServiceContext
}

func NewTerminalSessionLogic(ctx context.Context, svcCtx *svc.GpuContainerServiceContext) *TerminalSessionLogic {
	return &TerminalSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TerminalSessionLogic) Service(ctx context.Context, req *types.TerminalSessionRequest) (resp *types.Response, err error) {
	return &types.Response{
		Code:    http.StatusOK,
		Message: "SessionId generate success",
		Data: types.TerminalSessionResponseData{
			SessionId: "",
		},
	}, err
}
