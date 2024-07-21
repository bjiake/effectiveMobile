package interfaces

import (
	"context"
	"effectiveMobile/pkg/domain/task"
	"time"
)

type TaskRepository interface {
	Migrate(ctx context.Context) error
	Post(ctx context.Context, newTask task.Task) (*task.Task, error)
	Put(ctx context.Context, id int64, updateTask task.Task) (*task.Task, error)
	Get(ctx context.Context, id int64) (*task.Task, error)
	GetLaborCost(ctx context.Context, startTime time.Time, endTime time.Time) (task.Slice, error)
	GetAll(ctx context.Context) ([]task.Task, error)
	Delete(ctx context.Context, id int64) error
}
