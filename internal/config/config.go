package config

import (
	"gpu-container-service/pkg/util"
	clientset "k8s.io/client-go/kubernetes"
)

// Config 定义配置项结构体
type Config struct {
	ServicePort    string
	KubeconfigPath string
}

// NewConfig 创建配置实例
func NewConfig(servicePort, kubeconfigPath string) *Config {
	return &Config{
		ServicePort:    servicePort,
		KubeconfigPath: kubeconfigPath,
	}
}

// NewKubernetesClient 创建Kubernetes客户端
func NewKubernetesClient(kubeconfigPath string) (clientset.Interface, error) {
	return util.NewClientFactory(kubeconfigPath).Client()
}
