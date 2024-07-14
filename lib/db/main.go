package db

import (
	"calendly/lib/logger"
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func Connect(dbUrl string, maxIdleConnections, maxOpenConnections int, enableLogs bool) {
	db = establishConnection(dbUrl, enableLogs)

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(maxIdleConnections)
	sqlDB.SetMaxOpenConns(maxOpenConnections)
	logger.Info(context.TODO(), "Successfully connected to db", nil)
}

func establishConnection(dbUrl string, enableLogs bool) *gorm.DB {
	gormCfg := &gorm.Config{}
	if enableLogs {
		gormCfg.Logger = gormLogger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			gormLogger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  gormLogger.Info,
				IgnoreRecordNotFoundError: false,
				Colorful:                  false,
			},
		)
	}

	pgConfig := postgres.Config{
		// DriverName: "nrpgx",
		DSN: dbUrl,
	}

	dbConnection, err := gorm.Open(postgres.New(pgConfig), gormCfg)
	if err != nil {
		logger.Error(context.TODO(), "error while connecting to db", map[string]any{"error": err})
		panic(err)
	}

	return dbConnection
}

func Get() *gorm.DB {
	return db
}

func Close() {
	sqlDb, err := db.DB()
	if err != nil {
		logger.Error(context.TODO(), "error while fetching db to close db connection", map[string]any{"error": err})
		panic(err)
	}

	sqlDb.Close()
}
