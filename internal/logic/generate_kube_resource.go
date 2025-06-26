package logic

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gpu-container-service/internal/types"
	"gpu-container-service/pkg"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
)

var (
	SSHPort           int32 = 22
	JupyterPort       int32 = 8888
	GPUContainerLabel       = "mfy.com/gpu-container"
	GPUtypeLabel            = "mfy.com/gpu-type"
	DeafultDomain           = "containercloud-xian.xaidc.com"
)

func GenerateDeploymentTemplate(req *types.CreateInstanceRequest, namespace string) (*appsv1.Deployment, error) {
	logx.Info("Generate New Deployment Template")
	labels, err := ParseLabels(req.Name, req.Labels)
	if err != nil {
		return nil, fmt.Errorf("Parse Label Error: %v", err)
	}

	podTemplate, err := GeneratePodTemplate(req, labels)
	if err != nil {
		return nil, fmt.Errorf("Generate Pod Template Error: %v", err)
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: req.Name + "-",
			Namespace:    namespace,
			Labels:       labels,
			Annotations:  map[string]string{},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &req.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": req.Name,
				},
			},
			Template: podTemplate,
		},
	}, nil
}

func GeneratePodTemplate(req *types.CreateInstanceRequest, labels map[string]string) (v1.PodTemplateSpec, error) {
	logx.Info("Generate New Pod Template")

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

	container := v1.Container{
		Name:  req.Name,
		Image: req.Image,
		Ports: []v1.ContainerPort{
			{
				Name:          "ssh",
				ContainerPort: SSHPort,
			},
			{
				Name:          "http",
				ContainerPort: JupyterPort,
			},
		},
		Resources: ParseResources(req.ResourceRequest),
		Env:       env,
		EnvFrom: []v1.EnvFromSource{
			{
				SecretRef: &v1.SecretEnvSource{
					LocalObjectReference: v1.LocalObjectReference{
						Name: req.Name + "-secret",
					},
				},
			},
		},
	}

	return v1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{container},
			NodeSelector: map[string]string{
				GPUContainerLabel: "true",
				GPUtypeLabel:      req.ResourceRequest.GPUType,
			},
		},
	}, nil
}

func GenerateSecretTemplate(req *types.CreateInstanceRequest, deployment *appsv1.Deployment, namespace string) (*v1.Secret, error) {
	password, hashedPassword, err := pkg.GenerateJupyterPassword()
	if err != nil {
		return nil, fmt.Errorf("GenerateJupyterPassword Error: %v", err)
	}

	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name + "-secret",
			Namespace: namespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: appsv1.SchemeGroupVersion.String(),
					Kind:       "Deployment",
					Name:       deployment.Name,
					UID:        deployment.UID,
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

func GenerateServiceTemplate(req *types.CreateInstanceRequest, deployment *appsv1.Deployment, namespace string) (*v1.Service, error) {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name + "-service",
			Namespace: namespace,
			Labels: map[string]string{
				"app":       req.Name,
				"create-by": "mck",
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: appsv1.SchemeGroupVersion.String(),
					Kind:       "Deployment",
					Name:       deployment.Name,
					UID:        deployment.UID,
					Controller: new(bool),
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
