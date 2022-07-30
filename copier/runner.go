package copier

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

const (
	readableTimeLayout = "02.01.2006 15:04:05"
)

type Runner struct {
	copier *StreamCopier
}

func NewRunner(copier *StreamCopier) *Runner {
	return &Runner{
		copier: copier,
	}
}

func (r *Runner) ScheduleRecording(
	cronString string,
	duration time.Duration,
	url string,
	outputFunc GetOutputFunc,
	delay time.Duration,
) error {
	c := cron.New()
	entryID, err := c.AddFunc(cronString, func() { r.record(url, outputFunc, duration, delay) })
	if err != nil {
		return err
	}

	c.Start()

	entry := getEntryByID(c, entryID)
	if entry == nil {
		return fmt.Errorf("entry with id %v not found", entryID)
	}

	r.copier.logger.Info().
		Str("url", url).
		Str("delay", delay.String()).
		Str("next_call", entry.Next.Format(readableTimeLayout)).
		Str("duration", duration.String()).
		Msg("Starting listening")

	return nil
}

func (r *Runner) record(
	url string,
	outputFunc GetOutputFunc,
	duration time.Duration,
	delay time.Duration,
) {
	start := time.Now()
	finish := start.Add(duration)
	for {
		r.copier.logger.Debug().Msg("RUN")
		if time.Now().After(finish) {
			return
		}
		if err := r.copier.CopyStream(url, outputFunc); err != nil && err != ErrStreamClosed {
			r.copier.logger.Error().Err(err).Msg("Error copying stream")
			return
		}
		time.Sleep(delay)
	}
}

func getEntryByID(c *cron.Cron, id cron.EntryID) *cron.Entry {
	for _, entry := range c.Entries() {
		if entry.ID == id {
			return &entry
		}
	}

	return nil
}
