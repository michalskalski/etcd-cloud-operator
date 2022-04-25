package sts

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func k8Client() *kubernetes.Clientset {
	config := ctrl.GetConfigOrDie()
	return kubernetes.NewForConfigOrDie(config)
}

func GetPodsForApp(clientset *kubernetes.Clientset, ctx context.Context, namespace, appName string) ([]corev1.Pod, error) {
	labelSelector := fmt.Sprintf("app=%s", appName)
	opts := metav1.ListOptions{LabelSelector: labelSelector}
	list, err := clientset.CoreV1().Pods(namespace).List(ctx, opts)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
