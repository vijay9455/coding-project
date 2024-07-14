package migration

import (
	"calendly/lib/logger"
	"context"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func Run(dbConn *gorm.DB) {
	sqlConn, err := dbConn.DB()
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(sqlConn, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://db/migrations", "postgres", driver)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info(context.TODO(), "no new changes, so skipping migration", nil)
			return
		}
		logger.Error(context.TODO(), "error while running migration", nil)
		panic(err)
	}

	logger.Info(context.TODO(), "Successfully ran migration!", nil)
}
