package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type stats struct {
	Qtd int
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
		if members[0] == "userId" && members[1] == "movieId" && members[2] == "tag" && members[3] == "timestamp" {
			continue
		}
		movieID := members[1]
		m, ok := newStats[movieID]
		if !ok {
			m = &stats{}
			newStats[movieID] = m
		}
		m.Qtd += 1
	}

	bw := bufio.NewWriterSize(os.Stdout, 2048)
	defer bw.Flush()

	w := csv.NewWriter(bw)
	defer w.Flush()
	for k, v := range newStats {
		w.Write([]string{
			fmt.Sprintf("%s", k),
			fmt.Sprintf("%d", v.Qtd),
		})
	}
}
