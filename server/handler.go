package server

import (
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"cs7ns1/project3/entities"

	"github.com/withmandala/go-log"
)

func parseIPAndPort(addr string) (string, int, error) {
	result := strings.Split(addr, ":")

	if len(result) != 2 {
		return "", 0, errors.New("Failed to split addr: " + addr)
	}

	ip := result[0]
	port, err := strconv.Atoi(result[1])

	return ip, port, err
}

func handleIncomingMessage(conn *net.TCPConn, message entities.Message) {
	logger := log.New(os.Stderr)

	switch message.Method {
	case "MSG":
		Inbox.Push(message)
	case "FETCH_REQ":
		ip, port, err := parseIPAndPort(conn.LocalAddr().String())
		if err != nil {
			logger.Error(err)
			return
		}
		message := entities.Message{
			FromIP:   ip,
			FromPort: port,
			Method:   "FETCH_RESP",
			Body:     hostMap.String(";"),
		}
		conn.Write([]byte(message.String()))
	case "FETCH_RESP":
		body := message.Body
		hostMap.Set(body, ";")
	case "REGISTER":
		hostMap.Add(message.FromIP, message.FromPort)
	case "UNREGISTER":
		hostMap.Remove(message.FromIP)
	}
}

func handleIncomingStream(str string, builder *strings.Builder, delimiterCount *int, escapeChar *bool) (entities.Message, bool, error) {
	var message = entities.Message{}
	var hasNew = false
	var err error = nil

	for _, ch := range string(str) {
		switch ch {
		case '|':
			if *escapeChar {
				*escapeChar = false
			} else {
				*delimiterCount++
			}
		case '\\':
			*escapeChar = !*escapeChar
		}

		if *delimiterCount == 4 {
			str := strings.TrimSpace(builder.String())
			message, err = entities.NewMessage(str)
			hasNew = true
			*delimiterCount = 0
			builder.Reset()
		} else {
			builder.WriteRune(ch)
		}
	}

	return message, hasNew, err

}

func handleTCPConnection(conn *net.TCPConn) {
	defer conn.Close()

	logger := log.New(os.Stderr)

	delimiterCount := 0
	escapeChar := false

	var builder = new(strings.Builder)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Error(err)
		}

		message, hasNew, err := handleIncomingStream(string(buffer[:n]), builder, &delimiterCount, &escapeChar)
		if err != nil {
			logger.Error(err)
			continue
		}

		if !hasNew {
			continue
		}

		handleIncomingMessage(conn, message)
	}
}
