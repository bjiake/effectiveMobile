package service

import (
	"context"
	"effectiveMobile/pkg/domain/task"
	"fmt"
	"sort"
	"time"
)

func (s *service) TaskStart(ctx context.Context, id string, newTask task.Task) (*task.Task, error) {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return nil, err
	}

	newTask.StartTime = time.Now().UTC()
	result, err := s.rTask.Post(ctx, newTask)
	if err != nil {
		return nil, err
	}
	err = s.rPeople.AppendTask(ctx, idInt, *result)
	if err != nil {
		//Можно было сделать транзакцию на уровне бд, но я решил так быстрее
		e := s.rTask.Delete(ctx, result.ID)
		if e != nil {
			return nil, e
		}
		return nil, err
	}

	return result, nil
}

func (s *service) TaskFinish(ctx context.Context, taskId string) (*task.Task, error) {
	taskIdInt, err := s.checkIdParam(taskId)
	if err != nil {
		return nil, err
	}
	currTask, err := s.rTask.Get(ctx, taskIdInt)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	currTask.EndTime = &now
	duration := now.Sub(currTask.StartTime)
	currTask.TotalTime = &duration

	result, err := s.rTask.Put(ctx, taskIdInt, *currTask)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) TaskPut(ctx context.Context, id string, updateTask task.Task) (*task.Task, error) {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return nil, err
	}
	result, err := s.rTask.Put(ctx, idInt, updateTask)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) GetTask(ctx context.Context, startTimeStr string, endTimeStr string) ([]task.Task, error) {
	startTime, endTime, err := parseTimeStrings(startTimeStr, endTimeStr)
	if err != nil {
		return nil, err
	}

	result, err := s.rTask.GetLaborCost(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// Сортировка результатов
	sort.Sort(result)

	return result, nil
}

func parseTimeStrings(startTimeStr, endTimeStr string) (time.Time, time.Time, error) {
	layout := "2006-01-02T15:04:05.999999Z"

	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start time format: %w", err)
	}

	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end time format: %w", err)
	}

	return startTime, endTime, nil
}
func (s *service) GetAllTask(ctx context.Context) ([]task.Task, error) {
	result, err := s.rTask.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteTask(ctx context.Context, id string) error {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return err
	}

	err = s.rTask.Delete(ctx, idInt)
	if err != nil {
		return err
	}
	return nil
}
