package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type somedata struct {
	Age       int    `json:"Age"`
	Dead      bool   `json:"Dead"`
	OtherData string `json:"OtherData"`
}

type Data struct {
	sync.Mutex
	alldata map[string]somedata
}

func newData() *Data {
	return &Data{
		alldata: make(map[string]somedata),
	}
}

func (d *Data) Add(w http.ResponseWriter, r *http.Request) {
	var sd somedata
	err := json.NewDecoder(r.Body).Decode(&sd)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	d.Lock()
	d.alldata[id] = sd
	d.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func (d *Data) Get(w http.ResponseWriter, r *http.Request) {
	var sd somedata

	params := mux.Vars(r)
	sd = d.alldata[params["ID"]]
	err := json.NewEncoder(w).Encode(sd)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (d *Data) GetIDs(w http.ResponseWriter, r *http.Request) {
	keys := make([]string, 0, len(d.alldata))
	for k := range d.alldata {
		keys = append(keys, k)
	}

	err := json.NewEncoder(w).Encode(keys)

	if err != nil {
		return
	}
}

func (d *Data) Update(w http.ResponseWriter, r *http.Request) {
	var sd somedata

	params := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&sd)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d.Lock()
	d.alldata[params["ID"]] = sd
	d.Unlock()
}

func (d *Data) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	d.Lock()
	delete(d.alldata, params["ID"])
	d.Unlock()
	w.WriteHeader(http.StatusNoContent)
}
