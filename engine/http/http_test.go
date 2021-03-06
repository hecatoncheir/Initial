package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	"github.com/hecatoncheir/Configuration"
)

var (
	once       sync.Once
	goroutines sync.WaitGroup
)

func SetUpServer() {
	server := New("v1", nil)
	goroutines.Done()
	config := configuration.New()
	if config.ServiceName == "" {
		config.ServiceName = "Initial"
	}

	err := server.SetUp("", config.Development.HTTPServer.Host, config.Development.HTTPServer.Port)
	if err != nil {
		server.Log.Printf("Faild SetUpServer with error: %v", err)
	}
}

func TestHttpServerCanSendVersionOfAPI(test *testing.T) {
	goroutines.Add(1)
	go once.Do(SetUpServer)
	goroutines.Wait()

	config := configuration.New()
	if config.ServiceName == "" {
		config.ServiceName = "Initial"
	}

	iri := fmt.Sprintf("http://%v:%v/api/version", config.Development.HTTPServer.Host, config.Development.HTTPServer.Port)
	respose, err := http.Get(iri)
	if err != nil {
		test.Fatal(err)
	}

	encodedBody, err := ioutil.ReadAll(respose.Body)
	if err != nil {
		test.Fatal(err)
	}

	decodedBody := map[string]string{}

	err = json.Unmarshal(encodedBody, &decodedBody)
	if err != nil {
		test.Fatal(err)
	}

	if decodedBody["apiVersion"] != "v1" {
		fmt.Println("The api version should be the same.")
		test.Fail()
	}
}
