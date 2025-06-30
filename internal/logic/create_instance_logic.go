package logic

import (
	"context"
	"gpu-container-service/internal/svc"
	"gpu-container-service/internal/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateInstanceLogic struct {
	logx.Logger
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

func (l *CreateInstanceLogic) Service(ctx context.Context, req *types.CreateInstanceRequest, client clientset.Interface) (resp *types.CreateInstanceResponse, err error) {
	if err := req.Validate(req); err != nil {
		return &types.CreateInstanceResponse{
			Code:    int(types.ERR_INVALID_PARAMETER),
			Message: "request parameter error",
		}, err
	}
	logx.Info("request params is valid")

	namespace := "default"
	// Create Deployment
	deploymentTemplate, err := GenerateDeploymentTemplate(req, namespace)
	if err != nil {
		return &types.CreateInstanceResponse{
			Code:    int(types.ERR_GENERATE_RESOURCE),
			Message: "generate deployment template failed",
		}, err
	}
	logx.Info("generate deployment template success")
	deployment, err := client.AppsV1().Deployments(namespace).Create(ctx, deploymentTemplate, metav1.CreateOptions{})
	if err != nil {
		return &types.CreateInstanceResponse{
			Code:    int(types.ERR_RESOURCE_CREATE_FAILED),
			Message: "create deployment failed",
		}, err
	}
	logx.Info("create deployment %s success", deployment.Name)
	// Create Secret
	secretTemplate, err := GenerateSecretTemplate(req, deployment, namespace)
	if err != nil {
		return &types.CreateInstanceResponse{
			Code:    int(types.ERR_GENERATE_RESOURCE),
			Message: "generate secret template failed",
		}, err
	}
	logx.Info("generate secret template success")
	secret, err := client.CoreV1().Secrets(namespace).Create(ctx, secretTemplate, metav1.CreateOptions{})
	if err != nil {
		return &types.CreateInstanceResponse{
			Code:    int(types.ERR_RESOURCE_CREATE_FAILED),
			Message: "create secret failed",
		}, err
	}
	logx.Infof("create secret %s success", secret.Name)
	// Create Service
	serviceTemplate, err := GenerateServiceTemplate(req, deployment, namespace)
	if err != nil {
		return &types.CreateInstanceResponse{
			Code:    int(types.ERR_GENERATE_RESOURCE),
			Message: "generate service template failed",
		}, err
	}
	service, err := client.CoreV1().Services(namespace).Create(ctx, serviceTemplate, metav1.CreateOptions{})
	if err != nil {
		return &types.CreateInstanceResponse{
			Code:    int(types.ERR_RESOURCE_CREATE_FAILED),
			Message: "create service failed",
		}, err
	}
	logx.Infof("create service %s success", service.Name)

	return &types.CreateInstanceResponse{
		Code:    http.StatusOK,
		Message: "instance create success",
	}, err
}
