package resources

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const AppName = "go4grid-selenium4"
const HubName = AppName + "-hub"
const NodeName = AppName + "-node"

func HubDeployment() *appsv1.Deployment {

	name := HubName
	image := "selenium/hub:4.0.0"

	var replicas int32 = 1

	labels := make(map[string]string)
	labels["app"] = name

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
			Name:   name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
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
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 4442,
								},
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 4443,
								},
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 4444,
								},
							},
							ReadinessProbe: &apiv1.Probe{
								Handler: apiv1.Handler{
									HTTPGet: &apiv1.HTTPGetAction{
										Path:   "/wd/hub/status",
										Port:   intstr.IntOrString{IntVal: 4444},
										Scheme: "HTTP",
									},
								},
								InitialDelaySeconds: 30,
								PeriodSeconds:       5,
								SuccessThreshold:    2,
								TimeoutSeconds:      5,
							},
						},
					},
				},
			},
		},
	}
}

func NodeDeployment(browser string, replicas int32, hubHost, hubPort string) *appsv1.Deployment {

	name := NodeName + "-" + browser
	image := fmt.Sprintf("selenium/node-%s:4.0.0", browser)
	labels := make(map[string]string)
	labels["app"] = name

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
			Name:   name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
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
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 5555,
								},
							},
							Env: []apiv1.EnvVar{
								{
									Name:  "HUB_HOST",
									Value: hubHost,
								},
								{
									Name:  "HUB_PORT",
									Value: hubPort,
								},
							},
						},
					},
				},
			},
		},
	}
}
