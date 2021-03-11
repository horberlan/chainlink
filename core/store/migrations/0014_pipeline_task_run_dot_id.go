package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

const (
	up14 = `
ALTER TABLE pipeline_task_runs ADD COLUMN dot_id text NOT NULL; 
ALTER TABLE pipeline_task_runs ALTER COLUMN pipeline_task_spec_id DROP NOT NULL; 
ALTER TABLE pipeline_task_runs DROP CONSTRAINT pipeline_task_runs_pipeline_task_spec_id_fkey; 
`
	down14 = `
ALTER TABLE pipeline_task_runs DROP COLUMN dot_id;
`
)

func init() {
	Migrations = append(Migrations, &gormigrate.Migration{
		ID: "0014_pipeline_task_run_dot_id",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up14).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down14).Error
		},
	})
}
