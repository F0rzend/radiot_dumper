package main

import (
	"github.com/F0rzend/readiot-dumper/internal"
	"log"
	"net/http"
	"time"
)

const (
	path    = "https://stream.radio-t.com/"
	timeout = 10 * time.Second
)

func main() {
	dumper := internal.NewDumberService(
		internal.NewStreamCopierService(
			&http.Client{
				Timeout: 10 * time.Second,
			},
		),
		internal.NewDatedFileBuilder(
			internal.DatedFileOptions{
				Prefix:     "radio-t_",
				DateFormat: "02_01_2006",
				Extension:  "mp3",
			},
		),
		timeout,
	)

	log.Fatal(dumper.ListenAndCopy(path))
}
