package cli

import (
	"cs7ns1/project3/server"
	"fmt"
)

func ls() {
	fmt.Print(server.GetIndex("\n"))
}
