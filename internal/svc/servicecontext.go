package svc

import (
	"gpu-container-service/internal/config"
)

type CreateInstanceContext struct {
	Config config.Config
}

func NewCreateInstanceContext(c config.Config) *CreateInstanceContext {
	return &CreateInstanceContext{
		Config: c,
	}
}
