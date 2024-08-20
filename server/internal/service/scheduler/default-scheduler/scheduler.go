package defaultscheduler

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"proxyfinder/internal/config"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/domain/dto"
	serviceapiv1 "proxyfinder/internal/service/api"
	"proxyfinder/internal/service/checker"
	"proxyfinder/pkg/filter"
	"sync"
	"time"
)

type Scheduler struct {
	cfg          *config.Config
	log          *slog.Logger
	checker      checker.Checker
	proxyService serviceapiv1.ProxyService
	stopChan     chan os.Signal
}

func New(
	cfg *config.Config,
	log *slog.Logger,
	proxyService serviceapiv1.ProxyService,
	checker checker.Checker,
) *Scheduler {
	return &Scheduler{
		cfg:          cfg,
		log:          log,
		checker:      checker,
		proxyService: proxyService,
		stopChan:     make(chan os.Signal, 1),
	}
}

func (s *Scheduler) Run() {
	const op = "scheduler.Scheduler.Run"

	log := s.log.With(slog.String("op", op))
	log.Info("Start scheduler")

	signal.Notify(s.stopChan, os.Interrupt)

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Scheduler.Timeout)
	defer cancel()
	if s.cfg.Scheduler.StartImmediately {
		s.Refresh(ctx)
	}

	// Refreshing every hour
	ticker := time.NewTicker(s.cfg.Scheduler.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			s.Stop()
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Scheduler.Timeout)
			defer cancel()
			s.Refresh(ctx)
		}
	}
}

func (self *Scheduler) Refresh(ctx context.Context) error {
	log := self.log.With(slog.String("op", "Scheduler.Refresh"))
	log.Info("Start refreshing...")

	limit := 10

	options := filter.New()
	options.SetPage(1)
	options.SetPerPage(limit)
	options.UpdateLimitAndOffset()

	for range 1 {
		proxies, err := self.proxyService.GetAll(ctx, options)
		if len(proxies) == 0 {
			break
		}
		if err != nil {
			log.Error("Proxy storage failed", slog.Any("err", err))
			return err
		}

		var wg sync.WaitGroup
		for i := range proxies {
			wg.Add(1)
			go func(v dto.Proxy) {
				defer wg.Done()
				log := log.With(slog.String("proxy", fmt.Sprintf("%s:%d", v.Ip, v.Port)))

				err := self.RefreshOne(ctx, v)
				if err != nil {
					log.Error("Refresh failed", slog.Any("err", err))
					return
				}

				log.Debug("Done!")
			}(proxies[i])
		}
		wg.Wait()

		options.NextPage()
	}

	log.Info("Refresh completed.")

	return nil
}

func (self *Scheduler) RefreshOne(ctx context.Context, proxyDTO dto.Proxy) error {
	log := self.log.With(slog.String("op", "Scheduler.RefreshOne"))

	proxy := domain.Proxy{
		Id:           proxyDTO.Id,
		Ip:           proxyDTO.Ip,
		Port:         proxyDTO.Port,
		Protocol:     proxyDTO.Protocol,
		ResponseTime: proxyDTO.ResponseTime,
		StatusId:     proxyDTO.StatusId,
		CountryId:    proxyDTO.CountryId,
	}
	proxy, err := self.Check(ctx, proxy)
	if err != nil {
		log.Error("Check failed", slog.Any("err", err))
		return err
	}

	err = self.Update(ctx, proxy)
	if err != nil {
		if err.Error() == serviceapiv1.ErrIdNotFound {
			log.Error("Id not found", slog.Any("err", err))
			return err
		}

		log.Error("Update failed", slog.Any("err", err))
	}

	log.Info("Done!")
	return nil
}

func (self *Scheduler) Check(ctx context.Context, proxy domain.Proxy) (domain.Proxy, error) {
	log := self.log.With(slog.String("op", "Scheduler.Check"))
	log.Debug("Start checking...")

	checkCtx, checkCancel := context.WithTimeout(ctx, self.cfg.Checker.Timeout)
	defer checkCancel()

	start := time.Now()
	available, _ := self.checker.Check(checkCtx, proxy)
	proxy.ResponseTime = time.Now().Sub(start).Milliseconds()

	log.Debug("End checking", slog.Int64("response time", proxy.ResponseTime))

	if available {
		proxy.StatusId = domain.STATUS_AVAILABLE
	} else {
		proxy.StatusId = domain.STATUS_UNAVAILABLE
	}

	log.Debug("End checking", slog.Int64("StatusId", proxy.StatusId))

	return proxy, nil
}

func (self *Scheduler) Update(ctx context.Context, proxy domain.Proxy) error {
	dbCtx, dbCancel := context.WithTimeout(ctx, self.cfg.Database.Timeout)
	defer dbCancel()

	options := filter.New()
	options.AddField("id", filter.OpEq, proxy.Id, "int64")
	options.AddField("status_id", filter.OpEq, proxy.StatusId, "int64")
	options.AddField("response_time", filter.OpEq, proxy.ResponseTime, "int64")
	return self.proxyService.Update(dbCtx, options)
}

func (s *Scheduler) Stop() {
	s.log.Info("Gracefully stop")
	close(s.stopChan)
}
