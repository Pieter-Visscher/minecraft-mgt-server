package minecraft

import (
	"minecraft-mgt-server/k8s"

	"context"
	"log"

	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type Manager struct {
	K8s *k8s.Client
}

func (m *Manager) newBusyboxDeployment(name string) *appsv1.Deployment {
	replicas := int32(1)
	log.Printf("Starting deploy of new pod")
	
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": "minecraft-test",
				"instance": name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "minecraft-test",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "minecraft-test",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "busybox",
							Image:   "busybox:latest",
							// This command keeps the container running so you can inspect it
							Command: []string{"sh", "-c", "while true; do sleep 3600; done"},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse("100m"),
									corev1.ResourceMemory: resource.MustParse("128Mi"),
								},
							},
						},
					},
				},
			},
		},
	}
}

// CreateServer executes the creation in the cluster
func (m *Manager) CreateServer(ctx context.Context, name string) error {
	deployment := m.newBusyboxDeployment(name)
	
	_, err := m.K8s.Typed.AppsV1().Deployments("default").Create(ctx, deployment, metav1.CreateOptions{})
	return err
}
