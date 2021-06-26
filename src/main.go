package main

import "time"

func main() {
	parseArgs()
	interval := time.Duration(CLI.Interval) * time.Second

	watch(interval)
}
