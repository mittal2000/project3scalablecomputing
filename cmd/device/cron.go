package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/fishjump/cs7ns1_project3/p2p-server/entities"
	"github.com/robfig/cron/v3"
)

var (
	c *cron.Cron
)

func init() {
	c = cron.New()

	c.AddFunc("@every 1m", detectExternelHearbeat)
	c.AddFunc("@every 1m", detectInternelHearbeat)
	c.AddFunc("@every 1m", fetchExternalNodes)
}

func getUrl(host string, port int, path string) string {
	return fmt.Sprintf("https://%s:%d%s", host, port, path)
}

func detectExternelHearbeat() {
	nodes := external.Record.GetNodes()
	for _, name := range nodes {
		resp, err := client.Get(getUrl(name, 33000, "/healthz"))
		if err != nil {
			logger.Error(err)
		}

		if resp.StatusCode != http.StatusOK {
			external.Record.RemoveByName(name)
		}
	}
}

func detectInternelHearbeat() {
	nodes := internal.Record.GetNodes()
	for _, name := range nodes {
		resp, err := client.Get(getUrl(name, 443, "/healthz"))
		if err != nil {
			logger.Error(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			internal.Record.RemoveByName(name)
		}
	}
}

func fetchExternalNodes() {
	nodes := external.Record.GetNodes()
	for _, name := range nodes {
		register(name, externalPort, externalHostName)
		externalMessage(name, externalPort, externalHostName)
		externalUpdateList(name, externalPort, externalHostName)
	}
}

func externalUpdateList(host string, port int, local string) {
	request := &entities.ListRequest{
		Token: clientToken[host],
	}

	data, err := json.Marshal(request)
	if err != nil {
		logger.Error(err)
		return
	}

	resp, err := client.Post(getUrl(host, port, "/list"),
		"application/json", strings.NewReader(string(data)))
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()

	response := &entities.ListResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	json.Unmarshal(body, response)

	if resp.StatusCode != http.StatusOK {
		logger.Error(errors.New(response.Reason))
		return
	}

	if !response.Status {
		logger.Error(errors.New(response.Reason))
		return
	}

	for _, name := range response.Nodes {
		if _, find := external.Record.FindByName(name); !find {
			external.Record.Add(entities.GenToken(name), name)
		}
	}
}

func externalMessage(host string, port int, local string) {
	gwData, err := json.Marshal(deviceDataMap[externalHostName])
	if err != nil {
		logger.Error(err)
	}

	request := &entities.MessageRequest{
		Token: clientToken[host],
		Data:  string(gwData),
	}

	data, err := json.Marshal(request)
	if err != nil {
		logger.Error(err)
	}

	resp, err := client.Post(getUrl(host, port, "/message"),
		"application/json", strings.NewReader(string(data)))
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()

	response := &entities.MessageResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	json.Unmarshal(body, response)

	if resp.StatusCode != http.StatusOK {
		logger.Error(errors.New(response.Reason))
		return
	}

	if !response.Status {
		logger.Error(errors.New(response.Reason))
		return
	}
}

func register(host string, port int, local string) {
	request := &entities.RegisterRequest{
		Name: local,
	}

	data, err := json.Marshal(request)
	if err != nil {
		logger.Error(err)
		return
	}

	resp, err := client.Post(getUrl(host, port, "/register"),
		"application/json", strings.NewReader(string(data)))
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()

	response := &entities.RegisterResponse{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	json.Unmarshal(body, response)

	if resp.StatusCode != http.StatusOK {
		logger.Error(errors.New(response.Reason))
		return
	}

	if !response.Status {
		logger.Error(errors.New(response.Reason))
		return
	}

	clientToken[host] = response.Token
}
