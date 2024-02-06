package listener

type ChannelName string
type TableName string
type TriggerName string
type Operation string

const (
	Insert   Operation = "INSERT"
	Update   Operation = "UPDATE"
	Delete   Operation = "DELETE"
	Truncate Operation = "TRUNCATE"
)

type Connector interface {
	// Register регистрирует декодер событий
	// OK
	// InvalidArgument
	// AlreadyExists
	Register(TriggerName, Operation, EventDecoder) error

	// Subscribe подписывается на события в канале
	Subscribe(ChannelName) error

	// Unsubscribe отписывается от событий в канале
	Unsubscribe(ChannelName) error

	// Next возвращает следующее событие после обработки зарегистрированным декодером
	// OK
	// Aborted
	// InvalidEvent
	// EventDecoderNotFound
	Next() (interface{}, error)

	// Close закрывает соединение с базой данных
	Close() error
}

type Event struct {
	TableName   TableName   `json:"table_name"`
	TriggerName TriggerName `json:"trigger_name"`
	Operation   Operation   `json:"operation"`

	Old []byte `json:"old,omitempty"`
	New []byte `json:"new,omitempty"`
}

type EventDecoder func(e *Event) (interface{}, error)
