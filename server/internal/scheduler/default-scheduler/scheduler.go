package defaultscheduler

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"proxyfinder/internal/checker"
	"proxyfinder/internal/config"
	"proxyfinder/internal/domain"
	"proxyfinder/internal/storage"
	"sync"
	"time"
)


type Scheduler struct {
	cfg          *config.Config
	log          *slog.Logger
	checker      checker.Checker
	proxyStorage storage.ProxyStorage
	stopChan     chan os.Signal
}

func NewScheduler(cfg *config.Config, log *slog.Logger, proxyStorage storage.ProxyStorage, checker checker.Checker) *Scheduler {
	return &Scheduler{
		cfg:          cfg,
		log:          log,
		checker:      checker,
		proxyStorage: proxyStorage,
		stopChan:     make(chan os.Signal, 1),
	}
}

func (s *Scheduler) Run() {
	const op = "scheduler.Scheduler.Run"

	// log := s.log.With(slog.String("op", op))

	signal.Notify(s.stopChan, os.Interrupt)

	ticker := time.NewTicker(s.cfg.Scheduler.Interval)

	// s.Refresh(context.Background())

	select {
	case <-s.stopChan:
		s.Stop()
		return
	case <-ticker.C:
		s.Refresh(context.Background())
	}
}

func (s *Scheduler) Refresh(ctx context.Context) error {
	const op = "Scheduler.Refresh"

	log := s.log.With(slog.String("op", op))

	tx, err := s.proxyStorage.Begin()
	if err != nil {
		log.Error("BeginTx failed", slog.String("err", err.Error()))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Database.Timeout)
	defer cancel()

	proxies, err := s.proxyStorage.GetAll(ctx)
	if err != nil {
		log.Error("Proxy storage failed", slog.Any("err", err))
		return err
	}

	var wg sync.WaitGroup
	for _, v := range proxies {
		wg.Add(1)
		go func() {
			defer wg.Done()

			log := log.With(slog.String("proxy", fmt.Sprintf("%s:%d", v.Ip, v.Port)))

			log.Debug("Start checking...")

			checkCtx, cancel := context.WithTimeout(context.Background(), s.cfg.Checker.Timeout)
			start := time.Now()
			available, _ := s.checker.Check(checkCtx, &v)
			responseTime := time.Now().Sub(start).Milliseconds()
			log.Debug("End checking", slog.Int64("response time", responseTime))
			cancel()

			if available {
				v.StatusId = domain.STATUS_AVAILABLE
			} else {
				v.StatusId = domain.STATUS_UNAVAILABLE
			}

			log.Debug("End checking", slog.Int64("StatusId", v.StatusId))
			log.Debug("Start updating the proxy status...")

			dbCtx, cancel := context.WithTimeout(context.Background(), s.cfg.Checker.Timeout)
			err := s.proxyStorage.Update(dbCtx, tx, v.Id, &storage.ProxyUpdate{
				ResponseTime: &responseTime,
				StatusId:     &v.StatusId,
			})
			cancel()
			if err != nil {
				log.Error("Update failed", slog.String("err", err.Error()))
				return
			}
			log.Debug("Done!")
		}()
	}

	wg.Wait()

	log.Debug("Committing changes...")
	err = tx.Commit()
	if err != nil {
		log.Error("Commit failed", slog.String("err", err.Error()))
		tx.Rollback()
		return err
	}

	log.Debug("Refresh completed.")

	return nil
}

func (s *Scheduler) Stop() {
	s.log.Info("Gracefully stop")
	close(s.stopChan)
}
