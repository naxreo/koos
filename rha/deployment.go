package rha

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func int32Ptr(i int32) *int32 { return &i }

func CreateRHADeployment(clientset *kubernetes.Clientset, namespace string, name string, img string) bool {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}
	if img == "" {
		img = "bitnami/redis:6.2.6"
	}

	deployset := clientset.AppsV1().Deployments(namespace)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  name,
							Image: img,
							Env: []apiv1.EnvVar{
								{
									Name:  "REDIS_PASSWORD",
									Value: "ChangePassword",
								},
								{
									Name:  "ALLOW_EMPTY_PASSWORD",
									Value: "no",
								},
								{
									Name:  "REDIS_AOF_ENABLED",
									Value: "no",
								},
								{
									Name:  "REDIS_ALLOW_REMOTE_CONNECTIONS",
									Value: "yes",
								},
								{
									Name:  "REDIS_REPLICATION_MODE",
									Value: "master",
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "redis",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 6379,
								},
							},
						},
					},
				},
			},
		},
	}

	fmt.Printf("Creating deployment %s with %s ...\n", name, img)
	result, err := deployset.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error %q\n", err)
		return false
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return true
}

func CreateRHAHeadless(clientset *kubernetes.Clientset, namespace string, name string, selector string) bool {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}

	serviceset := clientset.CoreV1().Services(namespace)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:     "redis",
					Protocol: apiv1.ProtocolTCP,
					Port:     6379,
				},
			},
			Selector: map[string]string{
				"app": selector,
			},
			ClusterIP: "None",
			Type:      "ClusterIP",
		},
	}

	fmt.Printf("Creating %s headless service...\n", name)
	result, err := serviceset.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error %q\n", err)
		return false
	}
	fmt.Printf("Created headless service %q.\n", result.GetObjectMeta().GetName())
	return true
}

func CreateRHAService(clientset *kubernetes.Clientset, namespace string, name string, selector string) bool {
	if namespace == "" {
		namespace = apiv1.NamespaceDefault
	}

	serviceset := clientset.CoreV1().Services(namespace)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:     "redis",
					Protocol: apiv1.ProtocolTCP,
					Port:     6379,
				},
			},
			Selector: map[string]string{
				"app": selector,
			},
			Type: "ClusterIP",
		},
	}

	fmt.Printf("Creating %s Service...\n", name)
	result, err := serviceset.Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Error %q\n", err)
		return false
	}
	fmt.Printf("Created Service %q.\n", result.GetObjectMeta().GetName())
	return true
}

func IsDeployment(clientset *kubernetes.Clientset, namespace string, name string) bool {
	ns, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false
	}
	fmt.Println(ns)
	return true
}
