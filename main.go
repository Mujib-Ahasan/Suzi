package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

type plotC struct {
	results []Result
	attack  string
}

func main() {
	hfs := flag.NewFlagSet("suzi", flag.ExitOnError)

	url := hfs.String("url", " ", "Site where you want to attack")
	numOfReq := hfs.Int("req", 10, "Number of requests to send")
	timeout := hfs.Int("timeout", 5, "Request timeout in seconds")
	attacktype := hfs.String("atk", " ", "type of attack")
	plot := hfs.Bool("plot", false, "Do ya wanna plot da test as a timeseries ?")
	numCPUS := hfs.Int("cpus", runtime.NumCPU(), "Number of CPUs to use")
	method := hfs.String("method", "GET", "HTTP method to use (GET, POST, etc.)")
	rate := hfs.Int("rate", 1, "Number of requests per second")

	hfs.Parse(os.Args[1:])

	runtime.GOMAXPROCS(*numCPUS)

	// results := make([]Result, *numOfReq)

	var pc plotC

	switch strings.ToLower(*attacktype) {
	case "basic":
		pc = plotC{results: basicAttack(*url, *numOfReq, *rate, *method, *timeout), attack: "Basic"}
	case "burst":
		pc = plotC{results: burstAttack(*url, *numOfReq, *method, *timeout), attack: "Burst"}
	case "random":
		pc = plotC{results: randomLoadAttack(*url, *numOfReq, *method, *rate, *timeout), attack: "random"}
	default:
		fmt.Println("Unknown attack type:", *attacktype)
		return
	}

	if *plot {
		plotResults(pc)
	}
}
