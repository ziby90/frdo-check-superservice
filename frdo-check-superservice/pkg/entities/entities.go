package entities

import (
	"frdo-check-superservice/pkg/service"
)

type Entities struct {
	Service *service.Service
}

func NewEntities(service *service.Service) *Entities {
	return &Entities{Service: service}
}
