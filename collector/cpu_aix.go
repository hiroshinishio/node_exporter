// Copyright 2024 The Prometheus Authors
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
	"strconv"

	"github.com/go-kit/log"
	"github.com/power-devops/perfstat"
	"github.com/prometheus/client_golang/prometheus"
)

type cpuCollector struct {
	cpu    typedDesc
	logger log.Logger
}

func init() {
	registerCollector("cpu", defaultEnabled, NewCpuCollector)
}

func NewCpuCollector(logger log.Logger) (Collector, error) {
	return &cpuCollector{
		cpu:    typedDesc{nodeCPUSecondsDesc, prometheus.CounterValue},
		logger: logger,
	}, nil
}

func (c *cpuCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := perfstat.CpuStat()
	if err != nil {
		return err
	}

	for n, stat := range stats {
		ch <- c.cpu.mustNewConstMetric(float64(stat.User), strconv.Itoa(n), "user")
		ch <- c.cpu.mustNewConstMetric(float64(stat.Sys), strconv.Itoa(n), "system")
		ch <- c.cpu.mustNewConstMetric(float64(stat.Idle), strconv.Itoa(n), "idle")
		ch <- c.cpu.mustNewConstMetric(float64(stat.Wait), strconv.Itoa(n), "wait")
	}
	return nil
}
