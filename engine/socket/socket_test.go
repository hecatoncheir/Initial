package socket

import (
	"sync"
	"testing"

	"fmt"

	"github.com/hecatoncheir/Broker"
	"github.com/hecatoncheir/Configuration"
	"golang.org/x/net/websocket"
)

var (
	once       sync.Once
	goroutines sync.WaitGroup
)

func SetUpSocketServer() {
	testServer := New("v1.0", "", nil, nil)
	goroutines.Done()
	config := configuration.New()
	err := testServer.SetUp(config.Development.SocketServer.Host, config.Development.SocketServer.Port)
	if err != nil {
		fmt.Println("SetUpSocketServer faild with: ", err)
	}

	defer testServer.HTTPServer.Close()
}

func TestSocketServerCanHandleEvents(test *testing.T) {
	goroutines.Add(1)
	go once.Do(SetUpSocketServer)
	goroutines.Wait()

	config := configuration.New()
	if config.ServiceName == "" {
		config.ServiceName = "Initial"
	}

	iriOfWebSocketServer := fmt.Sprintf("ws://%v:%v", config.Development.SocketServer.Host,
		config.Development.SocketServer.Port)
	iriOfHTTPServer := fmt.Sprintf("http://%v:%v", config.Development.SocketServer.Host,
		config.Development.SocketServer.Port)

	socketConnection, err := websocket.Dial(iriOfWebSocketServer, "", iriOfHTTPServer)
	if err != nil {
		test.Error(err)
	}

	inputMessage := make(chan broker.EventData)

	go func() {
		defer socketConnection.Close()
		defer close(inputMessage)

		for {
			messageFromServer := broker.EventData{}
			err = websocket.JSON.Receive(socketConnection, &messageFromServer)
			if err != nil {
				test.Error(err)
				break
			}

			inputMessage <- messageFromServer
		}
	}()

	messageToServer := broker.EventData{Message: "Need api version"}
	err = websocket.JSON.Send(socketConnection, messageToServer)

	if err != nil {
		test.Error(err)
	}

	for messageFromServer := range inputMessage {
		if messageFromServer.Message == "Version of API" {
			break
		}

		if messageFromServer.Message != "Version of API" {
			test.Fail()
		}
	}
}
