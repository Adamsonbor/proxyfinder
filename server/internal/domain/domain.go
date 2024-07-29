package domain

import "time"

const (
	STATUS_AVAILABLE   = int64(1)
	STATUS_UNAVAILABLE = int64(2)
)

type Proxy struct {
	Id           int64     `json:"id"`
	Ip           string    `json:"ip"`
	Port         int       `json:"port"`
	Protocol     string    `json:"protocol"`
	ResponseTime int64     `json:"response_time" db:"response_time"`
	StatusId     int64     `json:"status_id" db:"status_id"`
	CountryId    int64     `json:"country_id" db:"country_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (p *Proxy) TableName() string { return "proxy" }

// available, not available, dont know
type Status struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (s *Status) TableName() string { return "status" }

type Country struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (c *Country) TableName() string { return "country" }

//
// type Protocol struct {
// 	Id        int64
// 	Name      string
// 	CreatedAt time.Time `db:"created_at"`
// 	UpdatedAt time.Time `db:"updated_at"`
// }
// func (p *Protocol) TableName() string { return "protocol" }
