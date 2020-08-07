package timer_test

import (
	"testing"
	"time"

	"github.com/PeiranMa/timer"
	"github.com/RussellLuo/timingwheel"
	"github.com/google/uuid"
)

func genD(i int) time.Duration {
	return time.Duration(i%10000) * time.Millisecond
}

func BenchmarkTimingWheel_StartStop(b *testing.B) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1m", 1000000},
		{"N-5m", 5000000},
		{"N-10m", 10000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]*timingwheel.Timer, c.N)
			for i := 0; i < len(base); i++ {
				base[i] = tw.AfterFunc(genD(i), func() {})
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				tw.AfterFunc(time.Second, func() {}).Stop()
			}

			b.StopTimer()
			for i := 0; i < len(base); i++ {
				base[i].Stop()
			}
		})
	}
}

func BenchmarkStandardTimer_StartStop(b *testing.B) {
	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1m", 1000000},
		{"N-5m", 5000000},
		{"N-10m", 10000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]*time.Timer, c.N)
			for i := 0; i < len(base); i++ {
				base[i] = time.AfterFunc(genD(i), func() {})
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				time.AfterFunc(time.Second, func() {}).Stop()
			}

			b.StopTimer()
			for i := 0; i < len(base); i++ {
				base[i].Stop()
			}
		})
	}
}

func Benchmarktimer_StartStop(b *testing.B) {
	twt := timer.NewTimeWheelTimer(time.Millisecond, 20)
	twt.Start()
	defer twt.Stop()

	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1m", 1000000},
		{"N-5m", 5000000},
		{"N-10m", 10000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]uuid.UUID, c.N)

			for i := 0; i < len(base); i++ {
				id := uuid.New()
				msg := &MessageB{
					ID:              timer.ItemID(id),
					TickDuration:    genD(i),
					TimeoutDuration: 10 * time.Second,
				}
				twt.Add(msg)
				base[i] = id
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				id := uuid.New()

				twt.Add(&MessageB{
					ID:              timer.ItemID(id),
					TickDuration:    time.Second,
					TimeoutDuration: 10 * time.Second,
				})
				twt.Remove(timer.ItemID(id))
			}
			b.StopTimer()
			for i := 0; i < len(base); i++ {
				twt.Remove(timer.ItemID(base[i]))
			}
		})
	}
}

func BenchmarkSimpleTimer_StartStop(b *testing.B) {
	st := timer.NewSimpleTimer()

	cases := []struct {
		name string
		N    int // the data size (i.e. number of existing timers)
	}{
		{"N-1m", 1000000},
		{"N-5m", 5000000},
		{"N-10m", 10000000},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			base := make([]uuid.UUID, c.N)

			for i := 0; i < len(base); i++ {
				id := uuid.New()
				msg := &MessageB{
					ID:              timer.ItemID(id),
					TickDuration:    genD(i),
					TimeoutDuration: 10 * time.Second,
				}
				st.Add(msg)
				base[i] = id
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				id := uuid.New()

				st.Add(&MessageB{
					ID:              timer.ItemID(id),
					TickDuration:    time.Second,
					TimeoutDuration: 10 * time.Second,
				})
				st.Remove(timer.ItemID(id))
			}
			b.StopTimer()
			for i := 0; i < len(base); i++ {
				st.Remove(timer.ItemID(base[i]))
			}
		})
	}
}

type MessageB struct {
	ID              timer.ItemID
	TickDuration    time.Duration
	TimeoutDuration time.Duration
	StartTime       time.Time
}

func (m *MessageB) Tick() bool {

	// *m.tickC <- time.Now().UTC()

	return false
}

func (m *MessageB) Timeout() bool {

	return true
}
func (m *MessageB) GetTickDuration() time.Duration {
	return m.TickDuration
}
func (m *MessageB) GetTimeoutDuration() time.Duration {
	return m.TimeoutDuration
}
func (m *MessageB) GetStartTime() time.Time {
	return m.StartTime
}
func (m *MessageB) SetStartTime(t time.Time) {
	m.StartTime = t
	return
}
func (m *MessageB) GetID() timer.ItemID {

	return m.ID
}

func (m *MessageB) SetID(id timer.ItemID) {
	m.ID = id
	return
}
