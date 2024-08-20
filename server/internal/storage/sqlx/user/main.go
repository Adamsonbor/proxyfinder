package userstorage

import (
	"context"
	"fmt"
	"proxyfinder/internal/domain"

	"github.com/jmoiron/sqlx"
)

const (
	getByQuery        = "SELECT * FROM users WHERE %s = ?"
	getByRefreshToken = `
		SELECT user.*
		FROM user
			JOIN session on session.user_id = user.id
		WHERE session.refresh_token = ?
		`
	insertQuery = "INSERT INTO users (name, email) VALUES (?, ?)"
	insertSessionQuery = "INSERT INTO session (user_id, refresh_token) VALUES (?, ?)"
)

type UserStorage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *UserStorage {
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
	res, err := self.db.ExecContext(ctx, insertQuery, user.Name, user.Email)
	if err != nil {
		return 0, err
	}

	return  res.LastInsertId()
}

func (self *UserStorage) NewSession(ctx context.Context, userID int64, refreshToken string) error {
	_, err := self.db.ExecContext(ctx, insertSessionQuery, userID, refreshToken)
	return err
}
