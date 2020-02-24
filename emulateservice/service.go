package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
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
	emulatorNum := flag.Int("n", 1, "nubmer of services")
	flag.Parse()

	wg := &sync.WaitGroup{}
	if err := Emulate(*emulatorNum, wg); err != nil {
		log.Println(err)
	}
	wg.Wait()
}

func GetConfigs() ([]string, error) {
	var idList []string

	resp, err := http.Get(DefaultURL)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&idList)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return idList, nil
}

func Emulate(num int, wg *sync.WaitGroup) error {
	s, err := GetConfigs()

	if err != nil {
		log.Println(err)
		return err
	}

	if len(s) < 1 {
		return errors.New("EMPTY_SLICE")
	}

	for i := 0; i < num; i++ {
		wg.Add(1)

		configID := s[rand.Intn(len(s))]
		serviceName := faker.Lorem().Word()
		fullURL := fmt.Sprintf("%v/%v/%v", DefaultURL, configID, serviceName)
		log.Println(fullURL)

		go emulate(fullURL, wg)
	}

	return nil
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
