package log

import log "github.com/sirupsen/logrus"

type utcFormatter struct {
	log.Formatter
}

func (u utcFormatter) Format(e *log.Entry) ([]byte, error) {
	e.Time = e.Time.Local() // set time to local time
	return u.Formatter.Format(e)
}

func init() {
	log.SetFormatter(&utcFormatter{&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}})
}
