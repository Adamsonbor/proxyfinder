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

type User struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	PhotoUrl    string    `json:"photo_url" db:"photo_url"`
	DateOfBirth time.Time    `json:"date_of_birth" db:"date_of_birth"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) TableName() string { return "user" }

type Favorits struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id" db:"user_id"`
	ProxyId   int64     `json:"proxy_id" db:"proxy_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (s *Favorits) TableName() string { return "favorits" }

type Session struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiredAt time.Time `json:"expired_at" db:"expired_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (s *Session) TableName() string { return "session" }
