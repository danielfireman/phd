package load

import (
	"log"
	"net/http"
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
	URL              string
	MaxQPS           int
	Duration         time.Duration
	results          chan result
}

type IntMetric struct {
	Max        int64
	Min        int64
	NumSamples int64
	sum        int64
}

func (m *IntMetric) Avg() float64 {
	return float64(m.sum) / float64(m.NumSamples)
}

func (m *IntMetric) sample(s int64) {
	m.NumSamples++
	if s > m.Max {
		m.Max = s
	} else if s < m.Min || m.Min == 0 {
		m.Min = s
	}
}

type Report struct {
	LatencyNanos IntMetric
	Duration     time.Duration
	NumErrors    int64
	QPS          float64
}

func req(c http.Client, url string) result {
	var code int
	s := time.Now()
	resp, err := c.Get(url)
	if err == nil {
		defer resp.Body.Close()
		code = resp.StatusCode
	}
	return result{
		code:     code,
		duration: time.Now().Sub(s),
		err:      err,
	}
}

func (g *Generator) Run() *Report {
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
					g.results <- req(client, g.URL)
				case <-f:
					wg.Done()
					return
				}
			}
		}()
	}
	wg.Wait()
	close(g.results)

	// Building report.
	d := time.Now().Sub(start)
	report := &Report{
		Duration:     d,
		LatencyNanos: IntMetric{},
	}
	for r := range g.results {
		switch {
		case r.err != nil:
			report.NumErrors++
		default:
			report.LatencyNanos.sample(r.duration.Nanoseconds())
		}
	}
	report.QPS = float64(report.NumErrors+report.LatencyNanos.NumSamples) / d.Seconds()
	return report
}
