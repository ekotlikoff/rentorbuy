package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	data "github.com/ekotlikoff/rentorbuy/internal"
)

func main() {
	var f = flag.String("f", "", "input data file")
	var i = flag.Bool("i", false, "interactive")
	flag.Parse()
	log.SetOutput(ioutil.Discard)
	if *i {
		p := tea.NewProgram(data.InitialModel())
		if err := p.Start(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		return
	}
	d, err := os.ReadFile(*f)
	if err != nil {
		log.Fatalf(err.Error())
	}
	s := data.LoadScenario(d)
	file := s.GenerateChart()
	defer os.Remove(file.Name())
	cmd := exec.Command("open", file.Name())
	cmd.Run()
	w := bufio.NewWriter(os.Stdout)
	w.WriteString("Opening chart... press enter to continue\n")
	w.Flush()
	r := bufio.NewReader(os.Stdin)
	r.ReadString('\n')
}
