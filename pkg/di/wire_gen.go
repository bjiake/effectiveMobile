package di

import (
	"context"
	http "effectiveMobile/pkg/api"
	"effectiveMobile/pkg/api/handler"
	"effectiveMobile/pkg/config"
	"effectiveMobile/pkg/db"
	"effectiveMobile/pkg/repo/people"
	"effectiveMobile/pkg/repo/task"
	"effectiveMobile/pkg/service"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	bd, err := db.ConnectToBD(cfg)
	if err != nil {
		return nil, err
	}
	// Repository
	peopleRepository := people.NewPeopleDataBase(bd)
	taskRepository := task.NewTaskDataBase(bd)

	//service - logic
	userService := service.NewService(peopleRepository, taskRepository)

	// Init Migrate
	err = userService.Migrate(context.Background())
	if err != nil {
		return nil, err
	}

	userHandler := handler.NewHandler(userService)
	serverHTTP := http.NewServerHTTP(userHandler)

	return serverHTTP, nil
}
