package copier

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
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
	ctx context.Context,
	cronString string,
	duration time.Duration,
	url string,
	outputFunc GetOutputFunc,
	delay time.Duration,
) error {
	c := cron.New()
	entryID, err := c.AddFunc(
		cronString,
		func() { r.record(ctx, url, outputFunc, duration, delay) },
	)
	if err != nil {
		return err
	}

	c.Start()

	entry := getEntryByID(c, entryID)
	if entry == nil {
		return fmt.Errorf("entry with id %v not found", entryID)
	}

	zerolog.Ctx(ctx).Info().
		Str("url", url).
		Str("delay", delay.String()).
		Str("next_call", entry.Next.Format(readableTimeLayout)).
		Str("duration", duration.String()).
		Msg("Starting listening")

	return nil
}

func (r *Runner) record(
	ctx context.Context,
	url string,
	outputFunc GetOutputFunc,
	duration time.Duration,
	delay time.Duration,
) {
	start := time.Now()
	finish := start.Add(duration)
	zerolog.Ctx(ctx).Info().Msg("Starting recording")
	for {
		if time.Now().After(finish) {
			zerolog.Ctx(ctx).Info().Msg("recording finished")
			return
		}

		if err := r.copier.CopyStream(ctx, url, outputFunc); err != nil && err != ErrStreamClosed {
			zerolog.Ctx(ctx).Error().Err(err).Msg("Error copying stream")
		}
		time.Sleep(delay)
		zerolog.Ctx(ctx).Debug().Msg("Retry recording")
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
