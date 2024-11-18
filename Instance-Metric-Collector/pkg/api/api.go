package api

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

func NewClientset() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Errorln(err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Errorln(err)
	}

	return kubeClient
}
