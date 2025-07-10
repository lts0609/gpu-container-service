package logic

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gpu-container-service/internal/types"
	"gpu-container-service/pkg"
	"hash/fnv"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
)

const (
	SSHPort     int32 = 22
	JupyterPort int32 = 8888
)

var (
	DeafultDomain = "containercloud-xian.xaidc.com"
	PodBaseName   = "gpu-instance"
)

func GeneratePodTemplate(req *types.CreateInstanceRequest, namespace string, gpuTypeLabelKey string, avaliableNodeLabelKey string) (*v1.Pod, error) {
	logx.Info("Generate New Pod Template")

	labels, err := ParseLabels(req.Name, req.Labels)
	if err != nil {
		return nil, fmt.Errorf("Parse Label Error: %v", err)
	}
	// generate jupyter enviroment
	env := []v1.EnvVar{
		{
			Name: "POD_NAME",
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		{
			Name:  "NB_PREFIX",
			Value: "/notebook/$(POD_NAME)",
		},
	}
	// calc pod hash-code by request uuid
	hash := fnv.New32()
	hash.Write([]byte(req.Uuid))
	hashcode := fmt.Sprintf("%x", hash.Sum32())[:8]

	container := v1.Container{
		Name:  PodBaseName + "-" + hashcode,
		Image: req.Image,
		Ports: []v1.ContainerPort{
			{
				Name:          "ssh",
				ContainerPort: SSHPort,
			},
			{
				Name:          "jupyter",
				ContainerPort: JupyterPort,
			},
		},
		Resources: ParseResources(req.ResourceRequest),
		Env:       env,
		EnvFrom: []v1.EnvFromSource{
			{
				SecretRef: &v1.SecretEnvSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: PodBaseName + "-" + hashcode + "-secret",
					},
				},
			},
		},
	}

	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      PodBaseName + "-" + hashcode,
			Namespace: namespace,
			Labels:    labels,
			Annotations: map[string]string{
				"uuid": req.Uuid,
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{container},
			NodeSelector: map[string]string{
				avaliableNodeLabelKey: "true",
				gpuTypeLabelKey:       req.ResourceRequest.GPU.Type,
			},
		},
	}, nil
}

func GenerateSecretTemplate(req *types.CreateInstanceRequest, pod *v1.Pod, namespace string) (*v1.Secret, error) {
	password, hashedPassword, err := pkg.GenerateJupyterPassword()
	if err != nil {
		return nil, fmt.Errorf("GenerateJupyterPassword Error: %v", err)
	}

	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name + "-secret",
			Namespace: namespace,
			Annotations: map[string]string{
				"uuid": req.Uuid,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: v1.SchemeGroupVersion.String(),
					Kind:       "Pod",
					Name:       pod.Name,
					UID:        pod.UID,
					Controller: new(bool),
				},
			},
		},
		Type: v1.SecretTypeOpaque,
		Data: map[string][]byte{
			"SSH_PASSWORD":     password,
			"NB_PASSWD":        password,
			"NB_HASHED_PASSWD": hashedPassword,
		},
	}, nil
}

func GenerateServiceTemplate(req *types.CreateInstanceRequest, pod *v1.Pod, namespace string) (*v1.Service, error) {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name + "-" + "service",
			Namespace: namespace,
			Annotations: map[string]string{
				"uuid": req.Uuid,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: v1.SchemeGroupVersion.String(),
					Kind:       "Pod",
					Name:       pod.Name,
					UID:        pod.UID,
				},
			},
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": req.Name,
			},
			Ports: []v1.ServicePort{
				{
					Name: "ssh",
					Port: SSHPort,
					TargetPort: intstr.IntOrString{
						IntVal: SSHPort,
					},
				},
				{
					Name: "jupyter",
					Port: JupyterPort,
					TargetPort: intstr.IntOrString{
						IntVal: JupyterPort,
					},
				},
			},
		},
	}, nil
}

func ParseResources(res types.Resources) v1.ResourceRequirements {
	requirements := v1.ResourceRequirements{
		Requests: make(v1.ResourceList),
		Limits:   make(v1.ResourceList),
	}

	if res.CPU != "" {
		if cpu, err := resource.ParseQuantity(res.CPU); err == nil {
			requirements.Requests[v1.ResourceCPU] = cpu
			requirements.Limits[v1.ResourceCPU] = cpu
		}
	}

	if res.Memory != "" {
		if mem, err := resource.ParseQuantity(res.Memory); err == nil {
			requirements.Requests[v1.ResourceMemory] = mem
			requirements.Limits[v1.ResourceMemory] = mem
		}
	}

	if res.Storage.SystemDisk != "" {
		if disk, err := resource.ParseQuantity(res.Storage.SystemDisk); err == nil {
			requirements.Requests[v1.ResourceEphemeralStorage] = disk
		}
	}

	if res.GPU.Num != "" {
		if gpu, err := resource.ParseQuantity(res.GPU.Num); err == nil {
			requirements.Requests[v1.ResourceName("nvidia.com/gpu")] = gpu
			requirements.Limits[v1.ResourceName("nvidia.com/gpu")] = gpu
		}
	}

	return requirements
}

func ParseLabels(name, labelStr string) (map[string]string, error) {
	labels := map[string]string{
		"app": name,
	}
	if len(labelStr) == 0 {
		return labels, nil
	}

	labelStrs := strings.Split(labelStr, ",")
	for ix := range labelStrs {
		labelStr := strings.Split(labelStrs[ix], "=")
		if len(labelStr) != 2 {
			return nil, fmt.Errorf("unexpected label spec: %s", labelStrs[ix])
		}
		if len(labelStr[0]) == 0 {
			return nil, fmt.Errorf("unexpected empty label key")
		}
		labels[labelStr[0]] = labelStr[1]
	}
	return labels, nil
}
