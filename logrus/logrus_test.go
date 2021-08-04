package logrus

import (
	"errors"
	"github.com/longzhoufeng/go-logger"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestName(t *testing.T) {
	l := NewLogger()

	if l.String() != "logrus" {
		t.Errorf("error: name expected 'logrus' actual: %s", l.String())
	}

	t.Logf("testing logger name: %s", l.String())
}

func TestWithFields(t *testing.T) {
	l := NewLogger(go_logger.WithOutput(os.Stdout)).Fields(map[string]interface{}{
		"k1": "v1",
		"k2": 123456,
	})

	go_logger.DefaultLogger = l

	go_logger.Log(go_logger.InfoLevel, "testing: Info")
	go_logger.Logf(go_logger.InfoLevel, "testing: %s", "Infof")
}

func TestWithError(t *testing.T) {
	l := NewLogger().Fields(map[string]interface{}{"error": errors.New("boom!")})
	go_logger.DefaultLogger = l

	go_logger.Log(go_logger.InfoLevel, "testing: error")
}

func TestWithLogger(t *testing.T) {
	// with *logrus.Logger
	l := NewLogger(WithLogger(logrus.StandardLogger())).Fields(map[string]interface{}{
		"k1": "v1",
		"k2": 123456,
	})
	go_logger.DefaultLogger = l
	go_logger.Log(go_logger.InfoLevel, "testing: with *logrus.Logger")

	// with *logrus.Entry
	el := NewLogger(WithLogger(logrus.NewEntry(logrus.StandardLogger()))).Fields(map[string]interface{}{
		"k3": 3.456,
		"k4": true,
	})
	go_logger.DefaultLogger = el
	go_logger.Log(go_logger.InfoLevel, "testing: with *logrus.Entry")
}

func TestJSON(t *testing.T) {
	go_logger.DefaultLogger = NewLogger(WithJSONFormatter(&logrus.JSONFormatter{}))

	go_logger.Logf(go_logger.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	go_logger.DefaultLogger = NewLogger()

	go_logger.Init(go_logger.WithLevel(go_logger.DebugLevel))
	go_logger.Logf(go_logger.DebugLevel, "test show debug: %s", "debug msg")

	go_logger.Init(go_logger.WithLevel(go_logger.InfoLevel))
	go_logger.Logf(go_logger.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	go_logger.DefaultLogger = NewLogger(ReportCaller())

	go_logger.Logf(go_logger.InfoLevel, "testing: %s", "WithReportCaller")
}
