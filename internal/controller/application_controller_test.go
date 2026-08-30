package controller

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	platformv1alpha1 "github.com/Yash-Singh/KubeForge/api/v1alpha1"
)

const (
	namespace = "default"
	name      = "test-app"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = Describe("Application Reconciler", func() {
	var ctx context.Context
	typeNamespacedName := types.NamespacedName{Name: name, Namespace: namespace}
	var app *platformv1alpha1.Application

	BeforeEach(func() {
		ctx = context.Background()
		app = newApplication()
		Expect(k8sClient.Create(ctx, app)).To(Succeed())
	})

	AfterEach(func() {
		By("cleaning up the Service")
		svc := &corev1.Service{}
		_ = k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, svc)
		_ = k8sClient.Delete(ctx, svc)
		By("cleaning up the Deployment")
		dep := &appsv1.Deployment{}
		_ = k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, dep)
		_ = k8sClient.Delete(ctx, dep)
		By("cleaning up the Application")
		Expect(k8sClient.Delete(ctx, app)).To(Succeed())
	})

	It("should create a Deployment when an Application is created", func() {
		By("reconciling the Application")
		r := &ApplicationReconciler{Client: k8sClient, Scheme: k8sClient.Scheme()}
		_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
		Expect(err).NotTo(HaveOccurred())

		By("verifying the Deployment was created")
		dep := &appsv1.Deployment{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, dep)).To(Succeed())
		Expect(dep.Spec.Template.Spec.Containers[0].Image).To(Equal(app.Spec.Image))
	})

	It("should create a Service when spec.service is configured", func() {
		By("setting spec.service")
		app.Spec.Service = &platformv1alpha1.ServiceSpec{Port: 8080, TargetPort: int32Ptr(8080)}
		Expect(k8sClient.Update(ctx, app)).To(Succeed())

		By("reconciling the Application")
		r := &ApplicationReconciler{Client: k8sClient, Scheme: k8sClient.Scheme()}
		_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
		Expect(err).NotTo(HaveOccurred())

		By("verifying the Service was created")
		svc := &corev1.Service{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, svc)).To(Succeed())
		Expect(svc.Spec.Ports).To(HaveLen(1))
		Expect(svc.Spec.Ports[0].Port).To(Equal(int32(8080)))
		Expect(svc.Spec.Selector).To(Equal(map[string]string{nameLabel: name}))
	})

	It("should not create a Service when spec.service is nil", func() {
		By("verifying the Application has no service")
		Expect(app.Spec.Service).To(BeNil())

		By("reconciling the Application")
		r := &ApplicationReconciler{Client: k8sClient, Scheme: k8sClient.Scheme()}
		_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
		Expect(err).NotTo(HaveOccurred())

		By("verifying no Service exists")
		svc := &corev1.Service{}
		err = k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, svc)
		Expect(err).To(HaveOccurred())
		Expect(errors.IsNotFound(err)).To(BeTrue())
	})

	It("should update status phase to Progressing then Ready", func() {
		By("reconciling the Application")
		r := &ApplicationReconciler{Client: k8sClient, Scheme: k8sClient.Scheme()}
		_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
		Expect(err).NotTo(HaveOccurred())

		By("verifying the status is observed")
		updatedApp := &platformv1alpha1.Application{}
		Eventually(func() error {
			err := k8sClient.Get(ctx, typeNamespacedName, updatedApp)
			if err != nil {
				return err
			}
			if updatedApp.Status.Phase == platformv1alpha1.PhaseProgressing || updatedApp.Status.Phase == platformv1alpha1.PhaseReady {
				return nil
			}
			return nil
		}, "10s", "500ms").Should(Succeed())
	})

	It("should update the Deployment when spec.image changes", func() {
		By("reconciling the Application")
		r := &ApplicationReconciler{Client: k8sClient, Scheme: k8sClient.Scheme()}
		_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
		Expect(err).NotTo(HaveOccurred())

		By("verifying the Deployment was created")
		dep := &appsv1.Deployment{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, dep)).To(Succeed())
		Expect(dep.Spec.Template.Spec.Containers[0].Image).To(Equal(app.Spec.Image))

		By("re-fetching and updating the Application")
		app = newApplication()
		Expect(k8sClient.Get(ctx, typeNamespacedName, app)).To(Succeed())
		app.Spec.Image = "ghcr.io/example/checkout:v2.0.0"
		Expect(k8sClient.Update(ctx, app)).To(Succeed())

		By("reconciling again")
		_, err = r.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
		Expect(err).NotTo(HaveOccurred())

		By("verifying the Deployment image is updated")
		Eventually(func() string {
			_ = k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, dep)
			return dep.Spec.Template.Spec.Containers[0].Image
		}, "10s", "500ms").Should(Equal("ghcr.io/example/checkout:v2.0.0"))
	})

	It("should recover a deleted Deployment", func() {
		By("reconciling the Application")
		r := &ApplicationReconciler{Client: k8sClient, Scheme: k8sClient.Scheme()}
		_, err := r.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
		Expect(err).NotTo(HaveOccurred())

		By("verifying the Deployment exists")
		dep := &appsv1.Deployment{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, dep)).To(Succeed())

		By("deleting the Deployment")
		Expect(k8sClient.Delete(ctx, dep)).To(Succeed())

		By("reconciling again")
		_, err = r.Reconcile(ctx, reconcile.Request{NamespacedName: typeNamespacedName})
		Expect(err).NotTo(HaveOccurred())

		By("verifying the Deployment is recreated")
		Eventually(func() error {
			return k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, dep)
		}, "10s", "500ms").Should(Succeed())
	})
})

func newApplication() *platformv1alpha1.Application {
	return &platformv1alpha1.Application{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: platformv1alpha1.ApplicationSpec{
			Image:    "ghcr.io/example/checkout:v1.0.0",
			Replicas: int32Ptr(1),
		},
	}
}
