package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)
}
