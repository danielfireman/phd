package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"runtime"
	"strconv"

	_ "expvar"
)

var (
	port       = flag.Int("port", 8999, "Port to run/expose service. Example: 8999")
	gomaxprocs = flag.Int("gomaxprocs", 2, "Number of max processors used by golang runtime. Example: 2")
)

func cpu(w http.ResponseWriter, r *http.Request) {
	m := allocateMem(getIntParam(r, "mem"))
	for i := 0; i < getIntParam(r, "cpu"); i++ {
		math.Sqrt(float64(m[0]))
	}
}

func allocateMem(b int) []byte {
	return make([]byte, b, b)
}

func ping(w http.ResponseWriter, r *http.Request) {
}

func quit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received command to quit.")
	os.Exit(1)
}

func main() {
	flag.Parse()

	runtime.GOMAXPROCS(*gomaxprocs)
	fmt.Printf("Listening at http://localhost:%d\n", *port)
	http.HandleFunc("/work", cpu)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/quit", quit)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

func getIntParam(r *http.Request, p string) int {
	f, err := strconv.Atoi(r.FormValue(p))
	if err != nil {
		f = 0
	}
	return f
}
