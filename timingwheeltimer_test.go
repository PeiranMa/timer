package timer_test

import (
	"testing"
	"time"

	"github.com/PeiranMa/timer"
	"github.com/google/uuid"
)

func Testtimer_RemoveFunc(t *testing.T) {
	twt := timer.NewTimeWheelTimer(time.Millisecond, 20)
	twt.Start()
	defer twt.Stop()
	t.Run("", func(t *testing.T) {
		tickC := make(chan time.Time, 2)
		timeoutC := make(chan time.Time, 2)
		id := uuid.New()
		msg := &MessageT{
			ID:              timer.ItemID(id),
			TickDuration:    1 * time.Second,
			TimeoutDuration: 2 * time.Second,
			tickC:           &tickC,
			timeoutC:        &timeoutC,
		}
		twt.Add(msg)
		twt.Remove(timer.ItemID(id))
		time.Sleep(3 * time.Second)
		select {
		case <-tickC:
			t.Errorf("Remove fail")
		case <-timeoutC:
			t.Errorf("Remove fail")
		default:
		}

	})

}

func Testtimer_AddFunc(t *testing.T) {
	twt := timer.NewTimeWheelTimer(time.Millisecond, 20)

	twt.Start()
	defer twt.Stop()
	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}

	timeoutDuration := 2 * time.Second
	for _, d := range durations {
		t.Run("", func(t *testing.T) {
			tickC := make(chan time.Time)
			timeoutC := make(chan time.Time)
			msg := &MessageT{
				TickDuration:    d,
				TimeoutDuration: timeoutDuration,
				tickC:           &tickC,
				timeoutC:        &timeoutC,
			}

			start := time.Now().UTC()

			err := 10 * time.Millisecond
			go func() {
				old := start
				var new time.Time
				var got time.Time
				var min time.Time
				for {
					select {
					case new = <-tickC:
						got = new.Truncate(time.Millisecond)
						min = old.Add(d).Truncate(time.Millisecond)

						if got.Before(min) || got.After(min.Add(err)) {
							t.Errorf("Tick(%s) expiration: want [%s, %s], got %s", d, min, min.Add(err), got)
						}
						old = got
					case <-timeoutC:
						return
					}

				}
			}()
			twt.Add(msg)
			timeoutTime := <-timeoutC
			got := timeoutTime.Truncate(time.Millisecond)

			min := start.Add(timeoutDuration).Truncate(time.Millisecond)

			if got.Before(min) || got.After(min.Add(d/2).Add(err)) {
				t.Errorf("Timeout(%s) expiration: want [%s, %s], got %s", d, min, min.Add(d/2).Add(err), got)
			}

		})

	}

}

func TestSimpleTimer_RemoveFunc(t *testing.T) {
	st := timer.NewSimpleTimer()

	t.Run("", func(t *testing.T) {
		tickC := make(chan time.Time, 2)
		timeoutC := make(chan time.Time, 2)
		id := uuid.New()
		msg := &MessageT{
			ID:              timer.ItemID(id),
			TickDuration:    1 * time.Second,
			TimeoutDuration: 2 * time.Second,
			tickC:           &tickC,
			timeoutC:        &timeoutC,
		}
		st.Add(msg)
		st.Remove(timer.ItemID(id))
		time.Sleep(3 * time.Second)
		select {
		case <-tickC:
			t.Errorf("Remove fail")
		case <-timeoutC:
			t.Errorf("Remove fail")
		default:
		}

	})

}

func TestSimpleTimer_AddFunc(t *testing.T) {
	st := timer.NewSimpleTimer()

	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}

	timeoutDuration := 2 * time.Second
	for _, d := range durations {
		t.Run("", func(t *testing.T) {
			tickC := make(chan time.Time)
			timeoutC := make(chan time.Time)
			msg := &MessageT{
				TickDuration:    d,
				TimeoutDuration: timeoutDuration,
				tickC:           &tickC,
				timeoutC:        &timeoutC,
			}

			start := time.Now().UTC()

			err := 10 * time.Millisecond
			go func() {
				old := start
				var new time.Time
				var got time.Time
				var min time.Time
				for {
					select {
					case new = <-tickC:
						got = new.Truncate(time.Millisecond)
						min = old.Add(d).Truncate(time.Millisecond)

						if got.Before(min) || got.After(min.Add(err)) {
							t.Errorf("Tick(%s) expiration: want [%s, %s], got %s", d, min, min.Add(err), got)
						}
						old = got
					case <-timeoutC:
						return
					}

				}
			}()
			st.Add(msg)
			timeoutTime := <-timeoutC
			got := timeoutTime.Truncate(time.Millisecond)

			min := start.Add(timeoutDuration).Truncate(time.Millisecond)

			if got.Before(min) || got.After(min.Add(d/2).Add(err)) {
				t.Errorf("Timeout(%s) expiration: want [%s, %s], got %s", d, min, min.Add(d).Add(err), got)
			}

		})

	}

}

type MessageT struct {
	ID              timer.ItemID
	TickDuration    time.Duration
	TimeoutDuration time.Duration
	StartTime       time.Time
	tickC           *chan time.Time
	timeoutC        *chan time.Time
}

func (m *MessageT) Tick() bool {

	*m.tickC <- time.Now().UTC()

	return true
}

func (m *MessageT) Timeout() bool {
	//do something
	tmp := time.Now().UTC()
	*m.timeoutC <- tmp
	*m.timeoutC <- tmp
	close(*m.tickC)
	close(*m.timeoutC)

	return true
}
func (m *MessageT) GetTickDuration() time.Duration {
	return m.TickDuration
}
func (m *MessageT) GetTimeoutDuration() time.Duration {
	return m.TimeoutDuration
}
func (m *MessageT) GetStartTime() time.Time {
	return m.StartTime
}
func (m *MessageT) SetStartTime(t time.Time) {
	m.StartTime = t
	return
}
func (m *MessageT) GetID() timer.ItemID {

	return m.ID
}

func (m *MessageT) SetID(id timer.ItemID) {
	m.ID = id
	return
}
