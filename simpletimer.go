package timer

import (
	"errors"

	"sync"
	"time"
)

type simpletimer struct {
	itemMap map[ItemID]*time.Timer
	mutex   sync.Mutex
}

func NewSimpleTimer() *simpletimer {

	t := &simpletimer{

		itemMap: make(map[ItemID]*time.Timer),
	}

	return t
}

func (t *simpletimer) wrapper(item Item) func() {

	f := func() {

		if time.Now().After(item.GetStartTime().Add(item.GetTimeoutDuration())) {

			item.Timeout()
			goto exit
		}

		if item.Tick() {

			time.AfterFunc(item.GetTickDuration(), t.wrapper(item))
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

func (t *simpletimer) Add(item Item) error {
	t.mutex.Lock()

	_, ok := t.itemMap[item.GetID()]
	t.mutex.Unlock()

	if ok {
		return errors.New("ID already in the map")
	}

	item.SetStartTime(time.Now())

	timer := time.AfterFunc(item.GetTickDuration(), t.wrapper(item))
	t.mutex.Lock()
	t.itemMap[item.GetID()] = timer
	t.mutex.Unlock()

	return nil
}

func (t *simpletimer) Remove(id ItemID) error {
	t.mutex.Lock()
	_, ok := t.itemMap[id]
	if !ok {
		t.mutex.Unlock()
		return errors.New("ID not in the TimeWheelTimer")
	}
	ok = t.itemMap[id].Stop() //should pay attention to whether my call stop the timer?
	delete(t.itemMap, id)
	t.mutex.Unlock()
	if !ok {
		return errors.New("ID in the map, but not in the timingwheel")
	}
	return nil

}
