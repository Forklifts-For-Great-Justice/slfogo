package slfogolib

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricHolder struct {
	metrics map[[2]string]float64
	msgVec  *prometheus.GaugeVec
	mtx     sync.Mutex
}

func NewMetricHolder() *MetricHolder {
	return &MetricHolder{
		msgVec: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "ForkliftsForGreatJustice",
				Subsystem: "slfogo",
				Name:      "messages_per_second",
				Help:      "Message sent to syslog per second",
			},
			[]string{"hostname", "service"}),
		metrics: map[[2]string]float64{},
	}
}

func (mh *MetricHolder) AddMetric(key [2]string) {
	mh.mtx.Lock()
	defer mh.mtx.Unlock()
	mh.metrics[key] += 1
}

func (mh *MetricHolder) update() {
	mh.mtx.Lock()
	defer mh.mtx.Unlock()
	for k, v := range mh.metrics {
		if v == 0 {
			mh.msgVec.DeleteLabelValues(k[0], k[1])
			delete(mh.metrics, k)
		} else {
			mh.msgVec.WithLabelValues(k[0], k[1]).Set(float64(v))
			mh.metrics[k] = 0
		}
	}
}

func (mh *MetricHolder) GetGauge() *prometheus.GaugeVec {
	return mh.msgVec
}

func (mh *MetricHolder) HandleMetrics(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
handleLoop:
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				break handleLoop
			} else {
				mh.update()
			}

		case <-ctx.Done():
			break handleLoop
		}
	}

}
