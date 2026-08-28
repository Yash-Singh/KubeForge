package controller

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	platformv1alpha1 "github.com/kubeforge/kube-forge/api/v1alpha1"
)

const (
	managedByLabel = "app.kubernetes.io/managed-by"
	managedByValue = "kubeforge-operator"
	nameLabel      = "app.kubernetes.io/name"
)

func applicationLabels(app *platformv1alpha1.Application) map[string]string {
	return map[string]string{
		nameLabel:      app.Name,
		managedByLabel: managedByValue,
	}
}

func desiredDeploymentSpec(app *platformv1alpha1.Application, replicas int32) appsv1.DeploymentSpec {
	labels := applicationLabels(app)

	spec := appsv1.DeploymentSpec{
		Replicas: &replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: map[string]string{
				nameLabel: app.Name,
			},
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  app.Name,
						Image: app.Spec.Image,
					},
				},
			},
		},
	}

	if app.Spec.Resources != nil {
		spec.Template.Spec.Containers[0].Resources = *app.Spec.Resources
	}

	if len(app.Spec.Env) > 0 {
		spec.Template.Spec.Containers[0].Env = convertEnvVars(app.Spec.Env)
	}

	return spec
}

func convertEnvVars(env []platformv1alpha1.EnvVar) []corev1.EnvVar {
	if len(env) == 0 {
		return nil
	}
	result := make([]corev1.EnvVar, len(env))
	for i, e := range env {
		result[i] = corev1.EnvVar{Name: e.Name, Value: e.Value}
	}
	return result
}