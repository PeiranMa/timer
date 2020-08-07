package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/PeiranMa/timer"
	"github.com/google/uuid"
)

func Twt(num int, tickDuration time.Duration, distribution string, timeRange int) {
	twt := timer.NewTimeWheelTimer(time.Millisecond, 20)

	twt.Start()
	defer twt.Stop()

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println(runtime.NumGoroutine())
			for i := 0; i < num; i++ {
				msg := &Message{
					ID:              timer.ItemID(uuid.New()),
					TickDuration:    tickDuration * time.Second,
					TimeoutDuration: time.Duration(timeoutGenerator(distribution, timeRange)) * time.Second,
				}

				twt.Add(msg)

			}
		}

	}

}
