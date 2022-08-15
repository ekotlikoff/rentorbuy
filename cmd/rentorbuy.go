package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	data "github.com/ekotlikoff/rentorbuy/internal"
)

func main() {
	var f = flag.String("f", "", "input data file")
	flag.Parse()
	d, err := os.ReadFile(*f)
	if err != nil {
		log.Fatalf(err.Error())
	}
	s := data.LoadScenario(d)
	s.Visualize()
	http.Handle("/", http.FileServer(http.Dir("./data")))
	http.ListenAndServe(":3000", nil)
}
