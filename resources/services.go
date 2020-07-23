package resources

import (
	apiv1 "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func HubService() *corev1.Service {

	name := HubName

	labels := make(map[string]string)
	labels["app"] = name

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
			Name:   name,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "hub-port",
					Protocol:   apiv1.ProtocolTCP,
					Port:       4444,
					TargetPort: intstr.IntOrString{IntVal: 4444},
				},
				{
					Name:       "pub-port",
					Protocol:   apiv1.ProtocolTCP,
					Port:       4442,
					TargetPort: intstr.IntOrString{IntVal: 4442},
				},
				{
					Name:       "sub-port",
					Protocol:   apiv1.ProtocolTCP,
					Port:       4443,
					TargetPort: intstr.IntOrString{IntVal: 4443},
				},
			},
			Selector: labels,
			Type:     apiv1.ServiceTypeNodePort,
		},
	}

}
