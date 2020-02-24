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
	mu         *sync.RWMutex
	ServiceMap map[string]Service
}

func (s *Service) CheckDead() bool {
	return s.DeadLine.After(time.Now())
}

func (s *Services) GetListOfServices() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	services := make([]string, 0, len(s.ServiceMap))

	for name, service := range s.ServiceMap {
		serviceFormat := fmt.Sprintf("Name: %v, Alive: %v, Deadline: %v", name, service.Alive, service.DeadLine)
		services = append(services, serviceFormat)
	}

	return services
}

func (s *Services) UpSet(name string, service Service) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ServiceMap[name] = service
}

func (s *Services) SetDeadLine(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
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
	s.mu.Lock()
	defer s.mu.Unlock()

	for name, service := range s.ServiceMap {
		if !service.CheckDead() {
			service.Alive = false
			s.ServiceMap[name] = service
		}
	}
}

func NewServices() *Services {
	return &Services{
		ServiceMap: make(map[string]Service),
		mu:         new(sync.RWMutex),
	}
}
