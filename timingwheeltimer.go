package timer

import (
	"errors"

	"sync"
	"time"

	"github.com/RussellLuo/timingwheel"
	"github.com/google/uuid"
)

type ItemID uuid.UUID // uuid

// func (item ItemID) String() string {
// 	return string(item)
// }

type Item interface {
	// TickDuration time.Duration
	Tick() bool
	Timeout() bool
	GetTickDuration() time.Duration
	GetTimeoutDuration() time.Duration
	GetStartTime() time.Time
	SetStartTime(time.Time)
	GetID() ItemID
}

type timeWheelTimer struct {
	tw      *timingwheel.TimingWheel
	itemMap map[ItemID]*timingwheel.Timer
	mutex   sync.Mutex
}

func NewTimeWheelTimer(wheelTickDuration time.Duration, wheelSize int64) *timeWheelTimer {

	t := &timeWheelTimer{
		tw:      timingwheel.NewTimingWheel(wheelTickDuration, wheelSize),
		itemMap: make(map[ItemID]*timingwheel.Timer),
	}

	return t
}

func (t *timeWheelTimer) wrapper(item Item) func() {

	f := func() {

		if time.Now().After(item.GetStartTime().Add(item.GetTimeoutDuration())) {

			item.Timeout()
			goto exit
		}

		if item.Tick() {

			t.tw.AfterFunc(item.GetTickDuration(), t.wrapper(item))
			return
		}
	exit:
		t.mutex.Lock()
		delete(t.itemMap, item.GetID())
		t.mutex.Unlock()

		return
	}

	return f
}

func (t *timeWheelTimer) Add(item Item) error {
	t.mutex.Lock()

	_, ok := t.itemMap[item.GetID()]
	t.mutex.Unlock()

	if ok {
		return errors.New("ID already in the map")
	}

	item.SetStartTime(time.Now())

	timer := t.tw.AfterFunc(item.GetTickDuration(), t.wrapper(item))
	t.mutex.Lock()
	t.itemMap[item.GetID()] = timer
	t.mutex.Unlock()

	return nil
}

func (t *timeWheelTimer) Start() {
	t.tw.Start()
}

func (t *timeWheelTimer) Stop() {
	t.tw.Stop()
}

func (t *timeWheelTimer) Remove(id ItemID) error {
	t.mutex.Lock()
	_, ok := t.itemMap[id]
	if !ok {
		t.mutex.Unlock()
		return errors.New("ID not in the TimeWheelTimer")
	}
	ok = t.itemMap[id].Stop()
	delete(t.itemMap, id)
	t.mutex.Unlock()
	if !ok {
		return errors.New("ID in the map, but not in the timingwheel")
	}
	return nil

}
