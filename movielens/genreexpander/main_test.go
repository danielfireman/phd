package main

import (
	"strings"
	"testing"
)

func TestProcessRow(t *testing.T) {
	dataTable := []struct {
		input string
		want  string
	}{
		{
			`1,Toy Story (1995),Adventure|Animation|Children|Comedy|Fantasy`,
			`1,"Toy Story (1995)",1995,5,1,1,1,0,1,0,0,0,1,0,0,0,0,0,0,0,0,0`,
		},
		{
			`2,Jumanji (1995),Adventure|Children|Fantasy`,
			`2,"Jumanji (1995)",1995,3,1,1,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0`,
		},
		{
			`3,Grumpier Old Men (1995),Comedy|Romance`,
			`3,"Grumpier Old Men (1995)",1995,2,0,0,0,0,1,0,0,0,0,0,0,0,0,1,0,0,0,0`,
		},
		{
			`4,Waiting to Exhale (1995),Comedy|Drama|Romance`,
			`4,"Waiting to Exhale (1995)",1995,3,0,0,0,0,1,0,0,1,0,0,0,0,0,1,0,0,0,0`,
		},
		{
			`5,Father of the Bride Part II (1995),Comedy`,
			`5,"Father of the Bride Part II (1995)",1995,1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0`,
		},
		{
			`6,Heat (1995),Action|Crime|Thriller`,
			`6,"Heat (1995)",1995,3,1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,1,0,0`,
		},
		{
			`7,Sabrina (1995),Comedy|Romance`,
			`7,"Sabrina (1995)",1995,2,0,0,0,0,1,0,0,0,0,0,0,0,0,1,0,0,0,0`,
		},
		{
			`8,Tom and Huck (1995),Adventure|Children`,
			`8,"Tom and Huck (1995)",1995,2,1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0`,
		},
		{
			`9,Sudden Death (1995),Action`,
			`9,"Sudden Death (1995)",1995,1,1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0`,
		},
		{
			`317,"Santa Clause, The (1994)",Comedy|Drama|Fantasy`,
			`317,"Santa Clause, The (1994)",1994,3,0,0,0,0,1,0,0,1,1,0,0,0,0,0,0,0,0,0`,
		},
		{
			`203,"To Wong Foo, Thanks for Everything! Julie Newmar (1995)",Comedy`,
			`203,"To Wong Foo, Thanks for Everything! Julie Newmar (1995)",1995,1,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0,0,0`,
		},
		{
			`126929,Li'l Quinquin,(no genres listed)`,
			`126929,"Li'l Quinquin",,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0`,
		},
		{
			`7789,"11'09""01 - September 11 (2002)",Drama`,
			`7789,"11'09""01 - September 11 (2002)",2002,1,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0`,
		},
	}
	for _, d := range dataTable {
		members, err := csvReader(strings.NewReader(d.input)).ReadAll()
		if err != nil {
			t.Errorf("Error processing input: %s. Err:%q", d.input, err)
		}
		if len(members) != 1 {
			t.Errorf("Invalid number of CSV records. got:%d want:1", len(members))
		}
		got := processRow(members[0])
		if d.want != got {
			t.Errorf("got:%s want:%s", got, d.want)
		}
	}
}
