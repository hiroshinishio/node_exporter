package collector

import (
	"github.com/power-devops/perfstat"
)

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	stats, err := perfstat.MemoryTotalStat()
	if err != nil {
		return nil, err
	}

	return map[string]float64{
		"total": float64(stats.RealTotal),
		"free":  float64(stats.RealFree),
	}, nil
}
