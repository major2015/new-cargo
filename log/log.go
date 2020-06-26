package log

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	log          = logrus.New()
	rawLog       = logrus.New()
	translations = make(map[string]string)
)

// LoadTranslations takes a map[string]interface and flattens it to map[string]string
// Because translations have been loaded - we internally override log the formatter
// Nested entries are accessible using dot notation.
// example:   `{"foo": {"bar": "baz"}}`
// flattened: `foo.bar: baz`
func LoadTranslations(things map[string]interface{}) {
	formatter := new(prefixed.TextFormatter)
	formatter.TimestampFormat = `Jan 02 15:04:05`
	formatter.FullTimestamp = true
	log.Formatter = &TranslationFormatter{formatter}
}

// TranslationFormatter implementation from prefixed.TextFormatter
type TranslationFormatter struct {
	*prefixed.TextFormatter
}

// Format returns
func (t *TranslationFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if code, ok := entry.Data["code"]; ok {
		if translation, ok := translations[code.(string)]; ok {
			entry.Message = translation
		}
	}
	return t.TextFormatter.Format(entry)
}

// RawFormatter a alias of prefixed.TextFormatter
type RawFormatter struct{}

// Format returns
func (f *RawFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

func init() {
	formatter := new(prefixed.TextFormatter)
	formatter.TimestampFormat = `Jan 02 15:04:05`
	formatter.FullTimestamp = true

	log.Formatter = formatter
	rawLog.Formatter = new(RawFormatter)
}

// Get returns log
func Get() *logrus.Logger {
	switch strings.ToLower(os.Getenv("TYK_LOGLEVEL")) {
	case "error":
		log.Level = logrus.ErrorLevel
	case "warn":
		log.Level = logrus.WarnLevel
	case "debug":
		log.Level = logrus.DebugLevel
	default:
		log.Level = logrus.InfoLevel
	}
	return log
}

// GetRaw returns rawLog
func GetRaw() *logrus.Logger {
	return rawLog
}
