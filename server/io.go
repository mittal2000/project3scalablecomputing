package server

import (
	"cs7ns1/project3/entities"
)

var (
	Inbox entities.Messages
	// outbox entities.Messages
	// requests entities.Messages
)

// func AddRequest(message entities.Message) {
// 	requests.Push(message)
// }

// func Send(message entities.Message) {
// 	outbox.Push(message)
// }

// func Read() (entities.Message, error) {
// 	return inbox.Pop()
// }
