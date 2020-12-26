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

// PGRepoFactory is used to create postgres repositories.
type PGRepoFactory struct {
	connectionString string
}

// NewPGRepoFactory returns a repo factory for the given connection information.
func NewPGRepoFactory(host string, port string, db string, user string, password string) *PGRepoFactory {
	return &PGRepoFactory{
		connectionString: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", 
			host, port, user, password, db),
	}
}

func (f *PGRepoFactory) createConnection(ctx context.Context) (*sqlx.DB, error) {
	return sqlx.ConnectContext(ctx, "postgres", f.connectionString)
}

func (f *PGRepoFactory) Migrate() error {
	db, err := f.createConnection(context.Background())
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

func (f *PGRepoFactory) GetGroupRepo(ctx context.Context) (GroupRepo, error) {
	db, err := f.createConnection(ctx)
	if err != nil {
		return nil, err
	}
	return &PGGroupRepo{
		db: db,
	}, nil
}
