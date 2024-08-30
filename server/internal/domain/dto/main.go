package dto

import "proxyfinder/internal/domain"

type Proxy struct {
	domain.Proxy
	Status     domain.Status  `json:"status"`
	Country    domain.Country `json:"country"`
}

type FavoriteCreate struct {
	UserId  int64 `json:"user_id"`
	ProxyId int64 `json:"proxy_id"`
}
