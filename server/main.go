package server

import (
	"fmt"
	"net"
	"os"
	"sync"

	"cs7ns1/project3/entities"

	"github.com/withmandala/go-log"
)

func Run(wg *sync.WaitGroup, config *entities.Config) {
	defer wg.Done()

	logger := log.New(os.Stderr)
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP(config.Host), Port: config.Port})
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Info(fmt.Sprintf("Local: <%s> \n", listener.Addr().String()))

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			logger.Error(err)
		}

		go handleTCPConnection(conn)
	}
}
