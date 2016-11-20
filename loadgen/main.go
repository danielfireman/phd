package main

// Examples of simple and quick running (for testing):
//
// To try out keep_duration flag.
// go run main.go --initial_qps=1 --step_duration=1s --step_size=1 --max_qps=2 \
// --num_warmup_steps=1 --timeout=100ms --workers=1 --addr=http://localhost:8080 \
// --cpus=1 --keep_duration=30s

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
	"strings"
	"sync/atomic"
	"time"
)

var (
	initialQps     = flag.Int("initial_qps", 50, "Initial QPS impressed on the server.")
	stepDuration   = flag.Duration("step_duration", 10*time.Second, "Duration of the load step. Example: 1m")
	stepSize       = flag.Int("step_size", 100, "Step size.")
	maxQPS         = flag.Int("max_qps", 1500, "Maximum QPS.")
	numWarmupSteps = flag.Int("num_warmup_steps", 2, "Number of steps to warmup. They are going to receive initial QPS.")
	timeout        = flag.Duration("timeout", 20*time.Millisecond, "HTTP client timeout")
	clientAddr     = flag.String("addr", "http://10.4.2.103:8080", "Client HTTP address")
	numWorkers     = flag.Int("workers", 32, "Client HTTP address")
	numCores       = flag.Int("cpus", 2, "Client HTTP address")
	msgSuffixes    = flag.String("msg_suffixes", "/numprimes/5000", "Suffix to add to the msg.")
	keepDuration   = flag.Duration("keep_duration", 0*time.Millisecond, "Time without increasing QPS after max.")
)

const (
	quitSuffix = "/quit"
)

// Shared variables, need to go trough atomic.
var succs, errs, reqs uint64

func main() {
	flag.Parse()

	suffixes := strings.Split(*msgSuffixes, ",")

	if *initialQps <= 0 {
		log.Fatalf("InitialQps must be positive.")
	}
	runtime.GOMAXPROCS(*numCores)
	workers := int(*numWorkers)
	fmt.Fprintf(os.Stderr, "RunningOn:%d Workers:%d", runtime.GOMAXPROCS(0), workers)

	pauseChan := make(chan struct{})
	work := make(chan struct{}, *maxQPS*workers)
	for i := 0; i < workers; i++ {
		client := http.Client{
			Timeout: *timeout,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   *timeout,
					KeepAlive: *timeout,
				}).Dial,
			},
		}
		go worker(client, work, pauseChan, suffixes)
	}

	qps := *initialQps
	fmt.Printf("ts,qps,totalReq,succReq,errReq,throughput\n")
	numSteps := 0
	var increaseLoadFinishTime time.Time
	for {
		step := time.Tick(*stepDuration)
		t := time.Tick(time.Duration(float64(1e9)/float64(qps)) * time.Nanosecond)
		start := time.Now()
		func() {
			for {
				<-t
				select {
				case <-pauseChan:
					dur := time.Now().Sub(start)
					close(pauseChan)
					time.Sleep(*stepDuration)
					pauseChan = make(chan struct{})
					fmt.Printf("%d,%d,%d,%d,%d,%.2f\n", time.Now().Unix(), qps, reqs, succs, errs, float64(succs)/dur.Seconds())
					atomic.StoreUint64(&reqs, 0)
					atomic.StoreUint64(&succs, 0)
					atomic.StoreUint64(&errs, 0)
					return
				case <-step:
					dur := time.Now().Sub(start)
					fmt.Printf("%d,%d,%d,%d,%d,%.2f\n", time.Now().Unix(), qps, reqs, succs, errs, float64(succs)/dur.Seconds())
					atomic.StoreUint64(&reqs, 0)
					atomic.StoreUint64(&succs, 0)
					atomic.StoreUint64(&errs, 0)
					return
				default:
					// Kind of flow control. We don't want the queue grows too big.
					if len(work) < workers*2 {
						work <- struct{}{}
					}
				}
			}
		}()

		if qps >= *maxQPS {
			if increaseLoadFinishTime.Nanosecond() == 0 {
				increaseLoadFinishTime = time.Now()
			}
			if increaseLoadFinishTime.Add(*keepDuration).Before(time.Now()) {
				close(work)
				resp, err := http.Get(*clientAddr + quitSuffix)
				if err == nil {
					io.Copy(ioutil.Discard, resp.Body)
				}
				return
			}
		} else {
			numSteps++
			if numSteps >= *numWarmupSteps {
				qps += *stepSize
			}
		}
	}
}

func worker(client http.Client, work chan struct{}, pauseChan chan struct{}, suffixes []string) {
	for {
		for _, suffix := range suffixes {
			<-work
			atomic.AddUint64(&reqs, 1)
			url := *clientAddr + suffix
			resp, err := client.Get(url)
			if err == nil {
				switch resp.StatusCode {
				case http.StatusOK:
					atomic.AddUint64(&succs, 1)
				case http.StatusTooManyRequests:
					pauseChan <- struct{}{}
				default:
					atomic.AddUint64(&errs, 1)
				}
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			} else {
				atomic.AddUint64(&errs, 1)
			}
		}
	}
}
