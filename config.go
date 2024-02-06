package listener

import (
	"errors"
	"fmt"
	"time"
)

// Config это минимальный набор данных для доступа к базе данных
type Config struct {
	// Host это адрес базы данных
	Host string

	// Port это порт по которому доступна база данных
	Port int

	// Name это название базы данных
	Name string

	// User это имя пользователя в базе данных
	User string

	// Password это пароль пользователя в базе данных
	Password string

	// MinReconnectInterval это минимальное время в секундах
	// до попытки восстановить соединение
	// После каждой попытки этот интервал удваивается до максимального
	MinReconnectInterval time.Duration

	// MaxReconnectInterval это максимальное время в секундах
	// до попытки восстановить соединение
	MaxReconnectInterval time.Duration
}

// String возвращает текстовое представление
func (c *Config) String() string {
	params := &params{}
	params.Add("host", c.Host)
	params.Add("port", fmt.Sprintf("%d", c.Port))
	params.Add("dbname", c.Name)
	params.Add("sslmode", "disable")
	params.Add("user", c.User)
	if c.Password != "" {
		params.Add("password", c.Password)
	}
	return params.String()
}

func (c *Config) Check() error {
	switch {
	case c.Host == "":
		return errors.New("empty config db host")
	case c.Port <= 0:
		return errors.New("empty config db port")
	case c.Name == "":
		return errors.New("empty config db name")
	case c.User == "":
		return errors.New("empty config db user")
	case c.MinReconnectInterval <= 0:
		return errors.New("empty config min_reconnect_interval")
	case c.MaxReconnectInterval <= 0:
		return errors.New("empty config max_reconnect_interval")
	}
	return nil
}
