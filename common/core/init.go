package core

import (
	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
	"mojito-coding-test/common/config"
)

const configFilePrefix string = "config"

var Logrus = newLogger()
var Logger *logrus.Entry // Set during boot.
var Config = newConfig()
var Validator = newValidator()

func newLogger() *logrus.Logger {
	return logrus.New()
}

func newConfig() *config.Config {
	cfg, err := config.New(configFilePrefix)
	if err != nil {
		Logrus.Fatalf("failed to load config - %+v", err)
	}

	return cfg
}

func Populate(values ...interface{}) {
	// Append core
	values = append(values, Logrus)
	values = append(values, Logger)
	values = append(values, Config)
	values = append(values, Validator)

	err := inject.Populate(values...)
	if err != nil {
		Logger.Fatal(err)
	}
}
