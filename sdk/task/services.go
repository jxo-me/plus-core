package task

import (
	"context"
)

const (
	SrvName = "TaskService"
)

var insService = tService{
	Services: []Service{},
}

type tService struct {
	Services []Service
}

func Services() *tService {
	return &insService
}

func (t *tService) String() string {
	return SrvName
}

func (t *tService) AddTasks(services ...Service) Service {
	t.Services = services
	return t
}

func (t *tService) Start(ctx context.Context) {
	for _, service := range t.Services {
		service.Start(ctx)
	}
}
