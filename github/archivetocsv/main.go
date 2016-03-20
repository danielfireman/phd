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
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
)

var (
	archivePath = flag.String("archive_path", "", "")
)

type Repo struct {
	ID uint64 `json:"id"`
}

func (r Repo) String() string {
	return fmt.Sprintf("%d", r.ID)
}

type Actor struct {
	ID uint64 `json:"id"`
}

func (a Actor) String() string {
	return fmt.Sprintf("%d", a.ID)
}

type Event struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	CraetedAt string `json:"created_at"`
	Actor     Actor  `json:"actor"`
	Repo      Repo   `json:"repo"`
	Payload   []byte `json:"payload"`
}

type Language struct {
	Name string
	LOC  int
}

type PullRequestEventPayload struct {
	Languages []string
}

func (e Event) String() string {
	if e.Type == "PullRequestEvent" {

	}
	return fmt.Sprintf("%s,%s,%s", e.ID, e.Type, e.CraetedAt)
}

func main() {
	flag.Parse()
	outputPath := filepath.Join(filepath.Dir(*archivePath), strings.Replace(filepath.Base(*archivePath), ".json.gz", ".csv", 1))
	fmt.Printf("Processing %s. Output will be written to %s\n", *archivePath, outputPath)

	reader, err := os.Open(*archivePath)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	events := []Event{}
	scanner := bufio.NewScanner(archive)
	for scanner.Scan() {
		e := Event{}
		if err := json.NewDecoder(bytes.NewBufferString(scanner.Text())).Decode(&e); err != nil {
			log.Fatal(err)
		}
		events = append(events, e)
	}

	o, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Close()

	if err := gocsv.MarshalFile(events, o); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Conversion completed.")
}
