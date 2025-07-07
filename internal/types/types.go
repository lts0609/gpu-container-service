package types

import (
	"time"
)

type Status string

const (
	StatusPreparing   Status = "Preparing"
	StatusCreating    Status = "Creating"
	StatusRunning     Status = "Running"     // 实例已暂停
	StatusTerminating Status = "Terminating" // 实例终止中
	StatusTerminated  Status = "Terminated"  // 实例已终止
	StatusFailed      Status = "Failed"      // 实例创建或运行失败
)

type (
	GPUResource struct {
		Type string `json:"type"`
		Num  string `json:"num"`
	}

	StorageResource struct {
		SystemDisk string        `json:"system_disk"`
		DataVolume []PVCResource `json:"data_volume,omitempty"`
	}

	PVCResource struct {
		PVC  string `json:"pvc"`
		Size string `json:"size"`
	}

	Resources struct {
		CPU     string          `json:"cpu"`
		Memory  string          `json:"memory"`
		Storage StorageResource `json:"storage"`
		GPU     GPUResource     `json:"gpu"`
	}
)

type CreateInstanceRequest struct {
	Uuid            string    `json:"uuid"`
	User            string    `json:"user,omitempty"`
	Name            string    `json:"name"`
	Image           string    `json:"image"`
	ChargeType      string    `json:"charge_type"`
	ResourceRequest Resources `json:"resource_request"`
	Labels          string    `json:"labels,omitempty"`
}

type Instance struct {
	Uuid            string
	Name            string
	User            string
	Status          Status
	Image           string
	CreateTime      time.Time
	StartTime       time.Time
	DeleteTime      time.Time
	ChargeType      string
	Links           map[string]string
	ResourceRequest Resources
}

type GetInstanceRequest struct {
	Uuid string `json:"uuid"`
}

type GetInstanceResponseData struct {
	ClusterName string `json:"cluster_name"`
	Item        Instance
}

type GetAllInstanceResponseData struct {
	ClusterName string     `json:"cluster_name"`
	Total       int        `json:"total"`
	Items       []Instance `json:"items"`
}

type DeleteInstanceRequest struct {
	Uuid string `json:"uuid"`
}

type DeleteInstanceResponseData struct{}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type TerminalSessionRequest struct {
	Namespace     string `json:"namespace"`
	PodName       string `json:"pod_name"`
	ContainerName string `json:"container_name"`
	Shell         string `json:"shell"`
}

type TerminalSessionResponseData struct {
	SessionId string `json:"session_id"`
}

func (c *CreateInstanceRequest) Validate(req *CreateInstanceRequest) error {
	return nil
}
