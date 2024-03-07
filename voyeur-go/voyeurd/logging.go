package voyeurd

import (
	"github.com/Masterminds/log-go"
	lgrs "github.com/Masterminds/log-go/impl/logrus"
	"github.com/sirupsen/logrus"
)

func initLogging() log.Logger {
	var devFlag = true // evenutally pass in

	var lg = logrus.New()
	if devFlag {
		lg.Level = logrus.DebugLevel
		lg.SetFormatter(&logrus.TextFormatter{})
	} else {
		lg.Level = logrus.InfoLevel // logrus default
		lg.SetFormatter(&logrus.JSONFormatter{})
	}

	log.Current = lgrs.New(lg)

	return log.Current
}
