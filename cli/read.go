package cli

import (
	"cs7ns1/project3/server"
	"fmt"
)

func read() {
	for first := true; true; {
		message, err := server.Inbox.Pop()
		if err != nil {
			if first {
				fmt.Println("No new message")
			}
			return
		}
		first = false
		fmt.Println("---------------")
		fmt.Printf("From: %s:%d\n", message.FromIP, message.FromPort)
		fmt.Println("---------------")
		fmt.Printf("%s\n", message.Body)
		fmt.Println("---------------")

	}
}
