package listener

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/lib/pq"
)

type connector struct {
	mutex    sync.RWMutex
	listener *pq.Listener
	channels map[ChannelName]struct{}
	decoders map[string]EventDecoder
}

func New(config *Config) (Connector, error) {
	err := config.Check()
	if err != nil {
		return nil, &Error{
			Code: InvalidArgument,
			Err:  err,
		}
	}
	listener := pq.NewListener(
		config.String(),
		config.MinReconnectInterval,
		config.MaxReconnectInterval,
		nil,
	)
	if listener == nil {
		return nil, &Error{
			Code: Internal,
			Err:  errors.New("empty postgres listener"),
		}
	}
	return &connector{
		listener: listener,
		channels: make(map[ChannelName]struct{}),
		decoders: make(map[string]EventDecoder),
	}, nil
}

func (c *connector) Subscribe(name ChannelName) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.channels[name]
	if ok {
		return &Error{
			Code: AlreadyExists,
			Err:  fmt.Errorf("A subscription to channel (%s) already exists", name),
		}
	}
	err := c.listener.Listen(string(name))
	if err != nil {
		return &Error{
			Code: Internal,
			Err:  err,
		}
	}
	c.channels[name] = struct{}{}
	return nil
}

func (c *connector) Unsubscribe(name ChannelName) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, ok := c.channels[name]
	if !ok {
		return nil
	}
	err := c.listener.Unlisten(string(name))
	if err != nil {
		return &Error{
			Code: Internal,
			Err:  err,
		}
	}
	delete(c.channels, name)
	return nil
}

func (c *connector) Register(
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
	_, ok := c.decoders[key]
	if ok {
		return &Error{
			Code: AlreadyExists,
			Err:  errors.New("decoder is already registered"),
		}
	}
	c.decoders[key] = dec
	return nil
}

func (c *connector) Next() (interface{}, error) {
	data, ok := <-c.listener.Notify
	if !ok {
		return nil, &Error{
			Code: Aborted,
			Err:  errors.New("closed channel"),
		}
	}
	event := new(Event)
	err := json.Unmarshal([]byte(data.Extra), event)
	if err != nil {
		return nil, &Error{
			Code: InvalidEvent,
			Err:  err,
		}
	}
	key := string(event.TriggerName) + ":" + string(event.Operation)
	f, ok := c.decoders[key]
	if !ok {
		return nil, &Error{
			Code: EventDecoderNotFound,
			Err:  errors.New("event_decoder not found"),
		}
	}
	return f(event)
}

func (c *connector) Close() error {
	return c.listener.Close()
}
