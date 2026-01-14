package k8s

import (
	//"context"
	"fmt"
	"path/filepath"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Client struct {
	Typed kubernetes.Interface
}

func Connect() (*Client, error) {
	// Try in-cluster config first (if running inside a Pod)
	config, err := rest.InClusterConfig()
	if err != nil {
		// Fallback to kubeconfig (for local development)
		fmt.Println("In-cluster config not found, falling back to local kubeconfig...")
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{Typed: clientset}, nil
}
