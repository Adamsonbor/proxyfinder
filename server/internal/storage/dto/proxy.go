package dto

import (
	"proxyfinder/internal/domain"
)

type ProxyDTO struct {
	domain.Proxy
	Status       domain.Status  `json:"status"`
	Country      domain.Country `json:"country"`
}
