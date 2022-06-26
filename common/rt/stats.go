package rt

import "runtime"

// Stats ...
type Stats struct {
	Allocated  uint64
	Objects    uint64
	InUse      uint64
	GCTime     uint64
	PauseTotal uint64
}

// GetMemStats ...
func GetMemStats(runGC bool) Stats {
	if runGC {
		runtime.GC()
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	return Stats{
		Allocated:  ms.Alloc,
		Objects:    ms.Mallocs - ms.Frees,
		InUse:      ms.HeapInuse,
		GCTime:     ms.PauseNs[(ms.NumGC+255)%256],
		PauseTotal: ms.PauseTotalNs,
	}
}
