package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)


type pgSequenceVal struct {
	Nextval int64 `db:"nextval"`
}

func getNextValFromSequence(ctx context.Context, tx *sqlx.Tx, sequenceName string) (nextval int64, err error) {
	pgSequenceVal := pgSequenceVal{}
	err = tx.GetContext(ctx, &pgSequenceVal, "SELECT nextval($1)", sequenceName)
	if err != nil {
		return 0, err
	}

	return pgSequenceVal.Nextval, nil
}
