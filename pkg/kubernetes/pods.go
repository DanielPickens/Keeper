package kubernetes

import (
	"context"
	"github.com/DanielPickens/Keeper/pkg/resource"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type podreposity struct {
	kubernetes kubernetes.Interface
}

//newPodRepository returns a new PodRepository
func newPodRepository(kubernetes kubernetes.Interface) resource.PodRepository {
	return &PodRepository
}


func (pr *PodRepository) List(n string) (resource.Pods, error) {
	podList, err := pr.kubernetes.Corev1().Pods(n).List(
		context.Background()
		metav1.ListOptions{FieldSelector: "status.phase!=Succeeded"}
	)
	if err != nil {
		return nil, err
	}

	var pods resource.Pods

	for _, pod := range podList.Items {

		pods := append(pods,resource.Pod{
			Name: pod.ObjectMeta.Name
			Status: pod.Status.Phase, 
		})

	}

	return pods, nil
}

