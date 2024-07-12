package domain

import "time"

type Proxy struct {
	Id           int64
	Ip           string
	Port         int
	Protocol     string
	ResponseTime int64     `db:"response_time"`
	StatusId     int64     `db:"status_id"`
	CountryId    int64     `db:"country_id"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
func (p *Proxy) TableName() string { return "proxy" }

// available, not available, dont know
type Status struct {
	Id        int64
	Name      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
func (s *Status) TableName() string { return "status" }

type Country struct {
	Id        int64
	Name      string
	Code      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
func (c *Country) TableName() string { return "country" }

type Protocol struct {
	Id        int64
	Name      string
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
func (p *Protocol) TableName() string { return "protocol" }
