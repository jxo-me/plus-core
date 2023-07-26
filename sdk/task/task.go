package task

import (
	"context"
	"github.com/jxo-me/plus-core/core/task"
)

const (
	SrvName = "TaskService"
)

var insService = tService{
	Services: []task.IService{},
}

type tService struct {
	Services []task.IService
}

func Services() *tService {
	return &insService
}

func (t *tService) String() string {
	return SrvName
}

func (t *tService) AddServices(services ...task.IService) task.TasksService {
	t.Services = services
	return t
}

func (t *tService) Start(ctx context.Context) {
	for _, service := range t.Services {
		service.Start(ctx)
	}
}
