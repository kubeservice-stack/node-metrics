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
	"time"

	"github.com/kubeservice-stack/node-metrics/pkg/schedule"
	"github.com/montanaflynn/stats"
)

func cpuMeanData(num int64) (float64, error) {
	now := time.Now()
	data, err := schedule.NodeCPURawDataStorage().Select("cpu_usage_active", nil, now.Unix()-num, now.Unix())
	if err != nil {
		return 0.0, err
	}

	d := stats.LoadRawData(Values(data))
	return d.Mean()
}

func cpuMaxData(num int64) (float64, error) {
	now := time.Now()
	data, err := schedule.NodeCPURawDataStorage().Select("cpu_usage_active", nil, now.Unix()-num, now.Unix())
	if err != nil {
		return 0.0, err
	}

	d := stats.LoadRawData(Values(data))
	return d.Max()
}
