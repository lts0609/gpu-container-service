package util

import (
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientFactory 定义客户端工厂接口
type ClientFactory interface {
	Client() (clientset.Interface, error)
}

// clientFactoryImpl 实现客户端工厂接口
type clientFactoryImpl struct {
	kubeconfigPath string
}

// NewClientFactory 创建客户端工厂实例
func NewClientFactory(kubeconfigPath string) ClientFactory {
	return &clientFactoryImpl{
		kubeconfigPath: kubeconfigPath,
	}
}

// Client 获取Kubernetes客户端
func (c *clientFactoryImpl) Client() (clientset.Interface, error) {
	// only use in test
	// kubeconfig, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", c.kubeconfigPath)
	if err != nil {
		return nil, err
	}

	return clientset.NewForConfig(kubeconfig)
}
