package task

import (
	"time"
)

// Task represents a task with its details.
// @swagger:model
type Task struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description"`
	StartTime   time.Time      `json:"startTime" validate:"required"`
	EndTime     *time.Time     `json:"endTime"`
	TotalTime   *time.Duration `json:"totalTime"`
}

type Slice []Task

func (s Slice) Len() int      { return len(s) }
func (s Slice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Slice) Less(i, j int) bool {
	// Если TotalTime не установлено, считаем его наименьшим
	if s[i].TotalTime == nil {
		return false
	}
	if s[j].TotalTime == nil {
		return true
	}
	// Сортировка от большего к меньшему
	return *s[i].TotalTime > *s[j].TotalTime
}
