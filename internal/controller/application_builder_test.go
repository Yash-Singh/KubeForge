package controller

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	platformv1alpha1 "github.com/Yash-Singh/KubeForge/api/v1alpha1"
)

var _ = Describe("applicationLabels", func() {
	It("should return the expected labels", func() {
		app := &platformv1alpha1.Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "checkout",
				Namespace: "apps",
			},
		}
		labels := applicationLabels(app)
		Expect(labels).To(HaveKeyWithValue(nameLabel, "checkout"))
		Expect(labels).To(HaveKeyWithValue(managedByLabel, managedByValue))
	})
})

var _ = Describe("desiredDeploymentSpec", func() {
	It("should build a DeploymentSpec from Application spec", func() {
		app := &platformv1alpha1.Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "checkout",
				Namespace: "apps",
			},
			Spec: platformv1alpha1.ApplicationSpec{
				Image:     "ghcr.io/example/checkout:v1.0.0",
				Replicas:  int32Ptr(2),
				Resources: resourceRequirements(),
				Env: []platformv1alpha1.EnvVar{
					{Name: "LOG_LEVEL", Value: "info"},
				},
			},
		}
		spec := desiredDeploymentSpec(app, *app.Spec.Replicas)

		Expect(spec.Template.Spec.Containers).To(HaveLen(1))
		c := spec.Template.Spec.Containers[0]
		Expect(c.Name).To(Equal("checkout"))
		Expect(c.Image).To(Equal("ghcr.io/example/checkout:v1.0.0"))
		Expect(c.Resources).To(Equal(*app.Spec.Resources))
		Expect(c.Env).To(Equal([]corev1.EnvVar{{Name: "LOG_LEVEL", Value: "info"}}))

		Expect(*spec.Replicas).To(Equal(int32(2)))
		Expect(spec.Selector.MatchLabels).To(Equal(map[string]string{nameLabel: "checkout"}))
		Expect(spec.Template.Labels).To(Equal(applicationLabels(app)))
	})

	It("should default replicas to 1 when nil", func() {
		app := &platformv1alpha1.Application{
			ObjectMeta: metav1.ObjectMeta{Name: "checkout"},
			Spec:       platformv1alpha1.ApplicationSpec{Image: "img"},
		}
		spec := desiredDeploymentSpec(app, 1)
		Expect(*spec.Replicas).To(Equal(int32(1)))
	})

	It("should not include env when empty", func() {
		app := &platformv1alpha1.Application{
			ObjectMeta: metav1.ObjectMeta{Name: "checkout"},
			Spec:       platformv1alpha1.ApplicationSpec{Image: "img"},
		}
		spec := desiredDeploymentSpec(app, 1)
		Expect(spec.Template.Spec.Containers[0].Env).To(BeNil())
	})
})

var _ = Describe("convertEnvVars", func() {
	It("should convert custom EnvVar slice to corev1.EnvVar slice", func() {
		input := []platformv1alpha1.EnvVar{
			{Name: "FOO", Value: "bar"},
			{Name: "BAZ", Value: "qux"},
		}
		output := convertEnvVars(input)
		Expect(output).To(Equal([]corev1.EnvVar{
			{Name: "FOO", Value: "bar"},
			{Name: "BAZ", Value: "qux"},
		}))
	})

	It("should return nil for empty input", func() {
		Expect(convertEnvVars(nil)).To(BeNil())
		Expect(convertEnvVars([]platformv1alpha1.EnvVar{})).To(BeNil())
	})
})

var _ = Describe("isDeploymentDegraded", func() {
	It("should return true when Available condition is False", func() {
		dep := &appsv1.Deployment{
			Status: appsv1.DeploymentStatus{
				Conditions: []appsv1.DeploymentCondition{
					{Type: "Available", Status: corev1.ConditionFalse},
				},
			},
		}
		Expect(isDeploymentDegraded(dep)).To(BeTrue())
	})

	It("should return true when Progressing has reason ProgressDeadlineExceeded", func() {
		dep := &appsv1.Deployment{
			Status: appsv1.DeploymentStatus{
				Conditions: []appsv1.DeploymentCondition{
					{Type: "Progressing", Status: corev1.ConditionTrue, Reason: "ProgressDeadlineExceeded"},
				},
			},
		}
		Expect(isDeploymentDegraded(dep)).To(BeTrue())
	})

	It("should return false for a healthy deployment", func() {
		dep := &appsv1.Deployment{
			Status: appsv1.DeploymentStatus{
				Conditions: []appsv1.DeploymentCondition{
					{Type: "Available", Status: corev1.ConditionTrue},
					{Type: "Progressing", Status: corev1.ConditionTrue, Reason: "NewReplicaSetAvailable"},
				},
			},
		}
		Expect(isDeploymentDegraded(dep)).To(BeFalse())
	})

	It("should return false for a deployment with no conditions", func() {
		dep := &appsv1.Deployment{}
		Expect(isDeploymentDegraded(dep)).To(BeFalse())
	})
})

func int32Ptr(i int32) *int32 {
	return &i
}

func resourceRequirements() *corev1.ResourceRequirements {
	return &corev1.ResourceRequirements{
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("100m"),
			corev1.ResourceMemory: resource.MustParse("128Mi"),
		},
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("512Mi"),
		},
	}
}
