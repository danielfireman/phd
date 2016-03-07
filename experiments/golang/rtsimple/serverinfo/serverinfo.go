package serverinfo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/pivotal-golang/bytefmt"
)

type GC struct {
	Num           uint32
	PauseTotal    time.Duration
	PauseTotalStr string
	CPUFraction   float64
	Enabled       bool
	Debug         bool
}

type Mem struct {
	TotalAlloc    uint64 // bytes allocated (even if freed)
	TotalAllocStr string
}

type ServerInfo struct {
	GC  GC
	Mem Mem
}

func Fetch(url string) (*ServerInfo, error) {
	if url == "" {
		return nil, fmt.Errorf("Error fetching server info: url can not be empty.")
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	v := &struct {
		MemStats runtime.MemStats `"json":"memstats"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		log.Fatalln(err.Error())
	}

	p := time.Duration(v.MemStats.PauseTotalNs) * time.Nanosecond
	return &ServerInfo{
		GC: GC{
			CPUFraction:   v.MemStats.GCCPUFraction,
			Num:           v.MemStats.NumGC,
			PauseTotal:    p,
			PauseTotalStr: p.String(),
			Enabled:       v.MemStats.EnableGC,
			Debug:         v.MemStats.DebugGC,
		},
		Mem: Mem{
			TotalAlloc:    v.MemStats.TotalAlloc,
			TotalAllocStr: bytefmt.ByteSize(v.MemStats.TotalAlloc),
		},
	}, nil
}
