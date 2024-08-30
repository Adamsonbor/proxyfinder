package userstorage

import (
	"context"
	"fmt"
	"proxyfinder/internal/domain"
	serviceapiv1 "proxyfinder/internal/service/api"

	"github.com/jmoiron/sqlx"
)

const (
	getByQuery        = "SELECT * FROM user WHERE %s = ?"
	getByRefreshToken = `
		SELECT 
			user.id, user.name, user.email, user.phone, user.photo_url,
			user.date_of_birth, user.created_at, user.updated_at
		FROM user
			JOIN session on user.id = session.user_id
		WHERE session.token = ?
	`
	getLatestSessionQuery = `SELECT * FROM session WHERE session.user_id = ? ORDER BY created_at DESC LIMIT 1`
	insertQuery        = "INSERT INTO user (name, email, phone, photo_url, date_of_birth) VALUES (?, ?, ?, ?, ?)"
	insertSessionQuery = "INSERT INTO session (user_id, token, expires_at) VALUES (?, ?, ?)"
)

type UserStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) serviceapiv1.UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (self *UserStorage) GetBy(ctx context.Context, fieldName string, value interface{}) (domain.User, error) {
	var user domain.User
	err := self.db.GetContext(ctx, &user, fmt.Sprintf(getByQuery, fieldName), value)
	return user, err
}

func (self *UserStorage) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	err := self.db.GetContext(ctx, &user, getByRefreshToken, refreshToken)
	return user, err
}

func (self *UserStorage) Save(ctx context.Context, user domain.User) (int64, error) {
	res, err := self.db.ExecContext(ctx, insertQuery, user.Name, user.Email, user.Phone, user.PhotoUrl, user.DateOfBirth)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (self *UserStorage) NewSession(ctx context.Context, userID int64, refreshToken string, expiresAt int64) error {
	_, err := self.db.ExecContext(ctx, insertSessionQuery, userID, refreshToken, expiresAt)
	return err
}
