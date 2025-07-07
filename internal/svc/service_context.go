package svc

import (
	"gpu-container-service/internal/config"
	"gpu-container-service/internal/types"
)

type GpuContainerServiceContext struct {
	Config    config.Config
	Instances map[string]types.Instance
}

func NewGpuContainerServiceContext(c config.Config) *GpuContainerServiceContext {
	return &GpuContainerServiceContext{
		Config:    c,
		Instances: make(map[string]types.Instance),
	}
}
