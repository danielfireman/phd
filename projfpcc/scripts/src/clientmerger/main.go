package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	rounds = flag.Int("rounds", 0, "Number of experiment rounds to be processed.")
	path   = flag.String("path", "", "Directory where the data is located.")
)

type record struct {
	ts         int64
	qps        int32
	throughput float32
}

func main() {
	flag.Parse()

	if *rounds == 0 {
		log.Fatalf("invalid round: %d", *rounds)
	}

	// Reading and consolidating input.
	recordsByQps := make(map[int32][]record)
	for round := 1; round <= *rounds; round++ {
		m1, err := processFile(filepath.Join(*path, fmt.Sprintf("client_%d_1", round)))
		if err != nil {
			log.Fatal(err)
		}
		m2, err := processFile(filepath.Join(*path, fmt.Sprintf("client_%d_2", round)))
		if err != nil {
			log.Fatal(err)
		}
		for k, rec1 := range m1 {
			rec2 := m2[k]
			totalQps := rec1.qps + rec2.qps
			recordsByQps[totalQps] = append(recordsByQps[totalQps], record{
				throughput: rec1.throughput + rec2.throughput,
				ts:         rec1.ts,
				qps:        totalQps,
			})
		}
	}

	// Generating output.
	for qps, records := range recordsByQps {
		outFile, err := os.Create(filepath.Join(*path, fmt.Sprintf("throughput_%d.csv", qps)))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing file: ", outFile.Name())
		w := bufio.NewWriter(outFile)
		fmt.Fprintf(w, "ts,qps,throughput\n")
		for _, rec := range records {
			fmt.Fprintf(w, "%d,%d,%.2f\n", rec.ts, rec.qps, rec.throughput)
		}
		w.Flush()
		outFile.Close()
	}
}

func processFile(f string) (map[int]record, error) {
	fmt.Println("Processing file: ", f)
	csvFile, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	m := make(map[int]record)
	scanner := bufio.NewScanner(csvFile)
	for lineno := 0; scanner.Scan(); lineno++ {
		// Ignoring header
		if lineno == 0 {
			continue
		}
		r := strings.Split(scanner.Text(), ",")
		id, _ := strconv.Atoi(r[0])
		ts, _ := strconv.ParseInt(r[1], 10, 64)
		qps, _ := strconv.ParseInt(r[2], 10, 32)
		throughput, _ := strconv.ParseFloat(r[6], 32)
		m[id] = record{
			ts:         ts,
			qps:        int32(qps),
			throughput: float32(throughput),
		}
	}
	return m, nil
}
