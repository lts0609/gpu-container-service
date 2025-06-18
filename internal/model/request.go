package model

// DeployCreateRequest 定义创建Pod的请求模型
type DeployCreateRequest struct {
	Name        string    `json:"name"`
	Namespace   string    `json:"namespace"`
	Image       string    `json:"image"`
	Resources   Resources `json:"resource"`
	Replicas    string    `json:"replicas"`
	Labels      string    `json:"labels"`
	Annotations string    `json:"annotations"`
}

// Resources 定义资源请求模型
type Resources struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
	GPU    string `json:"gpu"`
}

// Validate 验证请求参数
func (p *DeployCreateRequest) Validate() error {
	return nil
}
