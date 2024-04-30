package event

import "context"

type Event struct {
	Ctx  context.Context
	Data interface{} `json:"data"`
}
