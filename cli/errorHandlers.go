package cli

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func ExitWithError() {
	os.Exit(1)
}

func ExitOnError(error error, msg string) {
	if error != nil {
		log.Errorf("%s: %s", msg, error)
		ExitWithError()
	}
}

func ExitWithErrorMsg(errorMsg string) {
	log.Errorf(errorMsg)
	ExitWithError()
}
