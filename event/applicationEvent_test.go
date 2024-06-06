package event

import (
	"testing"
	"time"
)

func TestAppEventManager_RegisterListener(t *testing.T) {
	FinalAppEventManager.RegisterListener(&Listener1{})
	FinalAppEventManager.RegisterListener(&Listener2{})
	FinalAppEventManager.HandleEvents()

	FinalAppEventManager.PublishEvent(Event{Data: 1})
	go func() {
		FinalAppEventManager.PublishEvent(Event{Data: 2})
	}()

	go func() {
		FinalAppEventManager.PublishEvent(Event{Data: 3})
	}()
	go func() {
		time.Sleep(time.Second * 5)
		FinalAppEventManager.PublishEvent(Event{Data: 4})
	}()

	time.Sleep(time.Second * 10)
	FinalAppEventManager.Close()
	time.Sleep(time.Second * 10)
}
