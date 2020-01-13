package process

import (
	"sync/atomic"
	"time"

	"github.com/tradeo/golib/log"
	"github.com/tradeo/golib/util"
)

// CronJob defines a generic requirements for creating cron jobs.
type CronJob interface {
	Exec

	// Called before the cron cycle begins.
	Begin(period time.Duration) error

	// Called when the cron service is shutdown.
	End()
}

// Exec is a simple interface with just Exec method.
type Exec interface {
	Exec() error
}

// Cron is a process that calls periodically a given job.
// Implements the process.Process interface.
type Cron struct {
	period        time.Duration
	job           CronJob
	stopped       uint32
	done          chan struct{}
	tickerFactory util.TickerFactory
}

// NewCron creates cron instances.
func NewCron(job CronJob, period time.Duration) *Cron {
	return &Cron{
		job:           job,
		period:        period,
		stopped:       0,
		done:          make(chan struct{}),
		tickerFactory: util.NewTickerFactory(),
	}
}

// Stop stops the cron.
func (c *Cron) Stop() error {
	if atomic.CompareAndSwapUint32(&c.stopped, 0, 1) {
		close(c.done)
	}

	return nil
}

// Run starts the cron processing.
func (c *Cron) Run() error {
	ticker := c.tickerFactory.Create(c.period)
	defer ticker.Stop()
	defer c.job.End()

	err := c.job.Begin(c.period)
	if err != nil {
		return err
	}

	err = c.job.Exec()
	if err != nil {
		return err
	}

	for {
		select {
		case <-c.done:
			return nil
		case <-ticker.Chan():
			err := c.job.Exec()
			if err != nil {
				return err
			}
		}
	}
}

// CronJobLog a common cron job decorator adding some logging.
type CronJobLog struct {
	job    Exec
	logger log.Logger
	name   string
}

// NewCronJob creates new cron job.
// Implements the CronJob interface.
func NewCronJob(name string, job Exec, logger log.Logger) *CronJobLog {
	return &CronJobLog{
		job:    job,
		logger: logger,
		name:   name,
	}
}

// Begin implements the CronJob interface.
func (j *CronJobLog) Begin(period time.Duration) error {
	j.logger.Infof("cron job %s start period=%v", j.name, period)
	return nil
}

// Exec implements the CronJob interface.
func (j *CronJobLog) Exec() error {
	err := j.job.Exec()
	if err != nil {
		j.logger.Errorf("cron job %s failed: %v", j.name, err)
	}
	return err
}

// End implements the CronJob interface.
func (j *CronJobLog) End() {
	j.logger.Infof("cron job %s shutdown", j.name)
}
