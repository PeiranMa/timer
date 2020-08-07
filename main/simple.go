package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/PeiranMa/timer"
	"github.com/google/uuid"
)

func Simple(num int, tickDuration time.Duration) {
	st := timer.NewSimpleTimer()

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			fmt.Println(runtime.NumGoroutine())
			for i := 0; i < num; i++ {
				msg := &Message{
					ID:              timer.ItemID(uuid.New()),
					TickDuration:    tickDuration * time.Second,
					TimeoutDuration: time.Duration(timeoutGenerator("exp")) * time.Second,
				}

				st.Add(msg)

			}
		}

	}
}
