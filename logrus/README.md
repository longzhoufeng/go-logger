# logrus

[logrus](https://github.com/sirupsen/logrus) logger implementation for __go-core__ [meta logger](github.com/longzhoufeng/go-logger/tree/master/logger).

## Usage

```go
import (
	"os"
	"github.com/sirupsen/logrus"
	"github.com/longzhoufeng/go-logger/logger"
)

func ExampleWithOutput() {
	logger.DefaultLogger = NewLogger(logger.WithOutput(os.Stdout))
	logger.Infof("testing: %s", "Infof")
}

func ExampleWithLogger() {
	l := logrus.New() // *logrus.Logger
	logger.DefaultLogger = NewLogger(WithLogger(l))
	logger.Infof("testing: %s", "Infof")
}
```

