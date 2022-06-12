package main

import (
	syslog "log"
	"net/http"
	"os"
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
	FileDateFormat  string `yaml:"file_date_format" env:"FILE_DATE_FORMAT" env-default:"02_01_2006"`
	OutputDirectory string `yaml:"output_directory" env:"OUTPUT_DIRECTORY"`
	Delay           string `yaml:"delay"            env:"DELAY"            env-default:"10s"`
	LogLevel        string `yaml:"log_level"        env:"LOG_LEVEL"        env-default:"info"`
}

func Run() error {
	cfg := Config{}

	_, err := os.Stat(configFileName)
	if os.IsNotExist(err) {
		err = cleanenv.ReadEnv(&cfg)
	} else {
		err = cleanenv.ReadConfig(configFileName, &cfg)
	}

	if err != nil {
		return err
	}

	delay, err := time.ParseDuration(cfg.Delay)
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

	logger.Info().Str("url", cfg.SourceURL).Dur("delay", delay).Msg("Starting dumping")
	return streamCopier.ListenAndCopy(
		cfg.SourceURL,
		datedFileBuilder.GetOutput,
		delay,
	)
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
