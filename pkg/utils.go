package pkg

import (
	_ "fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

func GetClient() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()

	if err == nil {
		klog.Info("Using in-cluster config")
		return kubernetes.NewForConfig(config)
	} else {
		kubeconfigPath := os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			kubeconfigPath = os.Getenv("HOME") + "/.kube/config"
		}
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err == nil {
			klog.Info("Using k8s config: ", kubeconfigPath)
			return kubernetes.NewForConfig(config)
		} else {
			return nil, err
		}
	}
}
