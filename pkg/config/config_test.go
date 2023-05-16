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
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config(t *testing.T) {
	assert := assert.New(t)
	cfg := DefaultConfig()
	assert.Equal(cfg, Config{ConcurrentLimit: 0x1, Interval: 0xf, PartitionDuration: 0x12c, Retention: 0x15180, TimestampPrecision: 3})

	str := cfg.String()
	assert.Equal(str, "\n[config]\n  ## 携程并发限制\n  concurrent_limit = 1\n  ## 定时执行时间间隔(秒)\n  interval = 15\n  ## 每个partition的时间间隔(秒)\n  partition_duration = 300\n  ## 保留整个时间窗口\n  retention = 86400\n  ## timestamp时间进度(Nanoseconds:0; Microseconds:1; Milliseconds:2; Seconds:3)\n  accuracy = 3")
}
