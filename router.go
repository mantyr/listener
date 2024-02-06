package listener

import (
	"encoding/json"
	"errors"
	"fmt"
)

type router struct {
	decoders map[string]EventDecoder
}

func NewRouter() Router {
	return &router{
		decoders: make(map[string]EventDecoder),
	}
}

func (r *router) Register(
	triggerName TriggerName,
	operation Operation,
	dec EventDecoder,
) error {
	switch {
	case triggerName == "":
		return &Error{
			Code: InvalidArgument,
			Err:  errors.New("empty trigger_name"),
		}
	case operation == "":
		return &Error{
			Code: InvalidArgument,
			Err:  errors.New("empty operation"),
		}
	case dec == nil:
		return &Error{
			Code: InvalidArgument,
			Err:  errors.New("empty event_decoder"),
		}
	}
	switch operation {
	case Insert, Update, Delete:
	default:
		return &Error{
			Code: InvalidArgument,
			Err:  fmt.Errorf("unexpected operation (%s)", operation),
		}
	}
	key := string(triggerName) + ":" + string(operation)
	_, ok := r.decoders[key]
	if ok {
		return &Error{
			Code: AlreadyExists,
			Err:  errors.New("decoder is already registered"),
		}
	}
	r.decoders[key] = dec
	return nil
}

// Decode декодирует событие для последующей типизации
func (r *router) Decode(data []byte) (interface{}, error) {
	event := new(Event)
	err := json.Unmarshal([]byte(data), event)
	if err != nil {
		return nil, &Error{
			Code: InvalidEvent,
			Err:  err,
		}
	}
	key := string(event.TriggerName) + ":" + string(event.Operation)
	f, ok := r.decoders[key]
	if !ok {
		return nil, &Error{
			Code: EventDecoderNotFound,
			Err:  errors.New("event_decoder not found"),
		}
	}
	return f(event)
}
