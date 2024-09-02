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

func (self *FavoritsService) Save(ctx context.Context, options options.Options) (int64, error) {
	log := self.log.With(slog.String("op", "FavoritsService.Create"))

	err := self.ValidateOptions(options)
	if err != nil {
		log.Warn("failed to validate options", slog.String("err", err.Error()))
	}

	err = self.ValidateOrderOptions(options)
	if err != nil {
		log.Warn("failed to validate order options", slog.String("err", err.Error()))
		return 0, err
	}

	return self.storage.Save(ctx, options)
}

func (self *FavoritsService) Delete(ctx context.Context, options options.Options) error {
	log := self.log.With(slog.String("op", "FavoritsService.Delete"))

	err := self.ValidateOptions(options)
	if err != nil {
		log.Warn("failed to validate options", slog.String("err", err.Error()))
	}

	return self.storage.Delete(ctx, options)
}

func (self *FavoritsService) ValidateOrderOptions(options options.Options) error {
	order := map[string]int{
		"user_id": 0,
		"proxy_id": 1,
	}

	fields := options.Fields()
	for indx := 0; indx < len(fields)-1; indx++ {
		value := fields[indx].Name
		nextValue := fields[indx+1].Name

		if order[value] > order[nextValue] {
			return errors.New(apiv1.ErrInvalidFieldsOrder)
		}
	}

	return nil
}

func (self *FavoritsService) ValidateOptions(options options.Options) error {
	for _, opt := range options.Fields() {
		switch opt.Name {
		case "page", "perPage":
		case "user_id":
		case "proxy_id":
		default:
			return errors.New(apiv1.ErrInvalidField)
		}
	}
	return nil
}
