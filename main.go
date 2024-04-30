package main

import (
	"awesomeProject1/event"
	"time"
)

func main() {
	event.FinalAppEventManager.RegisterListener(&event.Listener1{})
	event.FinalAppEventManager.RegisterListener(&event.Listener2{})
	event.FinalAppEventManager.HandleEvents()

	event.FinalAppEventManager.PublishEvent(event.Event{Data: 1})
	go func() {
		event.FinalAppEventManager.PublishEvent(event.Event{Data: 2})
	}()

	go func() {
		event.FinalAppEventManager.PublishEvent(event.Event{Data: 3})
	}()
	go func() {
		time.Sleep(time.Second * 5)
		event.FinalAppEventManager.PublishEvent(event.Event{Data: 4})
	}()

	time.Sleep(time.Second * 10)
	event.FinalAppEventManager.Close()
	time.Sleep(time.Second * 10)

}
