/*
Copyright 2023 The KubeService-Stack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/kubeservice-stack/common/pkg/storage"
	"github.com/kubeservice-stack/node-metrics/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
)

type BaseCollector struct {
	metric map[string]*typedDesc
	logger log.Logger
}

// NewBaseCollector returns a new Collector exposing base average stats.
func NewBaseCollector(logger log.Logger) (Collector, error) {
	return &BaseCollector{
		metric: map[string]*typedDesc{
			"cpu_usage_active":     {prometheus.NewDesc(namespace+"_cpu_usage_active", "cpu usage active.", nil, nil), prometheus.GaugeValue},
			"cpu_usage_max_avg_1h": {prometheus.NewDesc(namespace+"_cpu_usage_max_avg_1h", "cpu usage max avg 1h.", nil, nil), prometheus.GaugeValue},
			"cpu_usage_max_avg_1d": {prometheus.NewDesc(namespace+"_cpu_usage_max_avg_1d", "cpu usage max avg 1d.", nil, nil), prometheus.GaugeValue},
			"cpu_usage_avg_5m":     {prometheus.NewDesc(namespace+"_cpu_usage_avg_5m", "cpu usage avg 5m.", nil, nil), prometheus.GaugeValue},
			"mem_usage_active":     {prometheus.NewDesc(namespace+"_mem_usage_active", "mem usage active.", nil, nil), prometheus.GaugeValue},
			"mem_usage_max_avg_1h": {prometheus.NewDesc(namespace+"_mem_usage_max_avg_1h", "mem usage max avg 1h.", nil, nil), prometheus.GaugeValue},
			"mem_usage_max_avg_1d": {prometheus.NewDesc(namespace+"_mem_usage_max_avg_1d", "mem usage max avg 1d.", nil, nil), prometheus.GaugeValue},
			"mem_usage_avg_5m":     {prometheus.NewDesc(namespace+"_mem_usage_avg_5m", "mem usage avg 5m.", nil, nil), prometheus.GaugeValue},
		},
		logger: logger,
	}, nil
}

func (c *BaseCollector) Update(ch chan<- prometheus.Metric) error {
	loads, err := getData() // 实时指标
	if err != nil {
		return fmt.Errorf("couldn't get load: %w", err)
	}
	for name, load := range loads {
		level.Debug(c.logger).Log("msg", "return load", "index", name, "load", load)
		ch <- c.metric[name].mustNewConstMetric(load)
	}
	return err
}

func Values(data []*storage.DataPoint) []float64 {
	var values []float64
	for _, v := range data {
		values = append(values, v.Value)
	}
	return values
}

func getData() (map[string]float64, error) {
	loads := make(map[string]float64)
	cpu_usage_active, err := cpuMeanData(int64(config.DefaultConfig().Interval) + int64(1))
	if err != nil {
		return loads, nil
	}
	loads["cpu_usage_active"] = cpu_usage_active

	cpu_usage_avg_5m, err := cpuMeanData(301)
	if err != nil {
		return loads, nil
	}
	loads["cpu_usage_avg_5m"] = cpu_usage_avg_5m

	cpu_usage_max_avg_1h, err := cpuMaxData(3601)
	if err != nil {
		return loads, nil
	}
	loads["cpu_usage_max_avg_1h"] = cpu_usage_max_avg_1h

	cpu_usage_max_avg_1d, err := cpuMaxData(3600*24 + 1)
	if err != nil {
		return loads, nil
	}
	loads["cpu_usage_max_avg_1d"] = cpu_usage_max_avg_1d

	mem_usage_active, err := memMeanData(int64(config.DefaultConfig().Interval) + int64(1))
	if err != nil {
		return loads, nil
	}
	loads["mem_usage_active"] = mem_usage_active

	mem_usage_avg_5m, err := memMeanData(301)
	if err != nil {
		return loads, nil
	}
	loads["mem_usage_avg_5m"] = mem_usage_avg_5m

	mem_usage_max_avg_1h, err := memMaxData(3601)
	if err != nil {
		return loads, nil
	}
	loads["mem_usage_max_avg_1h"] = mem_usage_max_avg_1h

	mem_usage_max_avg_1d, err := memMaxData(3600*24 + 1)
	if err != nil {
		return loads, nil
	}
	loads["mem_usage_max_avg_1d"] = mem_usage_max_avg_1d
	return loads, nil
}
