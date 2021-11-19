package entities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type (
	Config struct {
		Host    string  `json:"host"`
		Port    int     `json:"port"`
		Indexes Indexes `json:"indexes"`
	}

	Index struct {
		Url  string `json:"url"`
		Port int    `json:"port"`
	}

	Indexes []Index
)

func newIndex(url string, port int) *Index {
	return &Index{
		Url:  url,
		Port: port,
	}
}

func (c *Config) String() string {
	res, err := json.Marshal(c)
	if err != nil {
		return ""
	}

	return string(res)
}

func (c *Config) Set(value string) error {
	file, err := os.Open(value)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(jsonBytes), c)
}

func (i *Index) String() string {
	return fmt.Sprintf("%s:%d", i.Url, i.Port)
}

func (i *Indexes) String() string {
	var tmp []string
	for _, item := range *i {
		tmp = append(tmp, item.String())
	}

	return strings.Join(tmp, ";")
}

func (i *Indexes) Set(value string) error {
	regex := regexp.MustCompile("^([A-Za-z0-9.]+):([0-9]+)$")
	matches := regex.FindStringSubmatch(value)

	if len(matches) != 3 {
		return fmt.Errorf("require url format like host:port")
	}

	url := matches[1]
	port, err := strconv.Atoi(matches[2])
	if err != nil {
		return err
	}

	*i = append(*i, *newIndex(url, port))
	return nil
}
