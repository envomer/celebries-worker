package db

import (
	logger2 "any-days.com/celebs/logger"
	"github.com/evolidev/evoli/framework/filesystem"
	"github.com/evolidev/evoli/framework/logging"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Manager struct {
	db   *gorm.DB
	path string
}

var mgr *Manager

func Mgr() *Manager {
	if mgr == nil {
		mgr = newManager()
	}
	return mgr
}

func Db() *gorm.DB {
	return Mgr().db
}

func ResetDb() *gorm.DB {
	mgr = nil

	return Db()
}

func Clean() {
	driver := os.Getenv("DB_DRIVER")
	if driver == "postgres" {
		db := Db()
		db.Exec("DROP SCHEMA public CASCADE")
		db.Exec("CREATE SCHEMA public")

		logging.Info("Postgres DB cleaned")
	}

	filesystem.Delete("db.db")

	mgr = nil
}

func newManager() *Manager {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             5 * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Silent,   // Log level
			IgnoreRecordNotFoundError: false,           // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,            // Disable color
		},
	)

	driver := os.Getenv("DB_DRIVER")

	logger2.AppLog.Debug("DB_DRIVER: %s", driver)

	var dialect gorm.Dialector
	if driver == "postgres" {
		dsn := os.Getenv("DB_DSN")

		if dsn == "" {
			logging.Fatal("DB_DSN is empty")
		}

		dialect = postgres.Open(dsn)
	} else {
		dialect = sqlite.Open("db.db")
	}

	dbInstance, err := gorm.Open(dialect, &gorm.Config{
		Logger:                 newLogger,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		panic(err)
	}

	return &Manager{db: dbInstance}
}
