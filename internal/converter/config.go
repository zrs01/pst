package converter

import (
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

type Config struct {
	FontFamily string `yaml:"fontfamily,omitempty"`
	FontSize   int    `yaml:"fontsize,omitempty"`
	Logging    struct {
		Level string
	}
}

// NewConfig creates new instance of the configuration from the file
func NewConfig(f string) (*Config, error) {
	cfg := &Config{}
	cfg.FontFamily = "Arial"
	cfg.FontSize = 10
	cfg.Logging.Level = "INFO"

	if f != "" {
		if err := cfg.load(f); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}

func (c *Config) load(cf string) error {
	if err := configor.Load(c, cf); err != nil {
		return err
	}
	logrus.SetLevel(c.getLevel())
	return nil
}

// LogLevel returns the parsed logging level for logrus
func (c *Config) getLevel() logrus.Level {
	level, err := logrus.ParseLevel(c.Logging.Level)
	if err != nil {
		logrus.Warnf("failed to parse log level for %s, use INFO level", err)
		return logrus.InfoLevel
	}
	return level
}
