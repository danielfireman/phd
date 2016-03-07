package load

import (
	"net/http"
	"sync"
	"time"
)

type LoadPattern struct {
	Ops            int `json:"ops"`
	Mem            int `json:"mem"`
	QueryFormatStr string
}

func Run(maxQPS, nWorkers int, duration time.Duration, loadURLs ...string) *Summary {
	startTime := time.Now()
	cLevel := int(float32(nWorkers) / float32(len(loadURLs)))
	qps := int(float32(maxQPS) / float32(len(loadURLs)))
	resultsChan := make(chan RequestResult, nWorkers*1000)
	var wg sync.WaitGroup
	for _, q := range loadURLs {
		wg.Add(1)
		go func(q string) {
			defer wg.Done()
			g := Generator{
				ConcurrencyLevel: cLevel,
				StressURL:        q,
				MaxQPS:           qps,
				Duration:         duration,
			}
			g.Run(resultsChan)
		}(q)
	}
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	return summarize(resultsChan, startTime)
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

type Generator struct {
	ConcurrencyLevel int
	StressURL        string
	MaxQPS           int
	Duration         time.Duration
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
