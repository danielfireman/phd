package results

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/danielfireman/phd/experiments/golang/rtsimple/load"
)

type DurationMetric struct {
	Max        time.Duration
	Min        time.Duration
	NumSamples int64
	Avg        time.Duration
	Sum        time.Duration
}

func (m *DurationMetric) sample(s time.Duration) {
	m.NumSamples++
	m.Sum += s
	if s > m.Max {
		m.Max = s
	} else if s < m.Min || m.Min == 0 {
		m.Min = s
	}
}

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

type Summary struct {
	Latency    DurationMetric
	Duration   time.Duration
	NumErrors  int64
	QPS        float64
	ServerInfo ServerInfo
}

func Summarize(reqs <-chan load.RequestResult, varsURL string, start time.Time) *Summary {
	results := &Summary{}
	for r := range reqs {
		switch {
		case r.Err != nil:
			results.NumErrors++
		default:
			results.Latency.sample(r.Duration)
		}
	}
	results.Duration = time.Now().Sub(start)
	results.QPS = float64(results.NumErrors+results.Latency.NumSamples) / results.Duration.Seconds()
	results.Latency.Avg = time.Duration(float64(results.Latency.Sum) / float64(results.Latency.NumSamples))

	// Fetching RT information from the server.
	if varsURL == "" {
		return results
	}
	resp, err := http.Get(varsURL)
	if err != nil {
		log.Fatal(err)
	}
	v := &struct {
		MemStats runtime.MemStats `"json":"memstats"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		log.Fatal(err)
	}
	results.ServerInfo.GC = GC{
		CPUFraction: v.MemStats.GCCPUFraction,
		Num:         v.MemStats.NumGC,
		PauseTotal:  time.Duration(v.MemStats.PauseTotalNs) * time.Nanosecond,
		Enabled:     v.MemStats.EnableGC,
		Debug:       v.MemStats.DebugGC,
	}
	return results
}
