package main

import (
	"flag"
	"log"
	"os"

	"github.com/ekotlikoff/rentorbuy/internal/data"
)

func main() {
	var f = flag.String("f", "", "input data file")
	flag.Parse(file)
	d, err := os.ReadFile(*f)
	if err != nil {
		log.Fatalf(err)
	}
	s := data.LoadScenario(d)
	s.visualize()
}
