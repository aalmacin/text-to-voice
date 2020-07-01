package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Word struct {
	english string
	french  string
}

func main() {
	csvFile, err := os.Open("input.csv")

	if err != nil {
		panic(err)
	}

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		panic(err)
	}

	for _, line := range csvLines {
		word := Word{line[0], line[1]}
		fmt.Println(word)
	}

	defer csvFile.Close()
}
