package load

import (
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
	N                int
	ConcurrencyLevel int
	URL              string
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
	} else if s < m.Min {
		m.Min = s
	}
}

type Report struct {
	LatencyNanos IntMetric
	Duration     time.Duration
	NumErrors    int
	QPS          float64
}

func req(url string) result {
	var code int
	s := time.Now()
	resp, err := http.Get(url)
	if err == nil {
		code = resp.StatusCode
	}
	return result{
		code:     code,
		duration: time.Now().Sub(s),
		err:      err,
	}
}

func (g *Generator) Run() *Report {
	g.results = make(chan result, g.N)
	start := time.Now()

	wg := sync.WaitGroup{}
	c := make(chan struct{})

	// Kicking off the actual load generators.
	for i := 0; i < g.ConcurrencyLevel; i++ {
		go func() {
			for range c {
				wg.Add(1)
				g.results <- req(g.URL)
				wg.Done()
			}
		}()
	}
	for i := 0; i < g.N; i++ {
		c <- struct{}{}
	}
	// Closing token channel and waiting for the load generation to be finish.
	close(c)
	wg.Wait()

	// Building report.
	d := time.Now().Sub(start)
	report := &Report{
		Duration:     d,
		LatencyNanos: IntMetric{},
		QPS:          d.Seconds() / float64(g.N),
	}
	for r := range g.results {
		switch {
		case r.err != nil:
			report.NumErrors++
		default:
			report.LatencyNanos.sample(r.duration.Nanoseconds())
		}
	}
	return report
}
