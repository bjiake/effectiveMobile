package interfaces

import (
	"context"
	"effectiveMobile/pkg/domain/people"
	"effectiveMobile/pkg/domain/task"
)

type ServiceUseCase interface {
	Migrate(ctx context.Context) error

	// People
	InfoPeople(ctx context.Context, passportSerie string, passportNumber string) (*people.Info, error)
	Registration(ctx context.Context, newPeople people.Registration) (*int64, error)
	Login(ctx context.Context, people people.Registration) (int64, error)
	GetPeople(ctx context.Context, filter *people.Filter, pagination *people.Pagination) ([]people.Request, error)
	PutPeople(ctx context.Context, id string, updatePeople people.Info) (*people.Info, error)
	DeletePeople(ctx context.Context, id string) error

	//Task
	TaskStart(ctx context.Context, id string, newTask task.Task) (*task.Task, error)
	TaskFinish(ctx context.Context, taskId string) (*task.Task, error)
	TaskPut(ctx context.Context, id string, updateTask task.Task) (*task.Task, error)
	GetTask(ctx context.Context, startTimeStr string, endTimeStr string) ([]task.Task, error)
	GetAllTask(ctx context.Context) ([]task.Task, error)
	DeleteTask(ctx context.Context, id string) error
}
