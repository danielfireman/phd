package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type Repo struct {
	ID uint64 `json:"id"`
}

type Event struct {
	Type      string `json:"type"`
	CraetedAt string `json:"created_at"`
	Repo      Repo   `json:"repo"`
}

func (e Event) String() string {
	return fmt.Sprintf("%s,%d", e.Type, e.Repo.ID)
}

func main() {
	flag.Parse()
	archive, err := gzip.NewReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	scanner := bufio.NewScanner(archive)
	lines := 0
	for scanner.Scan() {
		e := Event{}
		if err := json.NewDecoder(bytes.NewBufferString(scanner.Text())).Decode(&e); err != nil {
			log.Fatal(err)
		}
		lines++
		t, err := time.Parse(time.RFC3339, e.CraetedAt)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "%s,%d,%d,%d,%d,%d,%d\n", e.String(), t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	}
	fmt.Fprintf(os.Stderr, "%d lines converted.\n", lines)
}
