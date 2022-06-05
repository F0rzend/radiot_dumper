package copier

import (
	"io"
	"log"
)

func logClosing(closable io.Closer) {
	if err := closable.Close(); err != nil {
		log.Println(err)
	}
}
