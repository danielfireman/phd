package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

var genreMap = map[string]int{
	"Action":             0,
	"Adventure":          1,
	"Animation":          2,
	"Children's":         3,
	"Comedy":             4,
	"Crime":              5,
	"Documentary":        6,
	"Drama":              7,
	"Fantasy":            8,
	"Film-Noir":          9,
	"Horror":             10,
	"Musical":            11,
	"Mystery":            12,
	"Romance":            13,
	"Sci-Fi":             14,
	"Thriller":           15,
	"War":                16,
	"Western":            17,
	"(no genres listed)": 18,
}

func main() {
	records, err := csv.NewReader(os.Stdin).ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, members := range records {
		// Very simple way to skip the title.
		if members[0] == "movieId" && members[1] == "title" && members[2] == "genres" {
			continue
		}
		genres := strings.Split(members[2], "|")

		// id, title, #genres counter of all genres
		newRow := make([]string, 3+len(genreMap))
		newRow[0] = members[0]
		newRow[1] = fmt.Sprintf("\"%s\"", members[1])
		newRow[2] = fmt.Sprintf("%d", len(genres))
		for _, g := range genres {
			newRow[3+genreMap[g]] = "1"
		}
		// Filling out genres
		for i, s := range newRow {
			if s == "" {
				newRow[i] = "0"
			}
		}
		fmt.Println(strings.Join(newRow, ","))
	}
}
