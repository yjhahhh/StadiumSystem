package log

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
	"time"

	"common/config"

	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)


var one sync.Once

func Init() {
	one.Do(func() {
		logPath, exists := config.GetString("LogPath")
		if !exists {
			panic("logPath is not set")
		}
		logrus.SetFormatter(&GinFormatter{})
		logName := fmt.Sprintf("%s/gin-log-", logPath)
		r, err := rotatelogs.New(logName + "%Y%m%d.log")
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(io.MultiWriter(os.Stdout, r))
	})
}

type GinFormatter struct{}

func (m *GinFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	b.WriteString(fmt.Sprintf("[%s] [%s] %s\n", entry.Time.Format("2006-01-02 15:04:05"), entry.Level, entry.Message))
	return b.Bytes(), nil
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		customTime := fmt.Sprintf("%dms", int(math.Ceil(float64(time.Since(startTime).Nanoseconds()/1000000))))
		logrus.Infof("Gin log method:%+v , path:%+v , status:%+v , tc:%+v",
			c.Request.Method, c.Request.RequestURI, c.Writer.Status(), customTime)
	}
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args)
}

func Info(args ...interface{}) {
	logrus.Infoln(args)
}

func Error(args ...interface{}) {
	logrus.Errorln(args)
}

func Debug(args ...interface{}) {
	logrus.Debugln(args)
}

func Warn(format string, args ...interface{}) {
	logrus.Warnln(args)
}

func Fatal(args ...interface{}) {
	logrus.Fatalln(args)
}