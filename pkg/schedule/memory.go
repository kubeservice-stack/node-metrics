//go:build !nomeminfo
// +build !nomeminfo

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
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kubeservice-stack/common/pkg/storage"
	"github.com/kubeservice-stack/node-metrics/pkg/config"
	"github.com/kubeservice-stack/node-metrics/pkg/util"
)

type Item struct {
	key   string
	value float64
}

func MemoryInfo() (float64, error) {
	file, err := os.Open(util.ProcFilePath("meminfo"))
	if err != nil {
		return 0.0, err
	}
	defer file.Close()

	var (
		memInfo = map[string]float64{}
		scanner = bufio.NewScanner(file)
	)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		// Workaround for empty lines occasionally occur in CentOS 6.2 kernel 3.10.90.
		if len(parts) == 0 {
			continue
		}
		key := parts[0][:len(parts[0])-1] // remove trailing : from key
		if key == "MemAvailable" || key == "MemTotal" {
			fv, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return 0.0, fmt.Errorf("invalid value in meminfo: %w", err)
			}
			switch len(parts) {
			case 2: // no unit
			case 3: // has unit, we presume kB
				fv *= 1024
			default:
				return 0.0, fmt.Errorf("invalid line in meminfo: %s", line)
			}
			memInfo[key] = fv
		} else {
			continue
		}
	}
	v, ok1 := memInfo["MemTotal"]
	s, ok2 := memInfo["MemAvailable"]
	if ok1 && ok2 && v > 0.0 {
		return 100 * (1 - s/v), nil
	}
	return 0.0, fmt.Errorf("invalid line in meminfo: %v", memInfo)
}

func StorageMemoryPoint() {
	t := time.Now().Unix()
	v, err := MemoryInfo()
	if err == nil {
		nodeMemoryDataStorage.InsertRows(
			[]storage.Row{
				storage.Row{
					Name:      "mem_usage_active",
					DataPoint: storage.DataPoint{Timestamp: t, Value: v},
				},
			})
	}
}

func StartMemory(cfg config.Config) {
	nodeMemorySecondsSchedule.Every(uint64(cfg.Interval)).Second().DoSafely(StorageMemoryPoint)

	go func() {
		<-nodeMemorySecondsSchedule.Start()
	}()
}
