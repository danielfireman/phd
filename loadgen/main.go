package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

var (
	initialQps   = flag.Int("initial_qps", 50, "Initial QPS impressed on the server.")
	stepDuration = flag.Duration("step_duration", 10*time.Second, "Duration of the load step. Example: 1m")
	stepSize     = flag.Int("step_size", 50, "Step size.")
	maxQPS       = flag.Int("max_qps", 1500, "Maximum QPS.")
	timeout      = flag.Duration("timeout", 20*time.Millisecond, "HTTP client timeout")
	clientAddr   = flag.String("addr", "http://10.4.2.103:8080", "Client HTTP address")
)

const (
	msgSuffix  = "/msg"
	quitSuffix = "/quit"
)

// Shared variables, need to go trough atomic.
var succs, errs, reqs uint64

func main() {
	if *initialQps <= 0 {
		log.Fatalf("InitialQps must be positive.")
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	workers := int(32)
	fmt.Fprintf(os.Stderr, "RunningOn:%d Workers:%d", runtime.GOMAXPROCS(0), workers)

	work := make(chan struct{}, *maxQPS*workers)
	for i := 0; i < workers; i++ {
		client := http.Client{
			Timeout: *timeout,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   *timeout,
					KeepAlive: *timeout,
				}).Dial,
				DisableKeepAlives: true,
			},
		}
		go worker(client, work)
	}

	qps := *initialQps
	fmt.Printf("qps,totalReq,succReq,errReq,throughput\n")
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
		fmt.Printf("%d,%d,%d,%d,%.2f\n", qps, reqs, succs, errs, float64(succs)/(*stepDuration).Seconds())

		if float64(reqs) < 0.80*float64(qps)*stepDuration.Seconds() {
			log.Fatalf("Client can not handle the load")
		}

		atomic.StoreUint64(&reqs, 0)
		atomic.StoreUint64(&succs, 0)
		atomic.StoreUint64(&errs, 0)

		if qps >= *maxQPS {
			close(work)
			resp, err := http.Get(*clientAddr + quitSuffix)
			if err == nil {
				io.Copy(ioutil.Discard, resp.Body)
			}
			return
		}
		qps += *stepSize
	}
}

func worker(client http.Client, work chan struct{}) {
	url := *clientAddr + msgSuffix
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
