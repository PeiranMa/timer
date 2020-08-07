package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/PeiranMa/timer"
)

type Message struct {
	ID              timer.ItemID
	TickDuration    time.Duration
	TimeoutDuration time.Duration
	StartTime       time.Time
}

func (m *Message) Tick() bool {

	time.Sleep(10 * time.Millisecond)

	return true
}

func (m *Message) Timeout() bool {

	return true
}
func (m *Message) GetTickDuration() time.Duration {
	return m.TickDuration
}
func (m *Message) GetTimeoutDuration() time.Duration {
	return m.TimeoutDuration
}
func (m *Message) GetStartTime() time.Time {
	return m.StartTime
}
func (m *Message) SetStartTime(t time.Time) {
	m.StartTime = t
	return
}
func (m *Message) GetID() timer.ItemID {

	return m.ID
}

func (m *Message) SetID(id timer.ItemID) {
	m.ID = id
	return
}

func timeoutGenerator(distribution string) int {
	switch {
	case distribution == "exp":
		rand.Seed(time.Now().UnixNano())
		preNum := rand.Intn(600) + 1
		x := math.Log(float64(preNum))
		x = (1 - x/math.Log(601)) * 60
		// math.P
		return int(x)
	case distribution == "uniform":
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(60)
		return num
	default:
		rand.Seed(time.Now().UnixNano())
		num := rand.Intn(60)
		return num

	}

}
