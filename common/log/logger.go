package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

var E = logrus.New()
var I = logrus.New()

type formatter struct {
	Prefix string
}

func (f *formatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Message = "[" + f.Prefix + "] " + entry.Message
	te := logrus.TextFormatter{}
	return te.Format(entry)
}

func InitLogger() {
	initSingleLogger(E)
	initSingleLogger(I)
}
func initSingleLogger(log *logrus.Logger) {
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{

		FullTimestamp: true,
	})
	log.SetFormatter(&formatter{Prefix: "External"})
}
