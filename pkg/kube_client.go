package pkg

import (
	"errors"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client interface {
	Config() (*restclient.Config, error)
	Client() (clientset.Interface, error)
}

type ClientImpl struct {
	kubeconfig string
	config     *restclient.Config
}

func NewClientBuilder(kubeconfig string) (*ClientImpl, error) {
	var config *restclient.Config
	var err error

	c := &ClientImpl{
		kubeconfig: kubeconfig,
	}

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", "")
	}

	if err != nil {
		return nil, errors.New("Error building config: " + err.Error())
	}
	c.config = config

	return c, nil
}

// Get the rest config
func (c ClientImpl) Config() (*restclient.Config, error) {
	config := c.config
	return restclient.AddUserAgent(config, "pod-creator"), nil
}

// Get the root client
func (c ClientImpl) Client() (clientset.Interface, error) {
	config, err := c.Config()
	if err != nil {
		return nil, err
	}
	return clientset.NewForConfig(config)
}
