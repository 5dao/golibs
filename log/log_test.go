package log

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestAll(t *testing.T) {
	log.Debugf("debug: %s", "start")
}
