package logs

import (
	"os"

	"github.com/rs/zerolog"
)

func workLog() (log *zerolog.Logger) {
	file, err := os.OpenFile(
		"log.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	logger := zerolog.New(file).With().Timestamp().Logger()

}
