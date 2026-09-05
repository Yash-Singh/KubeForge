package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	reconcileTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubeforge_reconcile_total",
			Help: "Total number of reconciliation attempts",
		},
		[]string{"namespace", "name", "result"},
	)

	reconcileDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "kubeforge_reconcile_duration_seconds",
			Help:    "Duration of reconciliation attempts in seconds",
			Buckets: prometheus.ExponentialBuckets(0.01, 2, 15),
		},
		[]string{"namespace", "name"},
	)

	reconcileErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubeforge_reconcile_errors_total",
			Help: "Total number of reconciliation errors by type",
		},
		[]string{"namespace", "name", "error_type"},
	)

	applicationPhase = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kubeforge_application_phase",
			Help: "Current phase of the application (1=Ready, 2=Progressing, 3=Degraded)",
		},
		[]string{"namespace", "name"},
	)

	applicationReplicas = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "kubeforge_application_replicas",
			Help: "Desired and ready replica counts for the application",
		},
		[]string{"namespace", "name", "state"},
	)

	resourceCreated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubeforge_resource_created_total",
			Help: "Total number of child resources created",
		},
		[]string{"namespace", "name", "resource_type"},
	)

	resourceUpdated = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kubeforge_resource_updated_total",
			Help: "Total number of child resources updated",
		},
		[]string{"namespace", "name", "resource_type"},
	)
)

func init() {
	metrics.Registry.MustRegister(
		reconcileTotal,
		reconcileDuration,
		reconcileErrors,
		applicationPhase,
		applicationReplicas,
		resourceCreated,
		resourceUpdated,
	)
}

const (
	phaseReady       = "Ready"
	phaseProgressing = "Progressing"
	phaseDegraded    = "Degraded"
)
