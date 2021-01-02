package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)


type pgSequenceVal struct {
	Nextval int64 `db:"nextval"`
}

func getNextValFromSequence(ctx context.Context, db *sqlx.DB, sequenceName string) (nextval int64, err error) {
	pgSequenceVal := pgSequenceVal{}
	err = db.GetContext(ctx, &pgSequenceVal, "SELECT nextval($1)", sequenceName)
	if err != nil {
		return 0, err
	}

	return pgSequenceVal.Nextval, nil
}
