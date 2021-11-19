package main

import (
	"flag"
	"os"
	"sync"

	"cs7ns1/project3/cli"
	"cs7ns1/project3/entities"
	"cs7ns1/project3/server"

	"github.com/withmandala/go-log"
)

type (
	Worker func(wg *sync.WaitGroup, config *entities.Config)
)

var (
	config  entities.Config
	workers map[string]Worker
)

func parseParams() {
	flag.StringVar(&config.Host, "host", "127.0.0.1", "localhost")
	flag.IntVar(&config.Port, "port", 8888, "port")
	flag.Var(&config.Indexes, "indexs", "known index list")
	flag.Var(&config, "config", "init use json config file")
	flag.Parse()
}

func init() {
	parseParams()
	workers = make(map[string]Worker)
	workers["cli"] = cli.Run
	workers["server"] = server.Run
}

func main() {
	logger := log.New(os.Stderr)
	logger.Info("Initialized with config:", config.String())

	wg := &sync.WaitGroup{}

	for _, worker := range workers {
		wg.Add(1)
		go worker(wg, &config)
	}

	wg.Wait()
}
