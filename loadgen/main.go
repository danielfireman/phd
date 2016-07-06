package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const url = "http://localhost:8080/msg"

var (
	initialQps   = flag.Int("initial_qps", 2000, "Initial QPS impressed on the server.")
	stepDuration = flag.Duration("step_duration", 10*time.Second, "Duration of the load step. Example: 1m")
	stepSize     = flag.Int("step_size", 200, "Step size.")
	maxQPS       = flag.Int("max_qps", 10000, "Maximum QPS.")
)

// Shared variables, need to go trough atomic.
var succs, errs, reqs uint64

func main() {
	if *initialQps <= 0 {
		log.Fatalf("InitialQps must be positive.")
	}

	client := http.Client{
		Timeout: 100 * time.Millisecond,
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			TLSHandshakeTimeout: 100 * time.Millisecond,
		},
	}
	workers := int(*maxQPS / 2)
	work := make(chan struct{}, *maxQPS*workers)
	for i := 0; i < workers; i++ {
		go worker(client, work)
	}

	qps := *initialQps
	for {
		step := time.Tick(*stepDuration)

		t := time.Tick(time.Duration(float64(1e9)/float64(qps)) * time.Nanosecond)
		func() {
			for {
				<-t
				select {
				case <-step:
					return
				default:
					// Kind of flow control. We don't want the queue grows too big.
					if len(work) < workers*2 {
						work <- struct{}{}
					}
				}
			}
		}()
		fmt.Printf("qps:%d nw:%d req:%d succ:%d err:%d\n", qps, workers, reqs, succs, errs)

		if float64(reqs) < 0.95*float64(qps)*stepDuration.Seconds() {
			log.Fatalf("Client can not handle the load")
		}

		oldSucc := succs
		oldReq := reqs

		atomic.StoreUint64(&reqs, 0)
		atomic.StoreUint64(&succs, 0)
		atomic.StoreUint64(&errs, 0)

		if qps >= *maxQPS {
			close(work)
			return
		}

		if float64(oldSucc) >= 0.95*float64(oldReq) {
			qps += *stepSize
		}
	}
}

func worker(client http.Client, work chan struct{}) {
	for {
		<-work
		atomic.AddUint64(&reqs, 1)
		resp, err := client.Get(url)
		if err == nil {
			if resp.StatusCode == 200 {
				atomic.AddUint64(&succs, 1)
			} else {
				atomic.AddUint64(&errs, 1)
			}
			io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		} else {
			atomic.AddUint64(&errs, 1)
		}
	}
}
