package models

import (
	"math"
	"sync/atomic"
)

type CompanyThresholds struct {
	cpuThreshold       atomic.Uint64
	memoryThreshold    atomic.Uint64
	mountThreshold     atomic.Uint64
	diskThreshold      atomic.Uint64
	networkRXThreshold atomic.Uint64
	networkTXThreshold atomic.Uint64
}

// Atomics does not support float64, so we use uint64 with math to convert it.
func floatStore(a *atomic.Uint64, v float64) {
	a.Store(math.Float64bits(v))
}

func floatLoad(a *atomic.Uint64) float64 {
	return math.Float64frombits(a.Load())
}

func (t *CompanyThresholds) SetCPU(v float64) {
	floatStore(&t.cpuThreshold, v)
}

func (t *CompanyThresholds) SetMemory(v float64) {
	floatStore(&t.memoryThreshold, v)
}

func (t *CompanyThresholds) SetMount(v float64) {
	floatStore(&t.mountThreshold, v)
}

func (t *CompanyThresholds) SetDisk(v float64) {
	floatStore(&t.diskThreshold, v)
}

func (t *CompanyThresholds) SetNetworkRX(v float64) {
	floatStore(&t.networkRXThreshold, v)
}

func (t *CompanyThresholds) SetNetworkTX(v float64) {
	floatStore(&t.networkTXThreshold, v)
}

func (t *CompanyThresholds) ExceedsCPU(usage float64) bool {
	return usage > floatLoad(&t.cpuThreshold)
}

func (t *CompanyThresholds) ExceedsMemory(usage float64) bool {
	return usage > floatLoad(&t.memoryThreshold)
}

func (t *CompanyThresholds) ExceedsMount(usage float64) bool {
	return usage > floatLoad(&t.mountThreshold)
}

func (t *CompanyThresholds) ExceedsDisk(usage float64) bool {
	return usage > floatLoad(&t.diskThreshold)
}

func (t *CompanyThresholds) ExceedsNetworkRX(rx float64) bool {
	return rx > floatLoad(&t.networkRXThreshold)
}

func (t *CompanyThresholds) ExceedsNetworkTX(tx float64) bool {
	return tx > floatLoad(&t.networkTXThreshold)
}

func (t *CompanyThresholds) ExceedsNetwork(rx, tx float64) bool {
	return rx > floatLoad(&t.networkRXThreshold) || tx > floatLoad(&t.networkTXThreshold)
}
