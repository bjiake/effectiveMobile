package handler

import (
	services "effectiveMobile/pkg/service/interface"
)

type Handler struct {
	service services.ServiceUseCase
}

func NewHandler(service services.ServiceUseCase) *Handler {
	return &Handler{
		service: service,
	}
}
