package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func valueToLineData(d *data) []opts.LineData {
	log.Print("Translating records values to a slice of data (Y axis)")
	items := make([]opts.LineData, 0)
	for _, val := range d.values {
		items = append(items, opts.LineData{Value: val})
	}
	return items
}

func calculateMAE(d1, d2 *data) float64 {
	n := len(d1.values)
	totalError := 0.0

	for i := 0; i < n; i++ {
		totalError += math.Abs(d1.values[i] - d2.values[i])
	}

	return totalError / float64(n)
}

type data struct {
	dates  []string
	values []float64
	min    int64
	max    int64
}

type graph struct {
	prices *data
	volume *data
}

func (g *graph) writeGraphToFile(w http.ResponseWriter, _ *http.Request) {
	// create a new line instance
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    "Prix et volume des runes d'air",
			Subtitle: "Graphique afin d'observer, ou non, la corrélation entre le prix et le volume pour les runes d'air à partir 03/03/2023 -> 29/08/2023",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Trigger: "axis",
		}),
		charts.WithLegendOpts(opts.Legend{
			Show: true,
		}),
	)

	// Put data into instance
	line.SetXAxis(g.prices.dates).
		AddSeries("Prix",
			valueToLineData(g.prices),
		).
		AddSeries("Volume",
			valueToLineData(g.volume),
		).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{
			Symbol: "roundRect",
		}))
	line.Render(w)
}

func (g *graph) fetchFile(path string, s *data, newMin int64, newMax int64) {
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
	s.min = math.MaxInt
	s.max = math.MinInt
	for _, r := range records[1:] {
		val, err := strconv.ParseInt(r[1], 10, 64)
		if err != nil {
			panic(err)
		}
		s.dates = append(s.dates, r[0])
		if val > s.max {
			s.max = val
		}
		if val < s.min {
			s.min = val
		}
	}

	for _, r := range records[1:] {
		val, err := strconv.ParseInt(r[1], 10, 64)
		if err != nil {
			panic(err)
		}
		s.values = append(s.values, float64(newMin)+(float64(val-s.min))*((float64(newMax-newMin))/(float64(s.max-s.min))))
	}
}

func main() {
	g := &graph{
		prices: &data{},
		volume: &data{},
	}
	g.fetchFile("data/price.csv", g.prices, -1, 1)
	g.fetchFile("data/volume.csv", g.volume, -1, 1)
	log.Printf("MAE: %.4f", calculateMAE(g.prices, g.volume))
	http.HandleFunc("/", g.writeGraphToFile)
	http.ListenAndServe(":8081", nil)
}
