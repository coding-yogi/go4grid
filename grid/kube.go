package grid

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type KubeApi struct {
	client    *kubernetes.Clientset
	namespace string
}

func (k8s *KubeApi) GetDeployments(label string) (*appsv1.DeploymentList, error) {
	deploymentsClient := k8s.client.AppsV1().Deployments(k8s.namespace)
	return deploymentsClient.List(context.TODO(), metav1.ListOptions{LabelSelector: label})
}

func (k8s *KubeApi) CreateDeployment(d *appsv1.Deployment) (*appsv1.Deployment, error) {
	deploymentsClient := k8s.client.AppsV1().Deployments(k8s.namespace)
	return deploymentsClient.Create(context.TODO(), d, metav1.CreateOptions{})
}

func (k8s *KubeApi) DeleteDeployment(name string) error {
	deploymentsClient := k8s.client.AppsV1().Deployments(k8s.namespace)
	return deploymentsClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (k8s *KubeApi) UpdateDeployment(d *appsv1.Deployment) (*appsv1.Deployment, error) {
	deploymentsClient := k8s.client.AppsV1().Deployments(k8s.namespace)
	return deploymentsClient.Update(context.TODO(), d, metav1.UpdateOptions{})
}

func (k8s *KubeApi) CreateService(s *corev1.Service) (*corev1.Service, error) {
	serviceClient := k8s.client.CoreV1().Services(k8s.namespace)
	return serviceClient.Create(context.TODO(), s, metav1.CreateOptions{})
}

func (k8s *KubeApi) GetServices(label string) (*corev1.ServiceList, error) {
	serviceClient := k8s.client.CoreV1().Services(k8s.namespace)
	return serviceClient.List(context.TODO(), metav1.ListOptions{LabelSelector: label})
}

func (k8s *KubeApi) DeleteService(name string) error {
	serviceClient := k8s.client.CoreV1().Services(k8s.namespace)
	return serviceClient.Delete(context.TODO(), name, metav1.DeleteOptions{})
}
