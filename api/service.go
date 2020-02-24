package main

import (
	"encoding/json"
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
		s.SearchDead()
		list := s.GetListOfServices()
		for _, service := range list {
			err := json.NewEncoder(w).Encode(service)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	})
}
