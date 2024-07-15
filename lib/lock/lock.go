package lock

import (
	"context"

	"gorm.io/gorm"
)

func WithLock(ctx context.Context, db *gorm.DB, resourceID string, exeFunc func(*gorm.DB) error) error {
	return db.Connection(func(tx *gorm.DB) error {
		defer tx.Exec("select pg_advisory_unlock(hashtext(?))", resourceID)
		tx.Exec("select pg_advisory_lock(hashtext(?))", resourceID)
		return exeFunc(db)
	})
}
