package env

import (
	"os"
	"strconv"
)

type loggerConfig struct {
	level  string
	asJson bool
}

func NewLoggerConfig() (*loggerConfig, error) {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	asJsonStr := os.Getenv("LOG_AS_JSON")
	asJson, err := strconv.ParseBool(asJsonStr)
	if err != nil {
		asJson = false
	}

	return &loggerConfig{
		level:  level,
		asJson: asJson,
	}, nil
}

func (c *loggerConfig) Level() string {
	return c.level
}

func (c *loggerConfig) AsJson() bool {
	return c.asJson
}
