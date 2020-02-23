package servicetrace

import (
	"errors"
	"log"
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
	return s.DeadLine.Before(time.Now())
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

func (s *Services) SearchDead(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			s.Mu.Lock()

			for name, service := range s.ServiceMap {
				log.Println(1)
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
	}
}

func NewServices() *Services {
	return &Services{
		ServiceMap: make(map[string]Service),
		Mu:         new(sync.RWMutex),
	}
}
