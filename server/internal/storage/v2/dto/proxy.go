package dto

import "proxyfinder/internal/domain"

type ProxyDTO struct {
	Id           int64          `json:"id"`
	Ip           string         `json:"ip"`
	Port         int            `json:"port"`
	Protocol     string         `json:"protocol"`
	ResponseTime int64          `json:"response_time" db:"response_time"`
	Status       domain.Status  `json:"status"`
	Country      domain.Country `json:"country"`
}
