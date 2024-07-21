package service

import (
	"context"
	"effectiveMobile/pkg/db"
	"effectiveMobile/pkg/domain/people"
	"strconv"
)

func (s *service) Registration(ctx context.Context, newPeople people.Registration) (*int64, error) {
	err := newPeople.Validate()
	if err != nil {
		return nil, err
	}
	result, err := s.rPeople.Registration(ctx, newPeople)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) Login(ctx context.Context, people people.Registration) (int64, error) {
	id, err := s.rPeople.Login(ctx, people)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *service) PutPeople(ctx context.Context, id string, updatePeople people.Info) (*people.Info, error) {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return nil, err
	}

	// Валидация обновляемых данных
	result, err := s.rPeople.Put(ctx, idInt, updatePeople)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetPeople(ctx context.Context, filter *people.Filter, pagination *people.Pagination) ([]people.Request, error) {
	result, err := s.rPeople.Get(ctx, filter, pagination)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) InfoPeople(ctx context.Context, passportSerie string, passportNumber string) (*people.Info, error) {
	if len(passportSerie) != 4 || len(passportNumber) != 6 {
		return nil, db.ErrValidate
	} else if _, err := strconv.Atoi(passportSerie); err != nil {
		return nil, db.ErrPassportSerie
	} else if _, err := strconv.Atoi(passportNumber); err != nil {
		return nil, db.ErrPassportNumber
	}
	passport := passportSerie + " " + passportNumber
	result, err := s.rPeople.Info(ctx, passport)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeletePeople(ctx context.Context, id string) error {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return err
	}

	err = s.rPeople.Delete(ctx, idInt)
	if err != nil {
		return err
	}
	return nil
}
