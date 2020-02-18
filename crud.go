package main

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type config struct {
	ID             int    `json:"ID"`
	IP             string `json:"IP"`
	Port           bool   `json:"Port"`
	MaxConnections int    `json:"MaxConnections"`
	OtherSettings  string `json:"OtherSettings"`
}

type configs struct {
	mu  *sync.RWMutex
	all map[string]config
}

type Data struct {
	*configs
}

func NewData() *Data {
	return &Data{
		&configs{
			all: make(map[string]config),
			mu:  new(sync.RWMutex),
		},
	}
}

// Add config
func (d *Data) Add(w http.ResponseWriter, r *http.Request) {
	var c config
	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	d.configs.add(id, c)

	w.WriteHeader(http.StatusCreated)
}

// Get config by ID
func (d *Data) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	c := d.configs.get(params["ID"])
	err := json.NewEncoder(w).Encode(c)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// GetIDs returns slice of IDs
func (d *Data) GetIDs(w http.ResponseWriter, r *http.Request) {
	keys := d.configs.getIDs()
	err := json.NewEncoder(w).Encode(keys)

	if err != nil {
		return
	}
}

// Update config by ID
func (d *Data) Update(w http.ResponseWriter, r *http.Request) {
	var c config

	params := mux.Vars(r)
	err := json.NewDecoder(r.Body).Decode(&c)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d.configs.set(params["ID"], c)
}

// Delete config by ID
func (d *Data) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	d.configs.delete(params["ID"])
	w.WriteHeader(http.StatusNoContent)
}

func (c *configs) set(id string, conf config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.all[id]; ok {
		c.all[id] = conf
	}
}

func (c *configs) add(id string, conf config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.all[id]; ok {
		return
	}

	c.all[id] = conf
}

func (c *configs) get(id string) config {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.all[id]
}

func (c *configs) getIDs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]string, 0, len(c.all))

	for k := range c.all {
		keys = append(keys, k)
	}

	return keys
}

func (c *configs) delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.all, id)
}
