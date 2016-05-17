package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

var genreMap = map[string]int{
	"Action":      0,
	"Adventure":   1,
	"Animation":   2,
	"Children's":  3,
	"Comedy":      4,
	"Crime":       5,
	"Documentary": 6,
	"Drama":       7,
	"Fantasy":     8,
	"Film-Noir":   9,
	"Horror":      10,
	"Musical":     11,
	"Mystery":     12,
	"Romance":     13,
	"Sci-Fi":      14,
	"Thriller":    15,
	"War":         16,
	"Western":     17,
}

var yearRegexp = regexp.MustCompile(`\(\d+\)`)

var csvReader = csv.NewReader

func main() {
	records, err := csvReader(os.Stdin).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, members := range records {
		// Very simple way to skip the title.
		if members[0] == "movieId" && members[1] == "title" && members[2] == "genres" {
			continue
		}
		fmt.Println(processRow(members))
	}
}

func processRow(members []string) string {
	// id, title, year, #genres
	nFixedFields := 4
	newRow := make([]string, nFixedFields+len(genreMap))
	newRow[0] = members[0]

	title := members[1]
	// Making sure title is treated as string
	newRow[1] = fmt.Sprintf("\"%s\"", title)

	// Extracting year
	newRow[2] = ""
	matches := yearRegexp.FindAllString(title, -1)
	if len(matches) > 0 {
		withParenthesis := matches[len(matches)-1]
		newRow[2] = withParenthesis[1 : len(withParenthesis)-1]
	}

	genres := strings.Split(members[2], "|")
	numGenres := 0
	for _, g := range genres {
		if g != "(no genres listed)" {
			newRow[nFixedFields+genreMap[g]] = "1"
			numGenres++
		}
	}
	newRow[nFixedFields-1] = fmt.Sprintf("%d", numGenres)
	// Filling out genres
	for i, s := range newRow[nFixedFields:] {
		if s == "" {
			newRow[i+nFixedFields] = "0"
		}
	}
	return strings.Join(newRow, ",")
}
