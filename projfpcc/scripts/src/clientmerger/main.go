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
	rounds     = flag.Int("rounds", 0, "Number of experiment rounds to be processed.")
	path       = flag.String("path", "", "Directory where the data is located.")
	numClients = flag.Int("clients", 1, "Number of client processes.")
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
		roundRecords := make(map[int32]record, *numClients)
		m1, err := processFile(filepath.Join(*path, fmt.Sprintf("client_%d_1", round)))
		if err != nil {
			log.Fatal(err)
		}
		for _, rec := range m1 {
			roundRecords[rec.qps] = rec
		}
		for client := 2; client <= *numClients; client++ {
			mCurr, err := processFile(filepath.Join(*path, fmt.Sprintf("client_%d_%d", round, client)))
			if err != nil {
				log.Fatal(err)
			}
			for qps, rec := range mCurr {
				consolidatedRec := roundRecords[qps]
				roundRecords[qps] = record{
					qps:        consolidatedRec.qps + rec.qps,
					ts:         consolidatedRec.ts,
					throughput: consolidatedRec.throughput + rec.throughput,
				}
			}
		}
		for qps, rec := range roundRecords {
			recordsByQps[qps] = append(recordsByQps[qps], rec)
		}
	}

	outFile, err := os.Create(filepath.Join(*path, "throughput.csv"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing file: ", outFile.Name())
	w := bufio.NewWriter(outFile)
	fmt.Fprintf(w, "load,ts,throughput\n")
	// Generating output.
	for _, records := range recordsByQps {
		for _, rec := range records {
			fmt.Fprintf(w, "%d,%d,%.2f\n", rec.qps, rec.ts, rec.throughput)
		}
	}
	w.Flush()
	outFile.Close()
}

func processFile(f string) (map[int32]record, error) {
	fmt.Println("Processing file: ", f)
	csvFile, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	m := make(map[int32]record)
	scanner := bufio.NewScanner(csvFile)
	for lineno := 0; scanner.Scan(); lineno++ {
		// Ignoring header
		if lineno == 0 {
			continue
		}
		r := strings.Split(scanner.Text(), ",")
		ts, _ := strconv.ParseInt(r[0], 10, 64)
		qps, _ := strconv.ParseInt(r[1], 10, 32)
		throughput, _ := strconv.ParseFloat(r[5], 32)
		m[int32(qps)] = record{
			ts:         ts,
			qps:        int32(qps),
			throughput: float32(throughput),
		}
	}
	return m, nil
}
