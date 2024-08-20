package proxyservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"proxyfinder/internal/domain/dto"
	apiv1 "proxyfinder/internal/service/api"
	serviceapiv1 "proxyfinder/internal/service/api"
	"proxyfinder/pkg/filter"
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

func (self *ProxyService) GetAll(ctx context.Context, options filter.Options) ([]dto.Proxy, error) {
	log := self.log.With(slog.String("op", "ProxyService.GetAll"))
	options, err := self.OptionsMap(options)
	if err != nil {
		log.Warn("failed to map options", slog.String("err", err.Error()))
		return nil, err
	}
	log.Debug("options", slog.Any("options", options))

	proxies, err :=  self.storage.GetAll(ctx, options)
	if err != nil {
		log.Warn("failed to get proxies", slog.String("err", err.Error()))
		return nil, err
	}

	return proxies, nil
}

func (self *ProxyService) Update(ctx context.Context, options filter.Options) error {
	log := self.log.With(slog.String("op", "ProxyService.Update"))

	if !options.Is() {
		return errors.New(apiv1.ErrIdNotFound)
	}

	err := self.IsValudUpdateOptions(options)
	if err != nil {
		log.Warn("failed to validate options", slog.String("err", err.Error()))
		return err
	}

	return self.storage.Update(ctx, options)
}

func (self *ProxyService) OptionsMap(options filter.Options) (filter.Options, error) {
	newOpts := filter.New()
	for _, field := range options.Fields() {
		newFieldName, err := self.MapFieldName(field.Name)
		if err != nil {
			return nil, err
		}

		err = newOpts.AddField(newFieldName, field.Op, field.Val, field.Type)
		if err != nil {
			return nil, err
		}
	}
	newOpts.SetLimit(options.Limit())
	newOpts.SetOffset(options.Offset())

	return newOpts, nil
}

// MapFieldName maps field name to column name in database
// Like: ip -> proxy.ip, country_name -> country.name
func (self *ProxyService) MapFieldName(fieldName string) (string, error) {
	switch fieldName {
	case "ip":
		return "proxy.ip", nil
	case "port":
		return "proxy.port", nil
	case "protocol":
		return "proxy.protocol", nil
	case "response_time":
		return "proxy.response_time", nil
	case "country_name":
		return "country.name", nil
	case "country_code":
		return "country.code", nil
	case "status_name":
		return "status.name", nil
	default:
		return "", fmt.Errorf(apiv1.ErrInvalidField)

	}
}

func (self *ProxyService) IsValudUpdateOptions(options filter.Options) error {
	err := errors.New(serviceapiv1.ErrIdNotFound)

	for _, v := range options.Fields() {
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
			return errors.New(serviceapiv1.ErrInvalidField)
		}
	}

	return err
}
