package serverinfo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"
)

type GC struct {
	Num         uint32
	PauseTotal  time.Duration
	CPUFraction float64
	Enabled     bool
	Debug       bool
}

type ServerInfo struct {
	GC GC
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
	return &ServerInfo{
		GC: GC{
			CPUFraction: v.MemStats.GCCPUFraction,
			Num:         v.MemStats.NumGC,
			PauseTotal:  time.Duration(v.MemStats.PauseTotalNs) * time.Nanosecond,
			Enabled:     v.MemStats.EnableGC,
			Debug:       v.MemStats.DebugGC,
		},
	}, nil
}
