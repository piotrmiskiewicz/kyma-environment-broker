package metricsv2

import (
	"context"
	"fmt"
	"github.com/kyma-project/kyma-environment-broker/internal/broker"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
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
			Buckets:   prometheus.LinearBuckets(50, 200, 15),
		}, []string{}),
		unbindHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: prometheusNamespacev2,
			Subsystem: prometheusSubsystemv2,
			Name:      "unbind_duration_millisecond",
			Help:      "The time of the unbind request",
			Buckets:   prometheus.LinearBuckets(50, 200, 15),
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

type BindingStatitics struct {
	db     storage.Bindings
	logger logrus.FieldLogger

	sync                                 sync.Mutex
	poolingInterval                      time.Duration
	MinutesSinceEarliestExpirationMetric prometheus.Gauge
}

func NewBindingStatsCollector(db storage.Bindings, poolingInterval time.Duration, logger logrus.FieldLogger) *BindingStatitics {
	return &BindingStatitics{
		db:              db,
		logger:          logger,
		poolingInterval: poolingInterval,
		MinutesSinceEarliestExpirationMetric: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: prometheusNamespacev2,
			Subsystem: prometheusSubsystemv2,
			Name:      "minutes_since_earliest_binding_expiration",
			Help:      "Specifies the time in minutes since the earliest binding expiration.",
		}),
	}
}

func (c *BindingStatitics) MustRegister(ctx context.Context) {
	prometheus.MustRegister(c.MinutesSinceEarliestExpirationMetric)
	go c.Job(ctx)
}

func (c *BindingStatitics) Job(ctx context.Context) {
	defer func() {
		if recovery := recover(); recovery != nil {
			c.logger.Errorf("panic recovered while handling in progress operation counter: %v", recovery)
		}
	}()

	if err := c.updateMetrics(); err != nil {
		c.logger.Error("failed to update metrics metrics", err)
	}

	ticker := time.NewTicker(c.poolingInterval)
	for {
		select {
		case <-ticker.C:
			if err := c.updateMetrics(); err != nil {
				c.logger.Error("failed to update operation stats metrics: ", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *BindingStatitics) updateMetrics() error {
	defer c.sync.Unlock()
	c.sync.Lock()

	stats, err := c.db.GetStatistics()
	if err != nil {
		return fmt.Errorf("cannot fetch in progress metrics from operations : %s", err.Error())
	}
	c.MinutesSinceEarliestExpirationMetric.Set(stats.MinutesSinceEarliestExpiration)
	return nil
}
