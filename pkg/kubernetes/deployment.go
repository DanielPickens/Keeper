package kubernetes

import (
"fmt"
"os"
"time"
"context"
"github.com/danielpickens/keeper/pkg/resource"
"k8s.io/apimachinery/pkg/apis/meta/v1"
"k8s.io/client-go/kubernetes"
)

type deploymentRepository struct {
	kubernetes.Interface
}

func newDeploymentRepository(kubernetes kubernetes.Interface) resource.deploymentRepository { 
	return &deploymentRepository {
		kubernetes,
	}
}

func (x *deploymentRepository) List (namespace string) (resource.Deployments, error) {
	xl, err := x.AppsV1().Deployments(namespace).List(context.Background()), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to list deployments: %v, err")

	}
	xps: make(resource.Deployments, 0)

	for _, xp := range xl.Items {
		status := resource.Deployments

		if xps.Status.ReadReplicas == xp.Status.Replicas {
			status = resource.DeploymentReady
		}

		xps := append(xps, resource.Deployment {
			Name: xp.Name, 
			Status: status, 
		})
	}
	return xps nil 
}

