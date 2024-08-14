package sqlxstorage

import (
	"context"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/storage"

	"github.com/jmoiron/sqlx"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (self *UserStorage) Begin(ctx context.Context) (*sqlx.Tx, error) {
	return self.db.BeginTxx(ctx, nil)
}

func (self *UserStorage) NewSession(ctx context.Context, tx *sqlx.Tx, user_id int64, refreshToken string) error {
	query := "INSERT INTO session (user_id, token) VALUES (?, ?)"

	_, err := self.db.ExecContext(ctx, query, user_id, refreshToken)
	return err
}

func (self *UserStorage) UpdateSession(ctx context.Context, tx *sqlx.Tx, user_id int64, refreshToken string) error {
	query := "UPDATE session SET token = ? WHERE user_id = ?"

	_, err := self.db.ExecContext(ctx, query, refreshToken, user_id)
	return err
}

func (self *UserStorage) Create(ctx context.Context, tx *sqlx.Tx, user *domain.User) (*domain.User, error) {
	query := "INSERT INTO user (email, name, photo_url, phone, date_of_birth) VALUES (?, ?, ?, ?, ?)"

	res, err := self.db.ExecContext(ctx, query, user.Email, user.Name, user.PhotoUrl, user.Phone, user.DateOfBirth)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	user.Id = id

	return user, nil
}

func (self *UserStorage) GetBy(ctx context.Context, field string, value interface{}) (*domain.User, error) {
	query := "SELECT * FROM user as u WHERE u." + field + " = ?"

	var user domain.User
	err := self.db.GetContext(ctx, &user, query, value)
	return &user, storage.ErrRecordNotFoundWrap(err)
}

func (self *UserStorage) GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error) {
	query := "SELECT u.* FROM user as u JOIN session as s ON u.id = s.user_id WHERE s.token = ?"

	var user domain.User
	err := self.db.GetContext(ctx, &user, query, refreshToken)

	return &user, storage.ErrRecordNotFoundWrap(err)
}
