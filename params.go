package listener

import (
	"fmt"
	"strings"
)

// params это список параметров для подключения к базе данных
type params struct {
	data []string
}

// Add добавляет параметр в список параметров
func (p *params) Add(key, value string) {
	p.data = append(
		p.data,
		fmt.Sprintf(
			"%s=%s",
			key,
			value,
		),
	)
}

// String возвращает строковое представление
func (p *params) String() string {
	return strings.Join(p.data, " ")
}
