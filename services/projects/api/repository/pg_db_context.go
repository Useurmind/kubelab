package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PGDBSystem is used to manage db contexts for http requests.
type PGDBSystem struct {
	connectionString string
}

func NewPGDBSystem(host string, port string, db string, user string, password string) *PGDBSystem {
	return &PGDBSystem{
		connectionString: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
			host, port, user, password, db),
	}
}

func (s *PGDBSystem) NewContext() DBContext {
	return &PGDBContext{
		connectionString: s.connectionString,
	}
}

// PGDBContext is used to create postgres repositories.
type PGDBContext struct {
	connectionString string
	db *sqlx.DB 
	tx *sqlx.Tx
}

// newConnection returns a new db connection that must be closed by the caller.
func (c *PGDBContext) newConnection(ctx context.Context) (*sqlx.DB, error) {
	return sqlx.ConnectContext(ctx, "postgres", c.connectionString)
}

// getSharedConnection returns the shared db connection that is closed when the context is closed when the http request ends.
func (f *PGDBContext) getSharedConnection(ctx context.Context) (*sqlx.Tx, error) {
	if f.db != nil {
		return f.tx, nil
	}

	db, err := f.newConnection(ctx)
	if err != nil {
		return nil, err
	}

	f.db = db
	f.tx, err = db.BeginTxx(ctx, nil)

	return f.tx, err
}

func (f *PGDBContext) Migrate() error {
	db, err := f.newConnection(context.Background())
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})

	exe, err := os.Executable()
	if err != nil {
		return err
	}
	exePath := filepath.Dir(exe)
	migrationPath := filepath.ToSlash(filepath.Join(exePath, "db"))

	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationPath), "postgres", driver)
	if err != nil {
		return err
	}
	defer m.Close()

	log.Info().Msgf("Migrating database to current version")
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		// ErrNoChange is ignored because that means there is just no work to be done
		return err
	}

	return nil
}

func (f *PGDBContext) GetGroupRepo(ctx context.Context) (GroupRepo, error) {
	tx, err := f.getSharedConnection(ctx)
	if err != nil {
		return nil, err
	}
	return &PGGroupRepo{
		tx: tx,
	}, nil
}

func (f *PGDBContext) GetProjectRepo(ctx context.Context) (ProjectRepo, error) {
	tx, err := f.getSharedConnection(ctx)
	if err != nil {
		return nil, err
	}
	return &PGProjectRepo{
		tx: tx,
	}, nil
}

func (r *PGDBContext) Commit() error {
	return r.tx.Commit()
}

func (r *PGDBContext) Rollback() error {
	return r.tx.Rollback()
}

func (r *PGDBContext) Close() error {
	return r.db.Close()
}