package service

import (
	"effectiveMobile/pkg/db"
	peopleI "effectiveMobile/pkg/repo/people/interface"
	taskI "effectiveMobile/pkg/repo/task/interface"

	"context"
	interfaces "effectiveMobile/pkg/service/interface"
	"strconv"
)

type service struct {
	rPeople peopleI.PeopleRepository
	rTask   taskI.TaskRepository
}

func NewService(
	peopleRepository peopleI.PeopleRepository,
	taskRepository taskI.TaskRepository,
) interfaces.ServiceUseCase {
	return &service{
		rPeople: peopleRepository,
		rTask:   taskRepository,
	}
}

func (s *service) Migrate(ctx context.Context) error {
	if err := s.rPeople.Migrate(ctx); err != nil {
		return err
	}
	if err := s.rTask.Migrate(ctx); err != nil {
		return err
	}

	return nil
}

func (s *service) checkIdParam(id string) (int64, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil || idInt <= 0 {
		return 0, db.ErrParamNotFound
	}
	return idInt, nil
}
