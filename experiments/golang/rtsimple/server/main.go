package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

var port = flag.Int("port", 8999, "Port to run/expose service. Example: '8999'")

func cpu(w http.ResponseWriter, r *http.Request) {
	mem := getIntParam(r, "mem")
	cpu := getIntParam(r, "cpu")
	m := make([]byte, mem, mem)
	for i := 0; i < cpu; i++ {
		m[0] = byte(math.Sqrt(float64(i)))
	}
	w.WriteHeader(http.StatusOK)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func quit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received command to quit.")
	os.Exit(1)
}

func main() {
	flag.Parse()
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
