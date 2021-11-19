package entities

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

type (
	Message struct {
		FromIP   string
		FromPort int
		Method   string
		Body     string
	}

	Messages struct {
		data  []Message
		mutex sync.Mutex
	}
)

func (m *Message) String() string {
	return fmt.Sprintf("%s|%d|%s|%s|", m.FromIP, m.FromPort, m.Method, m.Body)
}

func (m *Messages) empty() bool {
	return len(m.data) == 0
}

func (m *Messages) Empty() bool {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	return m.empty()
}

func (m *Messages) Push(message Message) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	m.data = append(m.data, message)
}

func (m *Messages) Pop() (Message, error) {
	defer m.mutex.Unlock()
	m.mutex.Lock()
	if m.empty() {
		return Message{}, errors.New("Queue is empty")
	}
	el := m.data[0]
	m.data = m.data[1:]
	return el, nil
}

func NewMessage(str string) (Message, error) {
	message := Message{}

	reader := strings.NewReader(str)
	buffer := new(bytes.Buffer)

	// FromIP
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			return Message{}, err
		}

		if ch == '|' {
			message.FromIP = buffer.String()
			buffer.Reset()
			break
		}

		buffer.WriteRune(ch)
	}

	// FromPort
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			return Message{}, err
		}

		if ch == '|' {
			message.FromPort, err = strconv.Atoi(buffer.String())
			if err != nil {
				return Message{}, err
			}
			buffer.Reset()
			break
		}

		buffer.WriteRune(ch)
	}

	// Method
	for {
		ch, _, err := reader.ReadRune()
		if err != nil {
			return Message{}, err
		}

		if ch == '|' {
			message.Method = buffer.String()
			buffer.Reset()
			break
		}

		buffer.WriteRune(ch)
	}

	// Body
	for {
		ch, _, err := reader.ReadRune()
		if err == io.EOF {
			message.Body = buffer.String()
			break
		}

		if err != nil {
			return Message{}, err
		}

		buffer.WriteRune(ch)
	}

	return message, nil
}
