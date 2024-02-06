package listener

import (
	"errors"
	"fmt"
	"sync"

	"github.com/lib/pq"
)

type connector struct {
	mutex    sync.RWMutex
	listener *pq.Listener
	channels map[ChannelName]struct{}
	router   Router
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
		router:   NewRouter(),
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
	return c.router.Register(triggerName, operation, dec)
}

func (c *connector) Next() (interface{}, error) {
	data, ok := <-c.listener.Notify
	if !ok {
		return nil, &Error{
			Code: Aborted,
			Err:  errors.New("closed channel"),
		}
	}
	return c.router.Decode([]byte(data.Extra))
}

func (c *connector) Notify() <-chan *pq.Notification {
	return c.listener.Notify
}

func (c *connector) Close() error {
	return c.listener.Close()
}
