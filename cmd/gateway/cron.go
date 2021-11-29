package main

import (
// importing the libraries again
// the encoding one is for the same pupose as previously
// the encoding/json package implements the encoding and decoding of the JSON (probably for the transfer
// of data which is handled in JSON format). The Marshal function is probably used to convert the
// JSON values to Go values and the Unmarshal function is probably used to convert the Go values back into
// the JSON values (the encoding and decoding, hence the name?)
	"encoding/json"
	// for the error handling most probably
	"errors"
	// this fmt implements formatted I/O with functions analogous to C's printf and scanf
	// this is from the official docs only (legit word to word)
	"fmt"
	// same as previously, for the file writing and stuff
	// the io/ioutil package implements some I/O utility functions
	"io/ioutil"
// the net/http package is for the Listening for connections probably
// it provides HTTP client and server implementations
// WOW MUCH SECURITY
	"net/http"
	// porbably for string manipulation
	// according to the official docs,  strings implements simple functions to manipulate
	// UTF-8 encoded strings.
	"strings"

	"github.com/fishjump/cs7ns1_project3/p2p-server/entities"
	// this is a contribution on Github I guess for cron in Go (however, I believe it
	// has also been included in the official docs and so can be trusted (what other option do we even have :P))
	// According to the official docs,
	// cron implements a cron spec parser and job runner.
	// Now, according to Wikipedia,
	// cron command-line utility, also known as cron job is a job scheduler
	//on Unix-like operating systems.
	//Users who set up and maintain software environments use cron to schedule jobs(commands or shell scripts)
	// to run periodically at fixed times, dates, or intervals
	// (basically scheduling)
	"github.com/robfig/cron/v3"
)

var (
	c *cron.Cron
)

func init() {
	c = cron.New()
// "the scheduling"
// calling certain functions every 1 minute, I believe the main use of the cron library
	c.AddFunc("@every 1m", detectExternelHearbeat)
	c.AddFunc("@every 1m", detectInternelHearbeat)
	c.AddFunc("@every 1m", fetchExternalNodes)
}
// this function is being used in detectExternelHearbeat and detectInternelHearbeat
// is basically just returning a URL when passed the name, port and the path
// use of fmt here to format the passed parameters as URL
func getUrl(host string, port int, path string) string {
	return fmt.Sprintf("https://%s:%d%s", host, port, path)
}

func detectExternelHearbeat() {
	// getting the external nodes, but where are we defining the external variable ?????
	nodes := external.Record.GetNodes()
	// go through each node in the list and ask for health
	for _, name := range nodes {
		// healthz is in entities folder, have to see it, this is why entities was imported here
		// wait this healthz is in p2p-server
		// should that not be imported?????
		resp, err := client.Get(getUrl(name, 33000, "/healthz"))
		if err != nil {
			logger.Error(err)
		}
// so if the response is not what we are expecting, the node is dead
// and so we remove that node from the external nodes
		if resp.StatusCode != http.StatusOK {
			external.Record.RemoveByName(name)
		}
	}
}
// so for gateway node, when detecting external heartbeat, getting the
// status of other gateways
// and when internal heartbeat, getting status of the devices?????
// same as the external heatbeat but this is for internal
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
		// the below 3 functions are the longest, but should not be the hardest to understand
		// might see later
		// but thing is that the same thing is there in the device and sensor and so
		// once this is understood then should not be much of a problem
		// however still have to look at all the files in the p2p-server folder
		// like healthz, list, message, register, server, unregister
		register(name, externalPort, externalHostName)
		externalMessage(name, externalPort, externalHostName)
		externalUpdateList(name, externalPort, externalHostName)
	}
}
// very long function, will take much time
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
	gwData, err := json.Marshal(gatewayDataMap[externalHostName])
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
