package main

import (
	"time"
)

func main() {
	num := 4000
	tickD := time.Duration(1)
	distribution := "exp"
	timeRange := 60

	// Twt(num, tickD, distribution, timeRange)
	Simple(num, tickD, distribution, timeRange)

}
