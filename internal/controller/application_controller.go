package controller

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/record"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	platformv1alpha1 "github.com/Yash-Singh/KubeForge/api/v1alpha1"
)

type ApplicationReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	DynamicClient dynamic.Interface
	Recorder      record.EventRecorder
}

// +kubebuilder:rbac:groups=platform.kubeforge.io,resources=applications,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=platform.kubeforge.io,resources=applications/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=platform.kubeforge.io,resources=applications/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=policy,resources=poddisruptionbudgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=keda.sh,resources=scaledobjects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=argoproj.io,resources=rollouts,verbs=get;list;watch;create;update;patch;delete

func (r *ApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	start := time.Now()
	logger := logf.FromContext(ctx)

	app := &platformv1alpha1.Application{}
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	logger = logger.WithValues("generation", app.Generation, "observedGeneration", app.Status.ObservedGeneration)

	replicas := ptr.To[int32](1)
	if app.Spec.Replicas != nil {
		replicas = app.Spec.Replicas
	}

	if err := r.reconcileDeployment(ctx, app, *replicas); err != nil {
		logger.Error(err, "Failed to reconcile Deployment")
		reconcileErrors.WithLabelValues(app.Namespace, app.Name, "deployment").Inc()
		r.Recorder.Event(app, corev1.EventTypeWarning, "ReconcileError", fmt.Sprintf("Failed to reconcile Deployment: %v", err))
		return ctrl.Result{}, err
	}

	if err := r.reconcilePodDisruptionBudget(ctx, app); err != nil {
		logger.Error(err, "Failed to reconcile PodDisruptionBudget")
		reconcileErrors.WithLabelValues(app.Namespace, app.Name, "pdb").Inc()
		r.Recorder.Event(app, corev1.EventTypeWarning, "ReconcileError", fmt.Sprintf("Failed to reconcile PodDisruptionBudget: %v", err))
		return ctrl.Result{}, err
	}

	if app.Spec.Service != nil {
		if err := r.reconcileService(ctx, app); err != nil {
			logger.Error(err, "Failed to reconcile Service")
			reconcileErrors.WithLabelValues(app.Namespace, app.Name, "service").Inc()
			r.Recorder.Event(app, corev1.EventTypeWarning, "ReconcileError", fmt.Sprintf("Failed to reconcile Service: %v", err))
			return ctrl.Result{}, err
		}
	}

	if err := r.reconcileNetworkPolicy(ctx, app); err != nil {
		logger.Error(err, "Failed to reconcile NetworkPolicy")
		reconcileErrors.WithLabelValues(app.Namespace, app.Name, "networkpolicy").Inc()
		r.Recorder.Event(app, corev1.EventTypeWarning, "ReconcileError", fmt.Sprintf("Failed to reconcile NetworkPolicy: %v", err))
		return ctrl.Result{}, err
	}

	if err := r.reconcileHorizontalPodAutoscaler(ctx, app); err != nil {
		logger.Error(err, "Failed to reconcile HorizontalPodAutoscaler")
		reconcileErrors.WithLabelValues(app.Namespace, app.Name, "hpa").Inc()
		r.Recorder.Event(app, corev1.EventTypeWarning, "ReconcileError", fmt.Sprintf("Failed to reconcile HorizontalPodAutoscaler: %v", err))
		return ctrl.Result{}, err
	}

	if err := r.reconcileKEDAScaledObject(ctx, app); err != nil {
		logger.Error(err, "Failed to reconcile KEDA ScaledObject")
		reconcileErrors.WithLabelValues(app.Namespace, app.Name, "keda").Inc()
		r.Recorder.Event(app, corev1.EventTypeWarning, "ReconcileError", fmt.Sprintf("Failed to reconcile KEDA ScaledObject: %v", err))
		return ctrl.Result{}, err
	}

	if err := r.reconcileArgoRollout(ctx, app); err != nil {
		logger.Error(err, "Failed to reconcile Argo Rollout")
		reconcileErrors.WithLabelValues(app.Namespace, app.Name, "rollout").Inc()
		r.Recorder.Event(app, corev1.EventTypeWarning, "ReconcileError", fmt.Sprintf("Failed to reconcile Argo Rollout: %v", err))
		return ctrl.Result{}, err
	}

	if err := r.updateStatus(ctx, app); err != nil {
		logger.Error(err, "Failed to update status")
		reconcileErrors.WithLabelValues(app.Namespace, app.Name, "status").Inc()
		return ctrl.Result{}, err
	}

	reconcileTotal.WithLabelValues(app.Namespace, app.Name, "success").Inc()
	reconcileDuration.WithLabelValues(app.Namespace, app.Name).Observe(time.Since(start).Seconds())

	return ctrl.Result{}, nil
}

func (r *ApplicationReconciler) reconcileDeployment(ctx context.Context, app *platformv1alpha1.Application, replicas int32) error {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
	}

	_, err := controllerutil.CreateOrPatch(ctx, r.Client, dep, func() error {
		if err := controllerutil.SetControllerReference(app, dep, r.Scheme); err != nil {
			return fmt.Errorf("setting owner reference: %w", err)
		}
		dep.Spec = desiredDeploymentSpec(app, replicas)
		return nil
	})
	return err
}

func (r *ApplicationReconciler) reconcilePodDisruptionBudget(ctx context.Context, app *platformv1alpha1.Application) error {
	pdb := desiredPodDisruptionBudget(app)
	if pdb == nil {
		return nil
	}

	_, err := controllerutil.CreateOrPatch(ctx, r.Client, pdb, func() error {
		if err := controllerutil.SetControllerReference(app, pdb, r.Scheme); err != nil {
			return fmt.Errorf("setting owner reference: %w", err)
		}
		return nil
	})
	return err
}

func (r *ApplicationReconciler) reconcileNetworkPolicy(ctx context.Context, app *platformv1alpha1.Application) error {
	np := desiredNetworkPolicy(app)
	if np == nil {
		return nil
	}

	_, err := controllerutil.CreateOrPatch(ctx, r.Client, np, func() error {
		if err := controllerutil.SetControllerReference(app, np, r.Scheme); err != nil {
			return fmt.Errorf("setting owner reference: %w", err)
		}
		return nil
	})
	return err
}

func (r *ApplicationReconciler) reconcileHorizontalPodAutoscaler(ctx context.Context, app *platformv1alpha1.Application) error {
	hpa := desiredHorizontalPodAutoscaler(app)
	if hpa == nil {
		return nil
	}

	_, err := controllerutil.CreateOrPatch(ctx, r.Client, hpa, func() error {
		if err := controllerutil.SetControllerReference(app, hpa, r.Scheme); err != nil {
			return fmt.Errorf("setting owner reference: %w", err)
		}
		return nil
	})
	return err
}

func (r *ApplicationReconciler) reconcileKEDAScaledObject(ctx context.Context, app *platformv1alpha1.Application) error {
	if app.Spec.KEDAScaledObject == nil || r.DynamicClient == nil {
		return nil
	}

	scaledObject := desiredKEDAScaledObject(app)

	existing, err := r.DynamicClient.Resource(scaledObjectGVR).Namespace(app.Namespace).Get(ctx, app.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = r.DynamicClient.Resource(scaledObjectGVR).Namespace(app.Namespace).Create(ctx, scaledObject, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	scaledObject.SetResourceVersion(existing.GetResourceVersion())
	_, err = r.DynamicClient.Resource(scaledObjectGVR).Namespace(app.Namespace).Update(ctx, scaledObject, metav1.UpdateOptions{})
	return err
}

func (r *ApplicationReconciler) reconcileArgoRollout(ctx context.Context, app *platformv1alpha1.Application) error {
	if app.Spec.ArgoRollout == nil || r.DynamicClient == nil {
		return nil
	}

	rollout := desiredArgoRollout(app)

	existing, err := r.DynamicClient.Resource(rolloutGVR).Namespace(app.Namespace).Get(ctx, app.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		_, err = r.DynamicClient.Resource(rolloutGVR).Namespace(app.Namespace).Create(ctx, rollout, metav1.CreateOptions{})
		return err
	}
	if err != nil {
		return err
	}

	rollout.SetResourceVersion(existing.GetResourceVersion())
	_, err = r.DynamicClient.Resource(rolloutGVR).Namespace(app.Namespace).Update(ctx, rollout, metav1.UpdateOptions{})
	return err
}

func (r *ApplicationReconciler) reconcileService(ctx context.Context, app *platformv1alpha1.Application) error {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name,
			Namespace: app.Namespace,
		},
	}

	port := app.Spec.Service.Port
	targetPort := port
	if app.Spec.Service.TargetPort != nil {
		targetPort = *app.Spec.Service.TargetPort
	}

	_, err := controllerutil.CreateOrPatch(ctx, r.Client, svc, func() error {
		if err := controllerutil.SetControllerReference(app, svc, r.Scheme); err != nil {
			return fmt.Errorf("setting owner reference: %w", err)
		}
		svc.Spec.Selector = map[string]string{
			nameLabel: app.Name,
		}
		svc.Spec.Ports = []corev1.ServicePort{
			{
				Protocol:   corev1.ProtocolTCP,
				Port:       port,
				TargetPort: intstr.FromInt32(targetPort),
			},
		}
		return nil
	})
	return err
}

func (r *ApplicationReconciler) updateStatus(ctx context.Context, app *platformv1alpha1.Application) error {
	logger := logf.FromContext(ctx)
	original := app.DeepCopy()
	previousPhase := app.Status.Phase

	app.Status.ObservedGeneration = app.Generation

	dep := &appsv1.Deployment{}
	key := client.ObjectKey{Name: app.Name, Namespace: app.Namespace}
	if err := r.Get(ctx, key, dep); err != nil {
		if apierrors.IsNotFound(err) {
			app.Status.Phase = platformv1alpha1.PhaseProgressing
			app.Status.DesiredReplicas = 0
			app.Status.ReadyReplicas = 0
			setCondition(app, platformv1alpha1.ConditionProgressing, "Reconciling", "Deployment not yet created")
		} else {
			return err
		}
	} else {
		app.Status.DesiredReplicas = ptr.Deref(dep.Spec.Replicas, 0)
		app.Status.ReadyReplicas = dep.Status.ReadyReplicas

		if isDeploymentDegraded(dep) {
			app.Status.Phase = platformv1alpha1.PhaseDegraded
			setCondition(app, platformv1alpha1.ConditionDegraded, "DeploymentFailed", "Deployment has unavailable or failing replicas")
		} else if dep.Status.ReadyReplicas >= ptr.Deref(dep.Spec.Replicas, 1) {
			app.Status.Phase = platformv1alpha1.PhaseReady
			setCondition(app, platformv1alpha1.ConditionReady, "ResourcesReady", "Application resources are ready")
		} else if dep.Status.ReadyReplicas > 0 {
			app.Status.Phase = platformv1alpha1.PhaseProgressing
			setCondition(app, platformv1alpha1.ConditionProgressing, "Scaling", "Waiting for all replicas to be ready")
		} else {
			app.Status.Phase = platformv1alpha1.PhaseProgressing
			setCondition(app, platformv1alpha1.ConditionProgressing, "Creating", "Waiting for first replica to become ready")
		}
	}

	if err := r.Status().Patch(ctx, app, client.MergeFrom(original)); err != nil {
		logger.Error(err, "Failed to patch Application status")
		return err
	}

	// Update metrics
	phaseValue := float64(0)
	switch app.Status.Phase {
	case phaseReady:
		phaseValue = 1
	case phaseProgressing:
		phaseValue = 2
	case phaseDegraded:
		phaseValue = 3
	}
	applicationPhase.WithLabelValues(app.Namespace, app.Name).Set(phaseValue)
	applicationReplicas.WithLabelValues(app.Namespace, app.Name, "desired").Set(float64(app.Status.DesiredReplicas))
	applicationReplicas.WithLabelValues(app.Namespace, app.Name, "ready").Set(float64(app.Status.ReadyReplicas))

	// Record events on phase transitions
	if previousPhase != app.Status.Phase {
		switch app.Status.Phase {
		case phaseReady:
			r.Recorder.Event(app, corev1.EventTypeNormal, "PhaseReady", fmt.Sprintf("Application is ready (%d/%d replicas)", app.Status.ReadyReplicas, app.Status.DesiredReplicas))
		case phaseProgressing:
			r.Recorder.Event(app, corev1.EventTypeNormal, "PhaseProgressing", fmt.Sprintf("Application is progressing (%d/%d replicas)", app.Status.ReadyReplicas, app.Status.DesiredReplicas))
		case phaseDegraded:
			r.Recorder.Event(app, corev1.EventTypeWarning, "PhaseDegraded", "Application is degraded")
		}
	}

	return nil
}

func setCondition(app *platformv1alpha1.Application, condType string, reason, message string) {
	now := metav1.Now()
	for i, c := range app.Status.Conditions {
		if c.Type == condType {
			if c.Status != metav1.ConditionTrue || c.Reason != reason || c.Message != message {
				app.Status.Conditions[i].Status = metav1.ConditionTrue
				app.Status.Conditions[i].Reason = reason
				app.Status.Conditions[i].Message = message
				app.Status.Conditions[i].LastTransitionTime = now
				app.Status.Conditions[i].ObservedGeneration = app.Generation
			}
			return
		}
	}
	app.Status.Conditions = append(app.Status.Conditions, metav1.Condition{
		Type:               condType,
		Status:             metav1.ConditionTrue,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: now,
		ObservedGeneration: app.Generation,
	})
}

func isDeploymentDegraded(dep *appsv1.Deployment) bool {
	for _, c := range dep.Status.Conditions {
		if c.Type == "Available" && c.Status == corev1.ConditionFalse {
			return true
		}
		if c.Type == "Progressing" && c.Reason == "ProgressDeadlineExceeded" {
			return true
		}
	}
	return false
}

func (r *ApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1alpha1.Application{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&policyv1.PodDisruptionBudget{}).
		Owns(&networkingv1.NetworkPolicy{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{}).
		Named("application").
		Complete(r)
}
