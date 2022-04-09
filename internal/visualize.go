package data

import (
	"os"

	"github.com/go-echarts/go-echarts/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func (s *Scenario) visualize() {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title: "When will you break even after buying?",
		}))
	line.AddSeries(s.House.Address, generateLineItems())
	f, _ := os.Create("line.html")
	line.Render(f)
}
func (s *Scenario) generateLineItems() []opts.LineData {
	items := make([]opts.LineData, 100)
	for m := 0; m < 12*10; m++ {
		items = append(items, opts.LineData{Value: s.moneyLostBuying - s.moneyLostRenting})
	}
	return items
}
