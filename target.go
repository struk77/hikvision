package main

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Target struct {
	sync.Mutex
	host      string
	registry  *prometheus.Registry
	pirStatus *PIRAlarmStatus
}

func NewTarget(host string) *Target {
	t := Target{
		host:     host,
		registry: prometheus.NewRegistry()}
	log.Println("new target:", host)
	t.registry.MustRegister(&t)
	return &t
}

var (
	PIREnableDesc = prometheus.NewDesc(
		"PIR_enable",
		"Status of PIR",
		nil, nil,
	)
)

func (t *Target) AddPIRStatus(s PIRAlarmStatus) {
	t.Lock()
	defer t.Unlock()
	t.pirStatus = &s
}

func (t *Target) Collect(ch chan<- prometheus.Metric) {
	t.Lock()
	defer t.Unlock()

	if t.pirStatus == nil {
		return
	}

	// Status of PIR
	ch <- prometheus.MustNewConstMetric(
		PIREnableDesc,
		prometheus.GaugeValue,
		float64(t.pirStatus.Enabled),
	)
}

func (t *Target) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(t, ch)
}
