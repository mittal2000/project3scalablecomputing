package cli

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"cs7ns1/project3/entities"
)

type (
	Cmd func()
)

var (
	banner string
	help   string
	cmd    map[string]Cmd
)

func init() {
	banner = `
	███████╗ ██████╗    ██████╗ ██████╗  ██████╗      ██╗███████╗ ██████╗████████╗██████╗ 
	██╔════╝██╔════╝    ██╔══██╗██╔══██╗██╔═══██╗     ██║██╔════╝██╔════╝╚══██╔══╝╚════██╗
	███████╗██║         ██████╔╝██████╔╝██║   ██║     ██║█████╗  ██║        ██║    █████╔╝
	╚════██║██║         ██╔═══╝ ██╔══██╗██║   ██║██   ██║██╔══╝  ██║        ██║    ╚═══██╗
	███████║╚██████╗    ██║     ██║  ██║╚██████╔╝╚█████╔╝███████╗╚██████╗   ██║   ██████╔╝
	╚══════╝ ╚═════╝    ╚═╝     ╚═╝  ╚═╝ ╚═════╝  ╚════╝ ╚══════╝ ╚═════╝   ╚═╝   ╚═════╝ 
																						  
	██████╗ ██████╗ ██████╗      ██████╗██╗     ██╗███████╗███╗   ██╗████████╗            
	██╔══██╗╚════██╗██╔══██╗    ██╔════╝██║     ██║██╔════╝████╗  ██║╚══██╔══╝            
	██████╔╝ █████╔╝██████╔╝    ██║     ██║     ██║█████╗  ██╔██╗ ██║   ██║               
	██╔═══╝ ██╔═══╝ ██╔═══╝     ██║     ██║     ██║██╔══╝  ██║╚██╗██║   ██║               
	██║     ███████╗██║         ╚██████╗███████╗██║███████╗██║ ╚████║   ██║               
	╚═╝     ╚══════╝╚═╝          ╚═════╝╚══════╝╚═╝╚══════╝╚═╝  ╚═══╝   ╚═╝         
	`

	help = `Useage:
ls                          - List all known nodes
send host message [ttl]     - Send a message to host, find host recusively, default ttl is 255
fetch                       - Fetch avaliable nodes from remote
read                        - Read new messages
help                        - Show this again
exit                        - Exit
`

	cmd = make(map[string]Cmd)
	cmd["ls"] = ls
	cmd["send"] = send
	cmd["read"] = read
	cmd["help"] = func() {
		fmt.Print(help)
	}
	cmd["exit"] = func() {
		os.Exit(0)
	}
	cmd[""] = func() {
		// Do nothing
	}
}

func Run(wg *sync.WaitGroup, config *entities.Config) {
	defer wg.Done()

	fmt.Println(banner)
	fmt.Println(help)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		cmdStr := scanner.Text()
		// TODO: BUG FIX, send host ****
		if fn, find := cmd[cmdStr]; find {
			fn()
		} else {
			fmt.Println("Error command input, please see help")
			cmd["help"]()
		}

		if scanner.Err() != nil {
			panic(scanner.Err().Error())
		}
	}
}
