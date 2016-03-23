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
)

type Repo struct {
	ID uint64 `json:"id"`
}

type Actor struct {
	ID uint64 `json:"id"`
}

type Event struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	CraetedAt string `json:"created_at"`
	Actor     Actor  `json:"actor"`
	Repo      Repo   `json:"repo"`
}

func (e Event) String() string {
	return fmt.Sprintf("%s,%s,%s,%d,%d", e.ID, e.Type, e.CraetedAt, e.Actor.ID, e.Repo.ID)
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
		fmt.Fprintf(w, "%s\n", e.String())
	}
	fmt.Fprintf(os.Stderr, "%d lines converted.\n", lines)
}
