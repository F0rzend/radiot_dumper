package main

import (
	"log"
	"net/http"
	"time"

	"github.com/F0rzend/radiot_dumper/copier"
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configFileName = "dumper.yml"
)

type Config struct {
	SourceURL       string `yaml:"source_url" env:"SOURCE_URL"`
	FilePrefix      string `yaml:"file_prefix" env:"FILE_PREFIX"`
	FileDateFormat  string `yaml:"file_date_format" env:"FILE_DATE_FORMAT" env-default:"02_01_2006"`
	OutputDirectory string `yaml:"output_directory" env:"OUTPUT_DIRECTORY"`
	Timeout         string `yaml:"timeout" env:"TIMEOUT" env-default:"10s"`
}

func main() {
	cfg := Config{}
	if err := cleanenv.ReadConfig(configFileName, &cfg); err != nil {
		log.Fatal(err)
	}

	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		log.Fatal(err)
	}

	dumper := copier.NewDumberService(
		copier.NewStreamCopierService(
			&http.Client{
				Timeout: 0,
			},
		),
		copier.NewDatedFileBuilder(
			copier.DatedFileOptions{
				OutputDirectory: cfg.OutputDirectory,
				Prefix:          cfg.FilePrefix,
				DateFormat:      cfg.FileDateFormat,
			},
		),
		timeout,
	)

	if err := dumper.ListenAndCopy(cfg.SourceURL); err != nil {
		log.Fatal(err)
	}
}
