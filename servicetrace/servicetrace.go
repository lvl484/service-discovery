package servicetrace

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	ServiceName = "SERVICE_NAME"

	ConnIdleTimeout = 30 * time.Second

	ErrNotFound = "NOT_FOUND"
)

type Service struct {
	Alive    bool
	DeadLine time.Time
}

type Services struct {
	Mu         *sync.RWMutex
	ServiceMap map[string]Service
}

func (s *Service) CheckDead() bool {
	return s.DeadLine.After(time.Now())
}

func (s *Services) GetListOfServices() []string {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	services := make([]string, 0, len(s.ServiceMap))

	for name, service := range s.ServiceMap {
		serviceFormat := fmt.Sprintf("Name: %v, Alive: %v, Deadline: %v", name, service.Alive, service.DeadLine)
		services = append(services, serviceFormat)
	}

	return services
}

func (s *Services) UpSet(name string, service Service) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.ServiceMap[name] = service
}

func (s *Services) SetDeadLine(name string) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	service, ok := s.ServiceMap[name]

	if !ok {
		return errors.New(ErrNotFound)
	}

	service.DeadLine = time.Now().Add(ConnIdleTimeout)
	service.Alive = true

	s.ServiceMap[name] = service

	return nil
}

func (s *Services) SearchDead() {
	s.Mu.Lock()

	for name, service := range s.ServiceMap {
		if !service.Alive {
			continue
		}

		if !service.CheckDead() {
			service.Alive = false
			s.ServiceMap[name] = service
		}
	}

	s.Mu.Unlock()
}

func NewServices() *Services {
	return &Services{
		ServiceMap: make(map[string]Service),
		Mu:         new(sync.RWMutex),
	}
}
