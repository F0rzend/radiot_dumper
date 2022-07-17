package main

import (
	"context"
	"fmt"
	syslog "log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/F0rzend/radiot_dumper/copier"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	configFileName = "dumper.yml"
)

type Config struct {
	SourceURL       string `yaml:"source_url"       env:"SOURCE_URL"`
	FilePrefix      string `yaml:"file_prefix"      env:"FILE_PREFIX"`
	Schedule        string `yaml:"schedule"         env:"SCHEDULE"`
	Duration        string `yaml:"duration"         env:"DURATION"`
	OutputDirectory string `yaml:"output_directory" env:"OUTPUT_DIRECTORY"`
	FileDateFormat  string `yaml:"file_date_format" env:"FILE_DATE_FORMAT" env-default:"02_01_2006"`
	Delay           string `yaml:"delay"            env:"DELAY"            env-default:"5s"`
	LogLevel        string `yaml:"log_level"        env:"LOG_LEVEL"        env-default:"info"`
}

func Run() error {
	cfg := Config{}

	_, err := os.Stat(configFileName)
	if os.IsNotExist(err) {
		fmt.Printf("Config file %s not found, using environment variables\n", configFileName)
		err = cleanenv.ReadEnv(&cfg)
	} else {
		fmt.Printf("Using config file %s\n", configFileName)
		err = cleanenv.ReadConfig(configFileName, &cfg)
	}

	if err != nil {
		return err
	}

	delay, err := time.ParseDuration(cfg.Delay)
	if err != nil {
		return err
	}
	duration, err := time.ParseDuration(cfg.Duration)
	if err != nil {
		return err
	}

	datedFileBuilder := copier.NewDatedFileBuilder(
		cfg.OutputDirectory,
		os.DirFS(cfg.OutputDirectory),
		cfg.FilePrefix,
		cfg.FileDateFormat,
	)

	logger := log.
		Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		Level(Must(zerolog.ParseLevel(cfg.LogLevel))).
		With().
		Caller().
		Logger()

	streamCopier := copier.NewStreamCopier(
		&http.Client{
			Timeout: 0,
		},
		logger,
	)

	runner := copier.NewRunner(streamCopier)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)

	if err = runner.ScheduleRecording(
		cfg.Schedule,
		duration,
		cfg.SourceURL,
		datedFileBuilder.GetOutput,
		delay,
	); err != nil {
		logger.Error().Err(err).Msg("Error scheduling recording")
	}

	<-ctx.Done()
	return nil
}

func main() {
	if err := Run(); err != nil {
		syslog.Fatal(err)
	}
}

func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
