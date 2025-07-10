package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	ClusterName        string `json:"cluster_name"`
	SchedulingPolicy   string `json:"scheduling_policy"`
	GpuTypeLabel       string `json:"gpu_type_label"`
	AvaliableNodeLabel string `json:"avaliable_node_label"`
}
