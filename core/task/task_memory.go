package task

type MemorySpec struct {
	TaskName   string
	RoutingKey string
}

type MemoryTask interface {
	GetSpec() *MemorySpec
	IHandler
}

type MemoryService interface {
	IService
	AddTasks(task ...MemoryTask) MemoryService
}
