package favoritsservice

import (
	"context"
	"errors"
	"log/slog"
	"proxyfinder/internal/domain"
	apiv1 "proxyfinder/internal/service/api"
	"proxyfinder/pkg/options"
)

type FavoritsService struct {
	log     *slog.Logger
	storage apiv1.FavoritsStorage
}

func New(log *slog.Logger, storage apiv1.FavoritsStorage) *FavoritsService {
	return &FavoritsService{
		log:     log,
		storage: storage,
	}
}

func (self *FavoritsService) GetAll(ctx context.Context, options options.Options, sort options.Options) ([]domain.Favorits, error) {
	log := self.log.With(slog.String("op", "FavoritsService.GetAll"))

	err := self.ValidateOptions(options)
	if err != nil {
		log.Warn("failed to validate options", slog.String("err", err.Error()))
	}

	return self.storage.GetAll(ctx, options, sort)
}

func (self *FavoritsService) ValidateOptions(options options.Options) error {
	for _, opt := range options.Fields() {
		switch opt.Name {
		case "user_id":
		case "proxy_id":
		default:
			return errors.New(apiv1.ErrInvalidField)
		}
	}
	return nil
}
