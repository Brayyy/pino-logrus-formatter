package pinoformat

import (
	"encoding/json"
	"fmt"
	"os"

	logrus "github.com/sirupsen/logrus"
)

var log *logrus.Entry

// PinoFormatter ...
type PinoFormatter struct {
	Name string
	Base map[string]interface{}
}

var (
	levelMap       map[logrus.Level]int
	cachedPid      int
	cachedHostname string
)

func init() {
	levelMap = map[logrus.Level]int{
		logrus.TraceLevel: 10,
		logrus.DebugLevel: 20,
		logrus.InfoLevel:  30,
		logrus.WarnLevel:  40,
		logrus.ErrorLevel: 50,
		logrus.FatalLevel: 60,
		logrus.PanicLevel: 60,
	}
	cachedPid = os.Getpid()
	cachedHostname, _ = os.Hostname()
}

// Format ...
func (f *PinoFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	pinoEntry := map[string]interface{}{
		"hostname":  cachedHostname,
		"level":     levelMap[entry.Level],
		"msg":       entry.Message,
		"pid":       cachedPid,
		"timestamp": float64(entry.Time.UnixNano() / 1e6),
	}
	// Name may not be set
	if f.Name != "" {
		pinoEntry["name"] = f.Name
	}
	// Add all additional fields, skip already defined keys
	for key, value := range entry.Data {
		if _, ok := pinoEntry[key]; !ok {
			pinoEntry[key] = value
		}
	}
	// JSON encode safely
	serialized, err := json.Marshal(pinoEntry)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
