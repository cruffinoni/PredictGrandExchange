package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type data struct {
	dates  []string
	values []float64
}

func fetchFile(path string, s *data) {
	fd, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully opened the CSV file")
	defer fd.Close()
	fileReader := csv.NewReader(fd)
	records, err := fileReader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	for _, r := range records[1:] {
		val, err := strconv.ParseFloat(r[1], 10)
		if err != nil {
			panic(err)
		}
		s.dates = append(s.dates, r[0])
		s.values = append(s.values, val)
	}
}
func main() {

}
