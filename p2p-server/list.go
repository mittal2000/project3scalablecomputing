package p2pserver

import (
	// all these libraries have been used before
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fishjump/cs7ns1_project3/p2p-server/entities"
)
// the functions are only longer because of the error handling
// this function is for the listing of all the available nodes on the server
func (s *P2PServer) list(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error(err)
	}
// making a variable of type ListRequest from the entities file
	request := &entities.ListRequest{}
	// unwrapping the body into request
	json.Unmarshal(body, request)
// s probably has record of all gateway nodes
// and if we do not find the gateway node, then return
// all sorts of errors
// but if we do find, then we return the list of devices that the gateway
// node has
	if _, find := s.Record.FindByToken(request.Token); !find {
		rw.WriteHeader(http.StatusBadRequest)

		response := entities.ListResponse{
			Status: false,
			Reason: "Wrong Token",
		}
// this is just to wrap up the response in the proper format to be written
		jsonStr, err := json.Marshal(response)
		if err != nil {
			logger.Error(err)
		}

		rw.Write(jsonStr)
	} else {
		// if we are able to find then we return an OK
		rw.WriteHeader(http.StatusOK)
// so basically we are listing all the nodes
		response := entities.ListResponse{
			Status: true,
			Reason: "",
			Nodes:  s.Record.GetNodes(),
		}

		jsonStr, err := json.Marshal(response)
		if err != nil {
			logger.Error(err)
		}

		rw.Write(jsonStr)
	}
}
// now message, register, serve, unregister are left alongwith that last
// very long fucntion in cron
// also scripts and the demoCA
// still not very clear about the certificates and stuff 
