package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"time"

	_ "github.com/lib/pq"
)

const DefaultDSN = "postgres://postgres:postgres@localhost:5432/oee?sslmode=disable"

func Open(dsn string) (*sql.DB, error) {
	if dsn == "" {
		dsn = DefaultDSN
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(30 * time.Second)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

//go:embed sql/schema.sql
var schema string

func Migrate(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx,
		schema,
	); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

type Repository struct {
	*WorkshiftRepository
	*EquipmentRepository
	*EquipmentLineRepository
	*DowntimeRepository
	*SpecificationsRepository
	*QualificationRepository
	*QualificationRuleRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		WorkshiftRepository:         NewWorkshiftRepository(db),
		EquipmentRepository:         NewEquipmentRepository(db),
		EquipmentLineRepository:     NewEquipmentLineRepository(db),
		DowntimeRepository:          NewDowntimeRepository(db),
		SpecificationsRepository:    NewSpecificationsRepository(db),
		QualificationRepository:     NewQualificationRepository(db),
		QualificationRuleRepository: NewQualificationRuleRepository(db),
	}
}
