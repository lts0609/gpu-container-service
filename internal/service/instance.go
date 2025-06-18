package service

import (
	"context"
	"gpu-container-service/internal/model"
	"gpu-container-service/internal/util"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

// Creator 实现资源创建服务接口
type Creator struct {
	client clientset.Interface
}

// NewCreator 创建Pod服务实例
func NewCreator(client clientset.Interface) *Creator {
	return &Creator{
		client: client,
	}
}

// CreateInstance 创建关联资源
func (s *Creator) CreateInstance(ctx context.Context, req model.DeployCreateRequest) (*appsv1.Deployment, *v1.Secret, *v1.Service, error) {
	deploymentTemplate, err := util.GenerateDeploymentTemplate(req)
	if err != nil {
		klog.Errorf("Generate Deployment Template Error: %v", err)
		return nil, nil, nil, err
	}

	deployment, err := s.client.AppsV1().Deployments(deploymentTemplate.Namespace).Create(ctx, deploymentTemplate, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("Create Deployment Error: %v", err)
		return nil, nil, nil, err
	}

	secretTemplate, err := util.GenerateSecretTemplate(req, deployment)
	if err != nil {
		klog.Errorf("Generate Secret Template Error: %v", err)
		return nil, nil, nil, err
	}

	secret, err := s.client.CoreV1().Secrets(secretTemplate.Namespace).Create(ctx, secretTemplate, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("Create Secret Error: %v", err)
		return nil, nil, nil, err
	}

	serviceTemplate, err := util.GenerateServiceTemplate(req, deployment)
	if err != nil {
		klog.Errorf("Generate Service Template Error: %v", err)
		return nil, nil, nil, err
	}

	service, err := s.client.CoreV1().Services(serviceTemplate.Namespace).Create(ctx, serviceTemplate, metav1.CreateOptions{})
	if err != nil {
		klog.Errorf("Create Service Error: %v", err)
		return nil, nil, nil, err
	}

	return deployment, secret, service, nil
}
