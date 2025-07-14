package logic

import (
	"context"
	"gpu-container-service/internal/svc"
	"gpu-container-service/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	DefaultNamespace          = "default"
	DefaultGpuTypeLabelKey    = "mfy.com/gpu-type"
	DefaultAvaliableNodeLabel = "mfy.com/gpu-container"
)

type CreateInstanceLogic struct {
	Logger logx.Logger
	ctx    context.Context
	svcCtx *svc.GpuContainerServiceContext
}

func NewCreateInstanceLogic(ctx context.Context, svcCtx *svc.GpuContainerServiceContext) *CreateInstanceLogic {
	return &CreateInstanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateInstanceLogic) Service(ctx context.Context, req *types.CreateInstanceRequest, client clientset.Interface) (resp *types.Response, err error) {
	if err := req.Validate(req); err != nil {
		return &types.Response{
			Code:    int(types.ERR_INVALID_PARAMETER),
			Message: "request parameter error",
		}, err
	}
	logx.Info("request params is valid")

	var (
		namespace             string
		gpuTypeLabelKey       string
		avaliableNodeLabelKey string
	)

	instance := types.Instance{
		Uuid:            req.Uuid,
		User:            req.User,
		Status:          types.StatusPreparing,
		Image:           req.Image,
		CreateTime:      time.Now().String(),
		ChargeType:      req.ChargeType,
		Links:           types.Links{},
		ResourceRequest: req.ResourceRequest,
	}

	if l.svcCtx.Config.SchedulingPolicy == "share" {
		namespace = "default"
	} else {
		namespace = DefaultNamespace
	}

	if l.svcCtx.Config.GpuTypeLabel != "" {
		gpuTypeLabelKey = l.svcCtx.Config.GpuTypeLabel
	} else {
		gpuTypeLabelKey = DefaultGpuTypeLabelKey
	}

	if l.svcCtx.Config.AvaliableNodeLabel != "" {
		avaliableNodeLabelKey = l.svcCtx.Config.AvaliableNodeLabel
	} else {
		avaliableNodeLabelKey = DefaultAvaliableNodeLabel
	}

	// Create Pod
	podTemplate, err := GeneratePodTemplate(req, namespace, gpuTypeLabelKey, avaliableNodeLabelKey)
	if err != nil {
		return &types.Response{
			Code:    int(types.ERR_GENERATE_RESOURCE),
			Message: "generate pod template failed",
		}, err
	}
	logx.Info("generate pod template success")
	pod, err := client.CoreV1().Pods(namespace).Create(ctx, podTemplate, metav1.CreateOptions{})
	if err != nil {
		return &types.Response{
			Code:    int(types.ERR_RESOURCE_CREATE_FAILED),
			Message: "create pod failed",
		}, err
	}
	logx.Info("create pod %s success", pod.Name)
	// Create Secret
	secretTemplate, err := GenerateSecretTemplate(req, pod, namespace)
	if err != nil {
		return &types.Response{
			Code:    int(types.ERR_GENERATE_RESOURCE),
			Message: "generate secret template failed",
		}, err
	}
	logx.Info("generate secret template success")
	secret, err := client.CoreV1().Secrets(namespace).Create(ctx, secretTemplate, metav1.CreateOptions{})
	if err != nil {
		return &types.Response{
			Code:    int(types.ERR_RESOURCE_CREATE_FAILED),
			Message: "create secret failed",
		}, err
	}
	logx.Infof("create secret %s success", secret.Name)
	// Create Service
	serviceTemplate, err := GenerateServiceTemplate(req, pod, namespace)
	if err != nil {
		return &types.Response{
			Code:    int(types.ERR_GENERATE_RESOURCE),
			Message: "generate service template failed",
		}, err
	}
	service, err := client.CoreV1().Services(namespace).Create(ctx, serviceTemplate, metav1.CreateOptions{})
	if err != nil {
		return &types.Response{
			Code:    int(types.ERR_RESOURCE_CREATE_FAILED),
			Message: "create service failed",
		}, err
	}
	logx.Infof("create service %s success", service.Name)

	instance.Name = pod.Name
	l.svcCtx.Instances[req.Uuid] = instance

	return &types.Response{
		Code:    http.StatusOK,
		Message: "instance create success",
	}, err
}
