package migrationsv2

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const (
	up13 = `
UPDATE keys SET next_nonce = 0 WHERE next_nonce IS NULL;
ALTER TABLE keys ALTER COLUMN next_nonce SET NOT NULL, ALTER COLUMN next_nonce SET DEFAULT 0;
`
	down13 = `
ALTER TABLE keys ALTER COLUMN next_nonce SET DEFAULT NULL;
`
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0013_not_null_keys_next_nonce",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up13).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down13).Error
		},
	})
}
