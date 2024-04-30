package event

import "fmt"

var FinalAppEventManager = NewAppEventManager()

type AppEventManager struct {
	ListenerContainer []Listener
	EventChan         chan Event
	Closed            bool
}

func (a *AppEventManager) RegisterListener(listener Listener) {
	if a.ListenerContainer == nil {
		a.ListenerContainer = make([]Listener, 0)
	}
	a.ListenerContainer = append(a.ListenerContainer, listener)
}

func (a *AppEventManager) PublishEvent(event Event) {
	if a.Closed {
		return
	}
	a.EventChan <- event
}

func NewAppEventManager() *AppEventManager {
	return &AppEventManager{
		EventChan: make(chan Event, 0),
	}
}

func (a *AppEventManager) HandleEvents() {
	go func() {
		for {
			select {
			case event, ok := <-a.EventChan:
				if ok {
					for _, listener := range a.ListenerContainer {
						err := listener.HandleEvent(event)
						if err != nil {
							return
						}
					}
				} else {
					// 为啥会一直打印？
					fmt.Println("manager is closed")
					break
				}
			}
		}
	}()
}

func (a *AppEventManager) Close() {
	a.Closed = true
	close(a.EventChan)
}
