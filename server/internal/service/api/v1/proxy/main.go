package proxyservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"proxyfinder/internal/domain/dto"
	apiv1 "proxyfinder/internal/service/api"
	serviceapiv1 "proxyfinder/internal/service/api"
	"proxyfinder/pkg/options"
)

type ProxyService struct {
	log     *slog.Logger
	storage apiv1.ProxyStorage
}

func New(log *slog.Logger, storage apiv1.ProxyStorage) *ProxyService {
	return &ProxyService{
		log:     log,
		storage: storage,
	}
}

func (self *ProxyService) GetAll(
	ctx context.Context,
	filter options.Options,
	sort options.Options,
) ([]dto.Proxy, error) {
	log := self.log.With(slog.String("op", "ProxyService.GetAll"))

	err := self.FieldsMap(filter)
	if err != nil {
		log.Warn("failed to map filter options", slog.String("err", err.Error()))
		return nil, err
	}

	log.Debug("options", slog.Any("filter", filter), slog.Any("sort", sort))

	err = self.FieldsMap(sort) 
	if err != nil {
		log.Warn("failed to map sort options", slog.String("err", err.Error()))
		return nil, err
	}
	fmt.Println(sort)

	proxies, err := self.storage.GetAll(ctx, filter, sort)
	if err != nil {
		log.Warn("failed to get proxies", slog.String("err", err.Error()))
		return nil, err
	}

	return proxies, nil
}

func (self *ProxyService) Update(ctx context.Context, filter options.Options) error {
	log := self.log.With(slog.String("op", "ProxyService.Update"))

	if !filter.Is() {
		return errors.New(apiv1.ErrIdNotFound)
	}

	err := self.IsValudUpdateOptions(filter)
	if err != nil {
		log.Warn("failed to validate options", slog.String("err", err.Error()))
		return err
	}

	return self.storage.Update(ctx, filter)
}

func (self *ProxyService) FieldsMap(opt options.Options) error {
	err := opt.MapField(func(field *options.Field) error {
		newFieldName, err := self.MapFieldName(field.Name)
		if err != nil {
			return err
		}

		field.Name = newFieldName
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// MapFieldName maps field name to column name in database
// Like: ip -> proxy.ip, country_name -> country.name
func (self *ProxyService) MapFieldName(fieldName string) (string, error) {
	switch fieldName {
	case "page", "perPage":
		return fieldName, nil
	case "id":
		return "proxy.id", nil
	case "ip":
		return "proxy.ip", nil
	case "port":
		return "proxy.port", nil
	case "protocol":
		return "proxy.protocol", nil
	case "response_time":
		return "proxy.response_time", nil
	case "country_id":
		return "country.id", nil
	case "country_name":
		return "country.name", nil
	case "country_code":
		return "country.code", nil
	case "status_id":
		return "status.id", nil
	case "status_name":
		return "status.name", nil
	default:
		self.log.Debug("invalid field", slog.String("field", fieldName))
		return "", fmt.Errorf(apiv1.ErrInvalidField)

	}
}

func (self *ProxyService) IsValudUpdateOptions(opt options.Options) error {
	err := errors.New(serviceapiv1.ErrIdNotFound)

	for _, v := range opt.Fields() {
		switch v.Name {
		case "ip":
		case "port":
		case "protocol":
		case "response_time":
		case "status_id":
		case "country_id":
		case "id":
			err = nil
		default:
			self.log.Debug("invalid field", slog.String("field", v.Name))
			return errors.New(serviceapiv1.ErrInvalidField)
		}
	}

	return err
}
