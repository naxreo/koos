package kube

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetNamespaceList(clientset *kubernetes.Clientset) ([]string, error) {
	var ns []string

	p, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return ns, err
	}
	for _, s := range p.Items {
		ns = append(ns, s.Name)
	}

	return ns, nil
}
