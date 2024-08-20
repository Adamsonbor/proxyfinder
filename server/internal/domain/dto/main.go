package dto

import "proxyfinder/internal/domain"

type Proxy struct {
	domain.Proxy
	Status  domain.Status  `json:"status"`
	Country domain.Country `json:"country"`
}
