//go:build !nocpu
// +build !nocpu

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
	"fmt"
	"time"

	"github.com/kubeservice-stack/common/pkg/storage"
	"github.com/kubeservice-stack/node-metrics/pkg/config"
	"github.com/kubeservice-stack/node-metrics/pkg/util"
	"github.com/prometheus/procfs"
)

func CPUInfo(oldTotal, oldIdle float64) (float64, float64, float64, error) {
	fs, err := procfs.NewFS(*util.ProcPath)
	if err != nil {
		return 0.0, oldTotal, oldIdle, fmt.Errorf("failed to open procfs: %w", err)
	}

	stat, err := fs.Stat()
	if err != nil {
		return 0.0, oldTotal, oldIdle, fmt.Errorf("failed to get /proc/stat: %w", err)
	}
	newTotal := stat.CPUTotal.Guest + stat.CPUTotal.GuestNice + stat.CPUTotal.Idle +
		stat.CPUTotal.Iowait + stat.CPUTotal.IRQ + stat.CPUTotal.Nice + stat.CPUTotal.SoftIRQ +
		stat.CPUTotal.Steal + stat.CPUTotal.Steal + stat.CPUTotal.System + stat.CPUTotal.User
	total := newTotal - oldTotal
	idle := stat.CPUTotal.Idle - oldIdle
	if total > 0.0 && idle > 0.0 {
		return (1 - idle/total) * 100, newTotal, stat.CPUTotal.Idle, nil
	}

	return 0.0, newTotal, stat.CPUTotal.Idle, fmt.Errorf("failed to get cpu, all cpu total == 0 ")
}

var (
	OldTotal float64 = 0.0
	OldIdle  float64 = 0.0
)

func StorageCPUPoint() {
	t := time.Now().Unix()
	v, oldTotal, oldIdle, err := CPUInfo(OldTotal, OldIdle)
	OldTotal = oldTotal
	OldIdle = oldIdle
	if err == nil {
		nodeCPURawDataStorage.InsertRows(
			[]storage.Row{
				storage.Row{
					Name:      "cpu_usage_active",
					DataPoint: storage.DataPoint{Timestamp: t, Value: v},
				},
			})
	}
}

func StartCPU(cfg config.Config) {
	nodeCPUSecondsSchedule.Every(uint64(cfg.Interval)).Second().DoSafely(StorageCPUPoint)

	go func() {
		<-nodeCPUSecondsSchedule.Start()
	}()
}
