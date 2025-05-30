// Copyright 2018 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !nocpu
// +build !nocpu

package collector

import (
	"log/slog"
	"strconv"
	"unsafe"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

type clockinfo struct {
	hz      int32
	tick    int32
	tickadj int32
	stathz  int32
	profhz  int32
}

const (
	CP_USER = iota
	CP_NICE
	CP_SYS
	CP_SPIN
	CP_INTR
	CP_IDLE
	CPUSTATES
)
const (
	CP_USER_O63 = iota
	CP_NICE_O63
	CP_SYS_O63
	CP_INTR_O63
	CP_IDLE_O63
	CPUSTATES_O63
)

type cpuCollector struct {
	cpu    typedDesc
	logger *slog.Logger
}

func init() {
	registerCollector("cpu", defaultEnabled, NewCPUCollector)
}

func NewCPUCollector(logger *slog.Logger) (Collector, error) {
	return &cpuCollector{
		cpu:    typedDesc{nodeCPUSecondsDesc, prometheus.CounterValue},
		logger: logger,
	}, nil
}

func (c *cpuCollector) Update(ch chan<- prometheus.Metric) (err error) {
	clockb, err := unix.SysctlRaw("kern.clockrate")
	if err != nil {
		return err
	}
	clock := *(*clockinfo)(unsafe.Pointer(&clockb[0]))
	hz := float64(clock.stathz)

	ncpus, err := unix.SysctlUint32("hw.ncpu")
	if err != nil {
		return err
	}

	var cpTime [][CPUSTATES]uint64
	for i := 0; i < int(ncpus); i++ {
		cpb, err := unix.SysctlRaw("kern.cp_time2", i)
		if err != nil && err != unix.ENODEV {
			return err
		}
		if err != unix.ENODEV {
			var times [CPUSTATES]uint64
			for n := 0; n < len(cpb); n += 8 {
				times[n/8] = *(*uint64)(unsafe.Pointer(&cpb[n]))
			}
			if len(cpb)/8 == CPUSTATES_O63 {
				copy(times[CP_INTR:], times[CP_INTR_O63:])
				times[CP_SPIN] = 0
			}
			cpTime = append(cpTime, times)
		}
	}

	for cpu, time := range cpTime {
		lcpu := strconv.Itoa(cpu)
		ch <- c.cpu.mustNewConstMetric(float64(time[CP_USER])/hz, lcpu, "user")
		ch <- c.cpu.mustNewConstMetric(float64(time[CP_NICE])/hz, lcpu, "nice")
		ch <- c.cpu.mustNewConstMetric(float64(time[CP_SYS])/hz, lcpu, "system")
		ch <- c.cpu.mustNewConstMetric(float64(time[CP_SPIN])/hz, lcpu, "spin")
		ch <- c.cpu.mustNewConstMetric(float64(time[CP_INTR])/hz, lcpu, "interrupt")
		ch <- c.cpu.mustNewConstMetric(float64(time[CP_IDLE])/hz, lcpu, "idle")
	}
	return err
}
