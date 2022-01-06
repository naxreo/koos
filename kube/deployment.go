package kube

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetDeploymentList(clientset *kubernetes.Clientset, namespace string) ([]string, error) {
	var retval []string
	deploy, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return retval, err
	}

	for _, v := range deploy.Items {
		//fmt.Printf("%v %v\n", i, v.Name)
		retval = append(retval, v.Name)
	}
	return retval, nil
}

func IsDeployment(clientset *kubernetes.Clientset, namespace string, name string) bool {
	_, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false
	}
	//fmt.Printf("%v", deploy.Name)
	return true
}
