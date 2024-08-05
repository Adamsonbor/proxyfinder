package smtpmail

import (
	"context"
	"fmt"
	"proxyfinder/internal/config"
	"proxyfinder/internal/domain"
)

type MailNotifier struct {
	cfg     *config.Config
	service *MailService
	proxies []int64
	email   string
}

func (self *MailNotifier) AddProxy(proxy *domain.Proxy) {
	self.proxies = append(self.proxies, proxy.Id)
}

func (self *MailNotifier) Update(ctx context.Context, proxies []domain.Proxy) {
	msg := `Proxy is available: %s:%d`

	for _, v := range proxies {
		if v.StatusId == domain.STATUS_AVAILABLE {
			self.service.SendMail(
				ctx,
				self.email,
				fmt.Sprintf("Proxy %s:%d is available", v.Ip, v.Port),
				[]byte(fmt.Sprintf(msg, v.Ip, v.Port)),
			)
		}
	}
}

func NewMailNotifier(cfg *config.Config, email string, service *MailService) *MailNotifier {
	return &MailNotifier{
		cfg:     cfg,
		email:   email,
		service: service,
	}
}
