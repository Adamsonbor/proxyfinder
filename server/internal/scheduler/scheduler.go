package scheduler

import (
	"proxyfinder/internal/domain"
	"time"
)

type Scheduler interface {
	Run(interval time.Duration)
	RunCycle(proxies []domain.Proxy)
}
