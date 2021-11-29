package p2pserver

import (
	// for the server listening and stuff
	"net/http"
)
// little function, if the node would be alive
// then it will be able to execute this
// and whosoever has called this, will be getting a Response
// this function is getting used in the cron file -> very complex

func (s *P2PServer) healthz(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	rw.WriteHeader(http.StatusOK)
}
