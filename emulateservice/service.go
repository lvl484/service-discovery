package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"syreclabs.com/go/faker"
)

const (
	DefaultURL = "http://localhost:8080/list"

	SleepTime = 10 * time.Second
)

type config struct {
	Name        string    `json:"Name"`
	Author      string    `json:"Author"`
	DateCreated time.Time `json:"DateCreated"`
	ConfigData  []byte    `json:"ConfigData"`
}

func main() {
	configID := flag.String("i", "SomeID", "id to get config")
	emulatorNum := flag.Int("n", 1, "nubmer of services")
	flag.Parse()

	wg := &sync.WaitGroup{}
	Emulate(*configID, *emulatorNum, wg)
	wg.Wait()
}

func GetConfigs() []string {
	var idList []string

	resp, err := http.Get(DefaultURL)

	if err != nil {
		log.Println(err)
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(&idList)

	if err != nil {
		log.Println(err)
		return nil
	}

	return idList
}

func Emulate(configID string, num int, wg *sync.WaitGroup) {
	for i := 0; i < num; i++ {
		wg.Add(1)

		serviceName := faker.Lorem().Word()
		fullURL := fmt.Sprintf("%v/%v/%v", DefaultURL, configID, serviceName)
		log.Println(fullURL)

		go emulate(fullURL, wg)
	}
}

func emulate(configURL string, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		var c config

		resp, err := http.Get(configURL)

		if err != nil {
			log.Println(err)
		}

		err = json.NewDecoder(resp.Body).Decode(&c)

		if err != nil {
			log.Println(err)
		}

		log.Println(c)

		time.Sleep(SleepTime)

		resp.Body.Close()
	}
}
