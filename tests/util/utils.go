package util

import (
	"os"

	celvalidatortypes "github.com/accuknox/kubernetes-cel-validator/resource-cel-validator/pkg/types"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"sigs.k8s.io/yaml"
)

func GetDeployment(deploymentFilePath string) (*appsv1.Deployment, error) {
	bytes, err := os.ReadFile(deploymentFilePath)
	if err != nil {
		return nil, err
	}
	deployment := &appsv1.Deployment{}
	err = yaml.Unmarshal(bytes, &deployment)
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func GetClusterRole(clusterRoleFilePath string) (*rbacv1.ClusterRole, error) {
	bytes, err := os.ReadFile(clusterRoleFilePath)
	if err != nil {
		return nil, err
	}
	clusterRole := &rbacv1.ClusterRole{}
	err = yaml.Unmarshal(bytes, &clusterRole)
	if err != nil {
		return nil, err
	}
	return clusterRole, nil
}

func GetPod(podFilePath string) (*v1.Pod, error) {
	bytes, err := os.ReadFile(podFilePath)
	if err != nil {
		return nil, err
	}
	pod := &v1.Pod{}
	err = yaml.Unmarshal(bytes, &pod)
	if err != nil {
		return nil, err
	}
	return pod, nil
}

func GetKubernetesResourcePrecondition(kubernetesResourcePreconditionPath string) (*celvalidatortypes.KubernetesResourcePrecondition, error) {
	bytes, err := os.ReadFile(kubernetesResourcePreconditionPath)
	if err != nil {
		return nil, err
	}
	kubernetesResourcePrecondition := &celvalidatortypes.KubernetesResourcePrecondition{}
	err = yaml.Unmarshal(bytes, &kubernetesResourcePrecondition)
	if err != nil {
		return nil, err
	}
	return kubernetesResourcePrecondition, nil
}
