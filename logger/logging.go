package logger

import "github.com/evolidev/evoli/framework/logging"

var AppLog = logging.NewLogger(&logging.Config{
	Name:        "app",
	PrefixColor: 165,
	//Path:        "log.log",
})

var WebLog = logging.NewLogger(&logging.Config{
	Name:        "web",
	PrefixColor: 169,
	//Path:        "log.log",
})

var VersionLog = logging.NewLogger(&logging.Config{
	Name:        "version",
	PrefixColor: 145,
	//Path:        "log.log",
})

var CronLog = logging.NewLogger(&logging.Config{
	Name:        "cron",
	PrefixColor: 135,
	//Path:        "cron.log",
})
