package metricsv2

import (
	"context"
	"github.com/kyma-project/kyma-environment-broker/internal/broker"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

type BindDurationCollector struct {
	bindHistorgam   *prometheus.HistogramVec
	unbindHistogram *prometheus.HistogramVec
	logger          logrus.FieldLogger
}

func NewBindDurationCollector(logger logrus.FieldLogger) *BindDurationCollector {
	return &BindDurationCollector{
		bindHistorgam: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: prometheusNamespacev2,
			Subsystem: prometheusSubsystemv2,
			Name:      "bind_duration_millisecond",
			Help:      "The time of the bind request",
			//Buckets:   prometheus.LinearBuckets(10, 2, 56),
		}, []string{}),
		unbindHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: prometheusNamespacev2,
			Subsystem: prometheusSubsystemv2,
			Name:      "unbind_duration_millisecond",
			Help:      "The time of the unbind request",
			//Buckets:   prometheus.LinearBuckets(10, 2, 56),
		}, []string{}),
		logger: logger,
	}
}

func (c *BindDurationCollector) Describe(ch chan<- *prometheus.Desc) {
	c.bindHistorgam.Describe(ch)
	c.unbindHistogram.Describe(ch)
}

func (c *BindDurationCollector) Collect(ch chan<- prometheus.Metric) {
	c.bindHistorgam.Collect(ch)
	c.unbindHistogram.Collect(ch)
}

func (c *BindDurationCollector) OnBindingExecuted(ctx context.Context, ev interface{}) error {
	obj := ev.(broker.BindRequestProcessed)
	c.bindHistorgam.WithLabelValues().Observe(float64(obj.ProcessingDuration.Milliseconds()))
	return nil
}

func (c *BindDurationCollector) OnUnbindingExecuted(ctx context.Context, ev interface{}) error {
	obj := ev.(broker.UnbindRequestProcessed)
	c.bindHistorgam.WithLabelValues().Observe(float64(obj.ProcessingDuration.Milliseconds()))
	return nil
}

type BindingCreationCollector struct {
	bindingCreated *prometheus.CounterVec
}

func NewBindingCreationCollector() *BindingCreationCollector {
	return &BindingCreationCollector{
		bindingCreated: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: prometheusNamespacev2,
			Subsystem: prometheusSubsystemv2,
			Name:      "binding_created_total",
			Help:      "The total number of created bindings",
		}, []string{"plan_id"}),
	}
}

func (c *BindingCreationCollector) Describe(ch chan<- *prometheus.Desc) {
	c.bindingCreated.Describe(ch)
}

func (c *BindingCreationCollector) Collect(ch chan<- prometheus.Metric) {
	c.bindingCreated.Collect(ch)
}

func (c *BindingCreationCollector) OnBindingCreated(ctx context.Context, ev interface{}) error {
	obj := ev.(broker.BindingCreated)
	c.bindingCreated.WithLabelValues(obj.PlanID).Inc()
	return nil
}
