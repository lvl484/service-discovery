package servicetrace

import (
	"errors"
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
	return s.DeadLine.Before(time.Now())
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

func (s *Services) SearchDead(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-ticker.C:
			for _, s := range s.ServiceMap {
				if !s.CheckDead() {
					s.Alive = false
				}
			}
		}
	}
}

func NewServices() *Services {
	return &Services{
		ServiceMap: make(map[string]Service),
		mu:         new(sync.RWMutex),
	}
}
