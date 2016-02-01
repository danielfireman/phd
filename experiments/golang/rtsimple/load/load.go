package load

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type result struct {
	err      error
	code     int
	duration time.Duration
}

type Generator struct {
	ConcurrencyLevel int
	StressURL        string
	ExpVarsURL       string
	MaxQPS           int
	Duration         time.Duration
	results          chan result
}

type IntMetric struct {
	Max        int64
	Min        int64
	NumSamples int64
	Avg        float64
	sum        int64
}

func (m *IntMetric) sample(s int64) {
	m.NumSamples++
	m.sum += s
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

type Results struct {
	LatencyNanos IntMetric
	Duration     time.Duration
	NumErrors    int64
	QPS          float64
	ServerInfo   ServerInfo
}

func req(c http.Client, url string) result {
	var code int
	s := time.Now()
	resp, err := c.Get(url)
	if err == nil {
		resp.Body.Close()
		code = resp.StatusCode
	}
	return result{
		code:     code,
		duration: time.Now().Sub(s),
		err:      err,
	}
}

func (g *Generator) Run() *Results {
	log.Println("Starting")
	g.results = make(chan result, g.ConcurrencyLevel*1000)
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(g.ConcurrencyLevel)

	// Triggering generators.
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	for i := 0; i < g.ConcurrencyLevel; i++ {
		share := float64(g.MaxQPS) / float64(g.ConcurrencyLevel)
		go func() {
			f := time.Tick(g.Duration)
			t := time.Tick(time.Duration(1e6/(share)) * time.Microsecond)
			for {
				select {
				case <-t:
					g.results <- req(client, g.StressURL)
				case <-f:
					wg.Done()
					return
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(g.results)
	}()

	// Building results.
	results := &Results{}
	for r := range g.results {
		switch {
		case r.err != nil:
			results.NumErrors++
		default:
			results.LatencyNanos.sample(r.duration.Nanoseconds())
		}
	}
	results.Duration = time.Now().Sub(start)
	results.QPS = float64(results.NumErrors+results.LatencyNanos.NumSamples) / results.Duration.Seconds()
	results.LatencyNanos.Avg = float64(results.LatencyNanos.sum) / float64(results.LatencyNanos.NumSamples)

	// Fetching RT information from the server.
	if g.ExpVarsURL == "" {
		return results
	}
	resp, err := client.Get(g.ExpVarsURL)
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
