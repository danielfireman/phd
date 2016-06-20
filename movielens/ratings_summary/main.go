package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type stats struct {
	Qtd        float64
	SumRatings float64
}

func main() {
	newStats := make(map[string]*stats)
	r := csv.NewReader(bufio.NewReader(os.Stdin))
	for {
		members, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		// Skipping the title.
		if members[0] == "userId" && members[1] == "movieId" && members[2] == "rating" && members[3] == "timestamp" {
			continue
		}
		movieID := members[1]
		m, ok := newStats[movieID]
		if !ok {
			m = &stats{}
			newStats[movieID] = m
		}
		m.Qtd += 1

		rating, err := strconv.ParseFloat(members[2], 64)
		if err != nil {
			log.Fatalf("Error processing line: %s. Err:%q", members, err)
		}
		m.SumRatings += rating
	}

	bw := bufio.NewWriterSize(os.Stdout, 2048)
	defer bw.Flush()

	w := csv.NewWriter(bw)
	defer w.Flush()
	for k, v := range newStats {
		w.Write([]string{
			fmt.Sprintf("%s", k),
			fmt.Sprintf("%.2f", v.Qtd),
			fmt.Sprintf("%.2f", v.SumRatings),
			fmt.Sprintf("%.2f", v.SumRatings/v.Qtd),
		})
	}
}
