package countryservice

import (
	"context"
	"errors"
	"log/slog"
	"proxyfinder/internal/domain"
	serviceapiv1 "proxyfinder/internal/service/api"
	"proxyfinder/pkg/options"
)

type CountryService struct {
	log     *slog.Logger
	storage serviceapiv1.CountryStorage
}

func New(log *slog.Logger, storage serviceapiv1.CountryStorage) *CountryService {
	return &CountryService{
		log:     log,
		storage: storage,
	}
}

func (self *CountryService) GetAll(ctx context.Context, filter options.Options, sort options.Options) ([]domain.Country, error) {
	log := self.log.With(slog.String("op", "CountryService.GetAll"))

	err := self.ValidateFilters(filter)
	if err != nil {
		log.Warn("failed to validate filters", slog.String("err", err.Error()))
		return nil, err
	}
	return self.storage.GetAll(ctx, filter, sort)
}

func (self *CountryService) ValidateFilters(options options.Options) error {
	log := self.log.With(slog.String("op", "CountryService.ValidateFilters"))

	for _, v := range options.Fields() {
		switch v.Name {
		case "page", "perPage", "name", "code":
		default:
			log.Warn("invalid field", slog.String("field", v.Name))
			return errors.New(serviceapiv1.ErrInvalidField)
		}
	}
	return nil
}
