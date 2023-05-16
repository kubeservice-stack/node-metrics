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

package schedule

import (
	"time"

	sched "github.com/kubeservice-stack/common/pkg/schedule"
	"github.com/kubeservice-stack/common/pkg/storage"
	"github.com/kubeservice-stack/node-metrics/pkg/config"
)

var (
	nodeCPUSecondsSchedule    = sched.NewScheduler()
	nodeMemorySecondsSchedule = sched.NewScheduler()

	nodeCPURawDataStorage storage.StorageInterface
	nodeMemoryDataStorage storage.StorageInterface

	err error
)

func InitDataStorage(cfg config.Config) {
	nodeCPURawDataStorage, err = storage.NewStorage(
		storage.WithPartitionDuration(time.Duration(cfg.PartitionDuration)*time.Second),
		storage.WithTimestampPrecision(cfg.TimestampPrecision),
		storage.WithRetention(time.Duration(cfg.Retention)*time.Second),
	)

	if err != nil {
		panic("init nodeCPURawDataStorage error")
	}

	nodeMemoryDataStorage, err = storage.NewStorage(
		storage.WithPartitionDuration(time.Duration(cfg.PartitionDuration)*time.Second),
		storage.WithTimestampPrecision(cfg.TimestampPrecision),
		storage.WithRetention(time.Duration(cfg.Retention)*time.Second),
	)

	if err != nil {
		panic("init nodeMemoryDataStorage error")
	}
}
