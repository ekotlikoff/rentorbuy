package data

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func (s *Scenario) GenerateChart() *os.File {
	line := charts.NewLine()
	tooltipFormatter := opts.FuncOpts(`function (params) {
		return Math.floor(params.value[1]).toLocaleString('en-US', { style: 'currency', currency: 'USD' })+' gain after '+Math.floor(params.value[0]/12)+' years';
	}`)
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			PageTitle: "Rent Or Buy?",
			Theme:     types.ThemeWesteros,
		}),
		charts.WithTitleOpts(opts.Title{
			Title: "When will you break even after buying?",
		}),
		charts.WithLegendOpts(opts.Legend{Show: true}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Formatter: tooltipFormatter}),
	)
	line.AddSeries(
		s.House.Address, s.generateLineItems(),
	)
	f, err := os.CreateTemp("", "rentorybuy-*.html")
	if err != nil {
		log.Fatalf("Failed to create temp file")
	}

	// Can we add custom html here with the scenario data?
	line.Render(f)
	return f
}

func (s *Scenario) generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 0, monthsToCalc)
	buying := s.monthlyMoneyMadeBuying(float64(monthsToCalc))
	renting := s.monthlyMoneyMadeRenting(float64(monthsToCalc))
	for m := 0; m < monthsToCalc; m++ {
		items = append(items, opts.LineData{
			Value: []float64{float64(m), buying[m] - renting[m]},
		})
	}
	return items
}

func (s *Scenario) writeToCSV() {
	f, err := os.Create("data/data.csv")
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(f)
	w.Write([]string{"month", "value"})
	buying := s.monthlyMoneyMadeBuying(float64(monthsToCalc))
	renting := s.monthlyMoneyMadeRenting(float64(monthsToCalc))
	for m := 0; m < monthsToCalc; m++ {
		w.Write([]string{strconv.Itoa(m), strconv.Itoa(int(buying[m] - renting[m]))})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
