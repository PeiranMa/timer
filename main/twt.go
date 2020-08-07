package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/PeiranMa/timer"
	"github.com/google/uuid"
)

func Twt(num int, tickDuration time.Duration) {
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
					TimeoutDuration: time.Duration(timeoutGenerator("uniform")) * time.Second,
				}

				twt.Add(msg)

			}
		}

	}

}
