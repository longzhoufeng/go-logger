package zap

import (
	"fmt"
	"github.com/longzhoufeng/go-core/debug/writer"
	"github.com/longzhoufeng/go-logger"
	"testing"
)

func TestName(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	if l.String() != "zap" {
		t.Errorf("name is error %s", l.String())
	}

	t.Logf("test logger name: %s", l.String())
}

func TestLogf(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}

	go_logger.DefaultLogger = l
	go_logger.Logf(go_logger.InfoLevel, "test logf: %s", "name")
}

func TestSetLevel(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	go_logger.DefaultLogger = l

	go_logger.Init(go_logger.WithLevel(go_logger.DebugLevel))
	l.Logf(go_logger.DebugLevel, "test show debug: %s", "debug msg")

	go_logger.Init(go_logger.WithLevel(go_logger.InfoLevel))
	l.Logf(go_logger.DebugLevel, "test non-show debug: %s", "debug msg")
}

func TestWithReportCaller(t *testing.T) {
	var err error
	go_logger.DefaultLogger, err = NewLogger(WithCallerSkip(0))
	if err != nil {
		t.Fatal(err)
	}

	go_logger.Logf(go_logger.InfoLevel, "testing: %s", "WithReportCaller")
}

func TestFields(t *testing.T) {
	l, err := NewLogger()
	if err != nil {
		t.Fatal(err)
	}
	go_logger.DefaultLogger = l.Fields(map[string]interface{}{
		"x-request-id": "123456abc",
	})
	go_logger.DefaultLogger.Log(go_logger.InfoLevel, "hello")
}

func TestFile(t *testing.T) {
	outFileWriter, err := writer.NewFileWriter(writer.WithPath("testdata"),writer.WithSuffix("log"))
	if err != nil {
		t.Errorf("logger setup error: %s", err.Error())
	}
	//var err error
	go_logger.DefaultLogger, err = NewLogger(go_logger.WithLevel(go_logger.TraceLevel), WithOutput(outFileWriter))
	if err != nil {
		t.Errorf("logger setup error: %s", err.Error())
	}
	go_logger.DefaultLogger = go_logger.DefaultLogger.Fields(map[string]interface{}{
		"x-request-id": "123456abc",
	})
	fmt.Println(go_logger.DefaultLogger)
	go_logger.DefaultLogger.Log(go_logger.InfoLevel, "hello")
}
