package cslb

import (
	"errors"
	"math"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	dto "github.com/prometheus/client_model/go"
)

const (
	DefaultMinSampleSize = 10
)

var (
	SampleNotEnoughError = errors.New("sample not enough")
	InvalidRatioError    = errors.New("invalid ratio error")
)

type Metrics struct {
	MinSampleSize int

	registerer prometheus.Registerer

	nodeCounter       *prometheus.CounterVec
	nodeFailedCounter *prometheus.CounterVec
}

// NewMetrics initializes a Metrics instance, it might return nil if it's unnecessary or invalid
func NewMetrics(maxNodeFailedRatio float64, minSampleSize int) *Metrics {
	// No metric is needed if maxNodeFailedRatio is invalid or equal to NodeFailedUnlimited
	if maxNodeFailedRatio < NodeFailedAny || maxNodeFailedRatio >= NodeFailedUnlimited {
		return nil
	}

	registerer := prometheus.NewRegistry()
	return &Metrics{
		MinSampleSize: minSampleSize,
		registerer:    registerer,
		nodeCounter: promauto.With(registerer).NewCounterVec(prometheus.CounterOpts{
			Name: "load_balancer_node_counter",
			Help: "Total number of node request",
		}, []string{"node"}),
		nodeFailedCounter: promauto.With(registerer).NewCounterVec(prometheus.CounterOpts{
			Name: "load_balancer_node_failed_counter",
			Help: "Total number of failed node request",
		}, []string{"node"}),
	}
}

func (m *Metrics) NodeInc(node Node) {
	m.nodeCounter.WithLabelValues(node.String()).Inc()
}

func (m *Metrics) NodeFailedInc(node Node) {
	m.nodeFailedCounter.WithLabelValues(node.String()).Inc()
}

func (m *Metrics) ResetNode(node Node) {
	m.nodeCounter.DeleteLabelValues(node.String())
	m.nodeFailedCounter.DeleteLabelValues(node.String())
}

func (m *Metrics) ResetAllNodes() {
	m.nodeCounter.Reset()
	m.nodeFailedCounter.Reset()
}

func getCounterValue(counter prometheus.Counter) (float64, error) {
	metric := &dto.Metric{}
	if err := counter.Write(metric); err != nil {
		return 0.0, err
	}
	return metric.Counter.GetValue(), nil
}

func (m *Metrics) GetNodeFailedRatio(node Node) (float64, error) {
	key := node.String()
	nodeCounter, _ := m.nodeCounter.GetMetricWithLabelValues(key)
	total, err := getCounterValue(nodeCounter)
	if err != nil {
		return 0.0, err
	}
	if total < float64(m.MinSampleSize) {
		return 0.0, SampleNotEnoughError
	}
	nodeFailedCounter, _ := m.nodeFailedCounter.GetMetricWithLabelValues(key)
	failed, err := getCounterValue(nodeFailedCounter)
	if err != nil {
		return 0.0, err
	}
	ratio := failed / total
	if math.IsNaN(ratio) {
		return 0.0, InvalidRatioError
	}
	return math.Max(math.Min(ratio, NodeFailedUnlimited), NodeFailedAny), nil
}
