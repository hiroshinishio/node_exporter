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

//go:build !nodiskstats
// +build !nodiskstats

package collector

import (
	"fmt"

	"github.com/go-kit/log"
	"github.com/power-devops/perfstat"
	"github.com/prometheus/client_golang/prometheus"
)

const diskstatsDefaultIgnoredDevices = ""

type diskstatsCollector struct {
	rbytes typedDesc
	wbytes typedDesc
	time   typedDesc

	deviceFilter deviceFilter
	logger       log.Logger
}

func init() {
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}

// NewDiskstatsCollector returns a new Collector exposing disk device stats.
func NewDiskstatsCollector(logger log.Logger) (Collector, error) {
	deviceFilter, err := newDiskstatsDeviceFilter(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse device filter flags: %w", err)
	}

	return &diskstatsCollector{
		rbytes: typedDesc{readBytesDesc, prometheus.CounterValue},
		wbytes: typedDesc{writtenBytesDesc, prometheus.CounterValue},
		time:   typedDesc{ioTimeSecondsDesc, prometheus.CounterValue},

		deviceFilter: deviceFilter,
		logger:       logger,
	}, nil
}

func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	stats, err := perfstat.DiskAdapterStat()
	if err != nil {
		return err
	}

	for _, stat := range stats {
		if c.deviceFilter.ignored(stat.Name) {
			continue
		}
		ch <- c.rbytes.mustNewConstMetric(float64(stat.Rblks/512), stat.Name)
		ch <- c.wbytes.mustNewConstMetric(float64(stat.Wblks/512), stat.Name)
		ch <- c.time.mustNewConstMetric(float64(stat.Time), stat.Name)
	}
	return nil
}
