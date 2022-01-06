package kube

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func Config() *kubernetes.Clientset {
	// create in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("Not in-cluster env\n")
		var kubeconfig string
		home := homedir.HomeDir()
		if home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		} else {
			panic("There is no homedir")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("Not out-of-cluster either\n")
			panic(err.Error())
		}
	}
	
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}
