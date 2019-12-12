package factory

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewDispatcher(t *testing.T) {
	var log = logrus.New()
	if dispatcher := CreateNewDispatcher(log); dispatcher == nil {
		t.Fatal("Failed to create new dispatcher.")
	}
}
