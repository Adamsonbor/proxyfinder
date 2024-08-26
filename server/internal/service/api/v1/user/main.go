package userservice

import (
	"context"
	"errors"
	"log/slog"
	"proxyfinder/internal/domain"
	serviceapiv1 "proxyfinder/internal/service/api"
)

const (
	ErrInvalidUser = "invalid user"
)

type UserService struct {
	log     *slog.Logger
	storage serviceapiv1.UserStorage
}

func New(log *slog.Logger, storage serviceapiv1.UserStorage) serviceapiv1.UserService {
	return &UserService{
		log:     log,
		storage: storage,
	}
}

func (self *UserService) UserInfo(ctx context.Context, id int64) (domain.User, error) {
	log := self.log.With(slog.String("op", "UserService.UserInfo"))

	user, err := self.storage.GetBy(ctx, "id", id)
	if err != nil {
		log.Debug("failed to get user info", slog.Int64("user_id", id), slog.Any("error", err))
		return domain.User{}, err
	}

	return user, nil
}

func (self *UserService) GetBy(ctx context.Context, fieldName string, value interface{}) (domain.User, error) {
	log := self.log.With(slog.String("op", "UserService.GetBy"))

	switch fieldName {
		case "refresh_token":
			return self.storage.GetByRefreshToken(ctx, self.MapFieldName(fieldName))
		case "id", "name", "email", "phone":
			return self.storage.GetBy(ctx, fieldName, value)
		default:
			log.Debug("invalid field name", slog.String("field_name", fieldName))
			return domain.User{}, errors.New(serviceapiv1.ErrInvalidField)
	}
}

func (self *UserService) Save(ctx context.Context, user domain.User) (int64, error) {
	log := self.log.With(slog.String("op", "UserService.Save"))

	err := self.IsValidUser(user)
	if err != nil {
		log.Debug("invalid user", slog.Any("user", user))
		return 0, err
	}

	return self.storage.Save(ctx, user)
}

func (self *UserService) NewSession(ctx context.Context, userId int64, token string) error {
	log := self.log.With(slog.String("op", "UserService.NewSession"))
	err := self.storage.NewSession(ctx, userId, token)
	if err != nil {
		log.Debug("failed to create new session", slog.Int64("user_id", userId))
		return err
	}

	return nil
}

func (self *UserService) MapFieldName(fieldName string) string {
	switch fieldName {
	case "refresh_token":
		return "session.refresh_token"
	default:
		return fieldName
	}
}

func (self *UserService) IsValidFieldName(fieldName string) error {
	switch fieldName {
	case "id", "name", "email", "phone", "refresh_token":
		return nil
	default:
		return errors.New(serviceapiv1.ErrInvalidField)
	}
}

func (self *UserService) IsValidUser(user domain.User) error {
	if user.Email == "" && user.Id == 0 {
		return errors.New(ErrInvalidUser)
	}

	return nil
}
