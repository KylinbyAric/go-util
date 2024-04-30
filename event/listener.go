package event

import "fmt"

type Listener interface {
	HandleEvent(event Event) error
}

type Listener1 struct {
}

func (a *Listener1) HandleEvent(event Event) error {
	fmt.Println("this listener1 event:%v", event)
	return nil
}

type Listener2 struct {
}

func (l *Listener2) HandleEvent(event Event) error {
	fmt.Println("this listener2 event:%v", event)
	return nil
}
