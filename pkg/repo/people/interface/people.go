package interfaces

import (
	"context"
	people "effectiveMobile/pkg/domain/people"
	"effectiveMobile/pkg/domain/task"
)

type PeopleRepository interface {
	Migrate(ctx context.Context) error
	Info(ctx context.Context, passportNumber string) (*people.Info, error)
	Registration(ctx context.Context, newPeople people.Registration) (*int64, error)
	Login(ctx context.Context, acc people.Registration) (int64, error)
	Get(ctx context.Context, filter *people.Filter, pagination *people.Pagination) ([]people.Request, error)
	Put(ctx context.Context, id int64, updatePeople people.Info) (*people.Info, error)
	Delete(ctx context.Context, id int64) error
	AppendTask(ctx context.Context, id int64, task task.Task) error
}
