package load

import (
	"net/http"
	"sync"
	"time"
)

type Generator struct {
	ConcurrencyLevel int
	StressURL        string
	MaxQPS           int
	Duration         time.Duration
}

type RequestResult struct {
	Err      error
	Code     int
	Duration time.Duration
}

func req(c http.Client, url string) RequestResult {
	var code int
	s := time.Now()
	resp, err := c.Get(url)
	if err == nil {
		resp.Body.Close()
		code = resp.StatusCode
	}
	return RequestResult{
		Code:     code,
		Duration: time.Now().Sub(s),
		Err:      err,
	}
}

func (g *Generator) Run(results chan<- RequestResult) {
	var wg sync.WaitGroup
	wg.Add(g.ConcurrencyLevel)
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	for i := 0; i < g.ConcurrencyLevel; i++ {
		share := float64(g.MaxQPS) / float64(g.ConcurrencyLevel)
		go func() {
			defer wg.Done()
			f := time.Tick(g.Duration)
			t := time.Tick(time.Duration(1e6/(share)) * time.Microsecond)
			for {
				select {
				case <-t:
					results <- req(client, g.StressURL)
				case <-f:
					return
				}
			}
		}()
	}
	wg.Wait()
}
