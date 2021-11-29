package main
// this is just a comment to test something out
// it is necessary for a program in go to have package main as the first line
// importing the necessary packages
import (
// the crypto/tls package is for the TLS protocol. TLS stands for Transport Layer Security and it is
// used for encrypting the communications between web servers and applications
	"crypto/tls"
// the crypto/x509 package is related to security and certificates
	"crypto/x509"
// the encoding/json package implements the encoding and decoding of the JSON (probably for the transfer
// of data which is handled in JSON format). The Marshal function is probably used to convert the
// JSON values to Go values and the Unmarshal function is probably used to convert the Go values back into
// the JSON values (the encoding and decoding, hence the name?)
	"encoding/json"
// the flag library is for command line flag parsing
// WHAT IS COMMAND LINE FLAG PARSING?????
// currently unsure what this is for, will update once have gone through the code
	"flag"
// the io/ioutil package implements some I/O utility functions
// WHAT ARE I/O UTILITY FUNCTIONS?????
	"io/ioutil"
// for generating random numbers
  "math/rand"
// the net/http package is for the Listening for connections probably
// it provides HTTP client and server implementations
// WOW MUCH SECURITY
	"net/http"
// the os package is for operating system functionality such as reading files
	"os"
// the sync package is for basic synchronisation like mutual exclusion locks
// might be useful for things like collision resolution or time sharig perhaps?????
	"sync"
// NOT VERY SURE WHAT THIS IS ?????
	p2pserver "github.com/fishjump/cs7ns1_project3/p2p-server"
	"github.com/fishjump/cs7ns1_project3/p2p-server/entities"
// The below import is for better log handling
	"github.com/withmandala/go-log"
	"strconv"
)

// declaring different variables here
// this language is very much like C
// the vaiables have been appropriately named to make it clear what purpose they serve
// in Go, the map[type]type is like a dictionary in Python
// Need to look at what is entities.GatewayData
// much to see, need to understand what all the other files are in for
//
var (
	dir              string
	externalHostName string
	internalHostName string
	initialIndexHost string
	externalPort     int
	internalPort     int

	clientToken    map[string]string
// for device might have to change this to deviceDataMap and entities.DeviceData
// similarly for sensors will have to change this to sensorDataMap and entities.SensorData
	sensorDataMap map[string]entities.SensorData

	external *p2pserver.P2PServer
	client   *http.Client

	wg sync.WaitGroup

	logger *log.Logger
)
// WHAT IS THE DIFFERENCE BETWEEN external message callback and internal message callback?????
func externalMsgCbk(name string, req *entities.MessageRequest) {
// making a GatewayData type variable
// This is containing 2 types
// One is a Name of type String and other is Data which is a map from a string to DeviceData
//  the string is probably the name of the device
// the DeviceData in itself has name and SensorData
// look at the entities.go file (p2p-server -> entities -> entities.go)
// that will make everything clear so for device, will have to change this from GatewayData to DeviceData
// and for sensors will have to change this to SensorData
	// data := &entities.SensorData{}
// the line below can be split into 2, err := json.Unmarshal([]byte(req.Data), data) this parses the JSON data
// and stores it into the GatewayData type object "data"
// and then if err != nil, we log the error (for this the last import was there)
// the json package's Unmarshal function is for
// func Unmarshal(data []byte, v interface{}); Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
//If v is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError.
//Unmarshal uses the inverse of the encodings that Marshal uses, allocating maps, slices, and pointers as necessary,
// with some additional rules which can be read about here -> https://pkg.go.dev/encoding/json#Unmarshal
	// if err := json.Unmarshal([]byte(req.Data), data); err != nil {
	// 	logger.Error(err)
	// 	return
	// }
// basically the gatewayDataMap stores details about the different gateways using mapping from the named
// of the gateway to the data which is a Gateway Data type
	// sensorDataMap[name] = *data
	max := 100
  	min := 10
  	sensor1 := rand.Intn(max - min) + min
	sensor2 := rand.Intn(max - min) + min
	sensor3 := rand.Intn(max - min) + min
	sensor4 := rand.Intn(max - min) + min
	sensor5 := rand.Intn(max - min) + min
	sensor6 := rand.Intn(max - min) + min
	sensor7 := rand.Intn(max - min) + min
	sensor8 := rand.Intn(max - min) + min
	sensorDataMap := entities.SensorData{Name: internalHostName, Data: map[string]string{ "sensor1": strconv.Itoa(sensor1), "sensor2": strconv.Itoa(sensor2), "sensor3": strconv.Itoa(sensor3), "sensor4": strconv.Itoa(sensor4), "sensor5": strconv.Itoa(sensor5), "sensor6": strconv.Itoa(sensor6), "sensor7": strconv.Itoa(sensor7), "sensor8": strconv.Itoa(sensor8)}}

// and now the map is encoded back into JSON
	str, err := json.Marshal(sensorDataMap)
	if err != nil {
		logger.Error(err)
	}
// using the encoded JSON from the deviceDataMap
// this might be there for deviceData but will not be there for sensorData I am guessing
	ioutil.WriteFile(dir+"/data.json", str, 0644)
}
// internalMsgCbk is not requried in Sensor if I understand correctly ?????
//func internalMsgCbk(name string, req *entities.MessageRequest) {
// in the internal message callback, data is DeviceData type rather than Gateway Data Type
// so this might mean that this is for the data which is coming into the server
// and the external message callback was for the data which is going out of the server
//	data := &entities.SensorData{}
//	if err := json.Unmarshal([]byte(req.Data), data); err != nil {
//		logger.Error(err)
//		return
//	}


//	deviceDataMap[externalHostName].Data[name] = *data

//	str, err := json.Marshal(deviceDataMap)
//	if err != nil {
//		logger.Error(err)
//	}

//	ioutil.WriteFile(dir+"/data.json", str, 0644)
//}

func runBackground(fn func()) {
// wg is a variable of the type sync.WaitGroup
// sync.WaitGroup A WaitGroup waits for a collection of goroutines to finish.
//The main goroutine calls Add to set the number of goroutines to wait for.
//Then each of the goroutines runs and calls Done when finished.
//At the same time, Wait can be used to block until all goroutines have finished.(this is copied from the
// official documentations)
	wg.Add(1)
	go func() {
		defer wg.Done()
		fn()
	}()
}

// now have to look at the init function again,
// ran out of brainpower
// but went through most of the code
// still nothing makes sense as to how it is working though haha
func init() {
	clientToken = make(map[string]string)
	
	
	logger := log.New(os.Stderr)

	flag.StringVar(&dir, "dir", ".", "directory to save data")
	flag.StringVar(&externalHostName, "host", "rasp-019.scss.tcd.ie", "")
	flag.IntVar(&externalPort, "port", 33000, "")
	flag.StringVar(&internalHostName, "subhost", "127.0.0.1", "")
	flag.IntVar(&internalPort, "subport", 443, "")
	flag.StringVar(&initialIndexHost, "index", "rasp-019.scss.tcd.ie", "")
	flag.Parse()
	// sensorDataMap := entities.SensorData{Name: internalHostName, Data: map[string]string{ "sensor1": strconv.Itoa(sensor1), "sensor2": strconv.Itoa(sensor2), "sensor3": strconv.Itoa(sensor3), "sensor4": strconv.Itoa(sensor4), "sensor5": strconv.Itoa(sensor5), "sensor6": strconv.Itoa(sensor6), "sensor7": strconv.Itoa(sensor7), "sensor8": strconv.Itoa(sensor8)}}
	external = p2pserver.NewServer(externalHostName, externalPort,
		dir+"/external.server.key",
		dir+"/external.server.crt",
		dir+"/ca.crt",
		externalMsgCbk)
// look at the file Record which is in the entities folder inside the p2p-server folder
	external.Record.Add(entities.GenToken(initialIndexHost), initialIndexHost)

	certPair, err := tls.LoadX509KeyPair(dir+"/client.crt", dir+"/client.key")
	if err != nil {
		logger.Error(err)
	}

	caCert, err := ioutil.ReadFile(dir + "/ca.crt")
	if err != nil {
		logger.Error(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certPair},
			},
		},
	}
}

func main() {
	runBackground(external.RunTLS)
	runBackground(c.Start)

	wg.Wait()
}
