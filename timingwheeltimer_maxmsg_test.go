package timer_test

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/PeiranMa/timer"
	"github.com/google/uuid"
)

func TestSimpleTimer(t *testing.T) {
	st := timer.NewSimpleTimer()
	var mutex sync.Mutex
	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
	timeoutDuration := 10 * time.Second

	msgCount := 100000
	t.Run("", func(t *testing.T) {
		C := make(chan timer.ItemID)
		exitChan := make(chan int)
		set := make(map[timer.ItemID]struct{})
		tmp := &set
		go func() {
			for {
				select {
				case rec := <-C:
					mutex.Lock()
					_, ok := (*tmp)[rec]
					if ok {
						mutex.Unlock()
						continue
					}
					(*tmp)[rec] = struct{}{}
					mutex.Unlock()
					if len((*tmp)) == msgCount {
						exitChan <- 1

					}

				}
			}

		}()
		for i := 0; i < msgCount; i++ {
			st.Add(&Message{
				ID:              timer.ItemID(uuid.New()),
				TickDuration:    durations[rand.Intn(len(durations))],
				TimeoutDuration: timeoutDuration,
				C:               &C,
			})
		}
		<-exitChan
		close(C)
	})

}

func TestTimeWheelTimer(t *testing.T) {
	twt := timer.NewTimeWheelTimer(time.Millisecond, 20)
	twt.Start()
	defer twt.Stop()

	var mutex sync.Mutex
	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
	timeoutDuration := 10 * time.Second

	msgCount := 1000000
	t.Run("", func(t *testing.T) {
		C := make(chan timer.ItemID)
		exitChan := make(chan int)
		set := make(map[timer.ItemID]struct{})
		tmp := &set
		go func() {
			for {
				select {
				case rec := <-C:
					mutex.Lock()
					_, ok := (*tmp)[rec]
					if ok {
						mutex.Unlock()
						continue
					}
					(*tmp)[rec] = struct{}{}
					mutex.Unlock()
					if len((*tmp)) == msgCount {
						exitChan <- 1

					}

				}
			}

		}()
		for i := 0; i < msgCount; i++ {
			twt.Add(&Message{
				ID:              timer.ItemID(uuid.New()),
				TickDuration:    durations[rand.Intn(len(durations))],
				TimeoutDuration: timeoutDuration,
				C:               &C,
			})
		}
		<-exitChan
		close(C)
	})
}

type Message struct {
	ID              timer.ItemID
	TickDuration    time.Duration
	TimeoutDuration time.Duration
	StartTime       time.Time
	C               *chan timer.ItemID
}

func (m *Message) Tick() bool {

	// *m.C <- m.ID

	return true
}

func (m *Message) Timeout() bool {

	*m.C <- m.ID

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
