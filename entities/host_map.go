package entities

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

type (
	HostMap struct {
		data  map[string]int
		mutex sync.Mutex
	}
)

func NewHostMap() HostMap {
	return HostMap{
		data:  make(map[string]int),
		mutex: sync.Mutex{},
	}
}

func (h *HostMap) Set(value string, delimiter string) error {
	defer h.mutex.Unlock()
	h.mutex.Lock()

	splited := strings.Split(value, delimiter)
	for _, value := range splited {
		host_port := strings.Split(value, ":")
		if len(host_port) != 2 {
			return errors.New("Failed to split " + value)
		}

		host := host_port[0]
		port, err := strconv.Atoi(host_port[1])
		if err != nil {
			return err
		}

		h.data[host] = port
	}

	return nil
}

func (h *HostMap) Foreach(f func(key string, value int)) {
	defer h.mutex.Unlock()
	h.mutex.Lock()
	for k, v := range h.data {
		f(k, v)
	}
}

func (h *HostMap) String(delimiter string) string {
	defer h.mutex.Unlock()
	h.mutex.Lock()

	buffer := new(bytes.Buffer)
	for key, value := range h.data {
		fmt.Fprintf(buffer, "%s:%d%s", key, value, delimiter)
	}
	return buffer.String()
}

func (h *HostMap) empty() bool {
	return len(h.data) == 0
}

func (h *HostMap) Empty() bool {
	defer h.mutex.Unlock()
	h.mutex.Lock()
	return h.empty()
}

func (h *HostMap) exist(key string) bool {
	_, find := h.data[key]
	return find
}

func (h *HostMap) Exist(key string) bool {
	defer h.mutex.Unlock()
	h.mutex.Lock()
	_, find := h.data[key]
	return find
}

func (h *HostMap) Get(key string) (int, error) {
	defer h.mutex.Unlock()
	h.mutex.Lock()

	if !h.exist(key) {
		return 0, errors.New("Not found host")
	}

	return h.data[key], nil
}

func (h *HostMap) Add(key string, value int) {
	defer h.mutex.Unlock()
	h.mutex.Lock()

	h.data[key] = value
}

func (h *HostMap) Remove(key string) error {
	defer h.mutex.Unlock()
	h.mutex.Lock()

	if !h.exist(key) {
		return errors.New("Not found host")
	}

	delete(h.data, key)

	return nil
}
