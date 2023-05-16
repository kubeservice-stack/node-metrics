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

package config

import (
	"fmt"
	"strconv"

	"github.com/kubeservice-stack/common/pkg/storage"
)

type Config struct {
	ConcurrentLimit    uint                       `yaml:"concurrent_limit" json:"concurrent_limit"`
	Interval           uint                       `yaml:"interval" json:"interval"`                     //second
	PartitionDuration  uint                       `yaml:"partition_duration" json:"partition_duration"` //second
	Retention          uint                       `yaml:"retention" json:"retention"`                   //second
	TimestampPrecision storage.TimestampPrecision `yaml:"accuracy" json:"accuracy"`                     //emun
}

// default config value
func DefaultConfig() Config {
	return Config{
		ConcurrentLimit:    1,
		Interval:           15,
		PartitionDuration:  5 * 60,
		Retention:          1 * 24 * 3600,
		TimestampPrecision: storage.Seconds,
	}
}

// default config
func (c Config) String() string {
	return fmt.Sprintf(`
[config]
  ## 携程并发限制
  concurrent_limit = %s
  ## 定时执行时间间隔(秒)
  interval = %s
  ## 每个partition的时间间隔(秒)
  partition_duration = %s
  ## 保留整个时间窗口
  retention = %s
  ## timestamp时间进度(Nanoseconds:0; Microseconds:1; Milliseconds:2; Seconds:3)
  accuracy = %s`,
		strconv.FormatUint(uint64(c.ConcurrentLimit), 10),
		strconv.FormatUint(uint64(c.Interval), 10),
		strconv.FormatUint(uint64(c.PartitionDuration), 10),
		strconv.FormatUint(uint64(c.Retention), 10),
		strconv.FormatUint(uint64(int(c.TimestampPrecision)), 10),
	)
}
