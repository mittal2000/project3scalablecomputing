package entities
// this entire file, as the name suggests, deals with the records in the sense that
// it manages them, their addition, removal (by token and by name), getting the list of tokens, nodes
// pretty important file, this contains all the control functions you could say; some of these can be
// used by the devices as well such as getNodes but they cannot add or remove Nodes as that
// control should only be given to the gateway node; FindByName and FindByToken could also
// be used by the devices as they would need that information to communicate with one another
import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

type (
	// declaring different variables simply
	// NodeRecord and TokenRecord are simple string to string maps
	NodeRecord  map[string]string
	TokenRecord map[string]string
// defining a new data type with 3 variables, 2 string to string maps and
// one mutex object
// according to the docs, it is a mutual exclusion lock
// its 2 main functions being used in this file are lock and Unlock
// As the name suggests, Lock locks m. If the lock is already in use,
//the calling goroutine blocks until the mutex is available (this is all from the official docs)
// and unlock unlocks m. It is a run-time error if m is not locked on entry to Unlock.
// we might have to do the error handling for unlock (not sure whether required or not)
	NodeTokenRecord struct {
		nodeRecord  NodeRecord
		tokenRecord TokenRecord
		m           sync.Mutex
	}
)
// generating a unique token based on the name of the device?????
func GenToken(name string) string {
	token := sha256.Sum256([]byte(name + time.Now().String()))

	return fmt.Sprintf("%x", token)
}
// func parameters_in function_name parameters_returned_type
// below function, no parameters are being passed in and so, nothing before function name
// however, r being returned, of type NodeTokenRecord and so the return type after function name
// is *NodeTokenRecord
func CreateNodeTokenRecord() *NodeTokenRecord {
	r := &NodeTokenRecord{}
	r.init()

	return r
}

func (r *NodeTokenRecord) init() {
	r.nodeRecord = make(NodeRecord)
	r.tokenRecord = make(TokenRecord)
}
// as the name suggests, this function is for getting the list of nodes (nodeRecord)
func (r *NodeTokenRecord) GetNodes() []string {

	r.m.Lock()
	// according to https://go.dev/tour/flowcontrol/12
	// A defer statement defers the execution of a function
	//  until the surrounding function returns.
	defer r.m.Unlock()

	var list []string
	for k := range r.nodeRecord {
		list = append(list, k)
	}

	return list
}

func (r *NodeTokenRecord) GetTokens() []string {
	r.m.Lock()
	defer r.m.Unlock()

	var list []string
	for k := range r.tokenRecord {
		list = append(list, k)
	}

	return list
}
// adding in names and tokens into the tokenRecord and nodeRecord maps
func (r *NodeTokenRecord) Add(token string, name string) {
	r.m.Lock()
	defer r.m.Unlock()

	r.tokenRecord[token] = name
	r.nodeRecord[name] = token
}

func (r *NodeTokenRecord) RemoveByToken(token string) {
	r.m.Lock()
	defer r.m.Unlock()
// erm so basically, get the token into name and the name into find, if find is not nil
// then go ahead and delete the record from the tokenRecord map and the nodeRecord map
	if name, find := r.tokenRecord[token]; find {
		delete(r.tokenRecord, token)
		delete(r.nodeRecord, name)
	}
}
// the same thing as above, just this time using the name to find the record rather than the token
// which makes sense as in the system, the gateway node might want to delete a node by the name
// this also explains why there are the redundant tokenRecord and nodeRecord, so that these
// 2 different operations can be carried out a bit more cleanly I guess?????
func (r *NodeTokenRecord) RemoveByName(name string) {
	r.m.Lock()
	defer r.m.Unlock()

	if token, find := r.nodeRecord[name]; find {
		delete(r.tokenRecord, token)
		delete(r.nodeRecord, name)
	}
}
// pretty self explanatory, lock r again, then find the record by the tokenRecord data and then
// return it
func (r *NodeTokenRecord) FindByToken(token string) (string, bool) {
	r.m.Lock()
	defer r.m.Unlock()

	val, find := r.tokenRecord[token]

	return val, find
}
// pretty self explanatory, lock r again, then find the record by the nodeRecord data and then
// return it
func (r *NodeTokenRecord) FindByName(name string) (string, bool) {
	r.m.Lock()
	defer r.m.Unlock()

	val, find := r.nodeRecord[name]

	return val, find
}
