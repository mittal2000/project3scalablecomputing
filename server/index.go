package server

import (
	"cs7ns1/project3/entities"
)

var (
	hostMap entities.HostMap
)

func init() {
	hostMap = entities.NewHostMap()
}

func GetIndex(delimiter string) string {
	return hostMap.String(delimiter)
}
