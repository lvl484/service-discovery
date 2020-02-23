package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lvl484/service-discovery/servicetrace"
)

func (d *Data) GetForService(s *servicetrace.Services) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		c, _ := d.configs.get(params[muxVarsID])
		err := json.NewEncoder(w).Encode(c)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		serviceName := params[servicetrace.ServiceName]

		var service servicetrace.Service

		s.UpSet(serviceName, service)
		if err := s.SetDeadLine(serviceName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func GetListOfServices(s *servicetrace.Services) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Mu.RLock()
		defer s.Mu.RUnlock()
		for name, service := range s.ServiceMap {
			serviceFormat := fmt.Sprintf("Name: %v, Alive: %v", name, service.Alive)
			err := json.NewEncoder(w).Encode(serviceFormat)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	})
}
