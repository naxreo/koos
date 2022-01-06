package kube

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type serviceMeta struct {
	svcName string
	svcIp   string
	svcType string
}

func GetServiceList(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	var retval []string
	svc, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return retval, err
	}
	for _, v := range svc.Items {
		retval = append(retval, v.Name)
	}
	return retval, nil
}

func IsService(clientset *kubernetes.Clientset, namespace string, name string) bool {
	_, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false
	}
	//fmt.Printf("%v %v %v", svc.Name, svc.Spec.ClusterIP, svc.Spec.Type)
	return true
}

func GetService(clientset *kubernetes.Clientset, namespace string, name string) ([3]string, error) {
	var retval [3]string
	svc, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return retval, err
	}

	retval = [3]string{svc.Name, svc.Spec.ClusterIP, string(svc.Spec.Type)}
	return retval, nil
}
