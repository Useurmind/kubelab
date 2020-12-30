package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/useurmind/kubelab/services/projects/api/models"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
)

type pgGroup struct {
	ID int64 `db:"id"`
	Name string `db:"name"`
	Data types.JSONText `db:"data"`
}

// PGGroupRepo is an implementation of the GroupRepo interface to store groups in a postgres database.
type PGGroupRepo struct {
	db *sqlx.DB
}

func (r *PGGroupRepo) CreateOrUpdate(ctx context.Context, group *models.Group) (*models.Group, error) {
	pggroup := pgGroup{
		ID: group.Id,
		Name: group.Name,
	}

	bytes, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	pggroup.Data = types.JSONText(bytes)

	if (pggroup.ID == 0) {
		// insert
		res, err := r.db.NamedExec("INSERT INTO groups (name, data) VALUES (:name, :data)", pggroup)
		if err != nil {
			return nil, err
		}

		id, _ := res.LastInsertId()
		group.Id = id
	} else {
		// update
		res, err := r.db.NamedExec("UPDATE groups SET name=:name, data=:data WHERE id=:id", pggroup)
		if err != nil {
			return nil, err
		}

		if rows, _ := res.RowsAffected(); rows == 0 {
			return nil, fmt.Errorf("Could not update group %d, affected rows 0", group.Id)
		}
	}

	return group, nil
}

func (r *PGGroupRepo) Get(ctx context.Context, groupID int) (*models.Group, error) {
	pggroup := pgGroup{}
	err := r.db.GetContext(ctx, &pggroup, "SELECT * FROM groups WHERE id = $1 LIMIT 1", groupID)
	if err != nil {
		return nil, err
	}

	group := models.Group{
		Id: pggroup.ID,
		Name: pggroup.Name,
	}
	err = pggroup.Data.Unmarshal(&group)
	if err != nil {
		log.Error().
			Err(err).
			Int("groupID", groupID).
			Str("data", pggroup.Data.String()).
			Msgf("Could not unmarshal group data of group from json")
		return nil, err
	}

	return &group, nil
}

func (r *PGGroupRepo) List(ctx context.Context, startIndex int, count int) ([]*models.Group, error) {
	pggroups := make([]pgGroup, 0)
	err := r.db.SelectContext(ctx, &pggroups, "SELECT * FROM groups LIMIT $1 OFFSET $2", count, startIndex)
	if err != nil {
		return nil, err
	}

	groups := make([]*models.Group, len(pggroups))
	for i, pggroup := range pggroups {
		group := models.Group{
			Id: pggroup.ID,
			Name: pggroup.Name,
		}
		err = pggroup.Data.Unmarshal(&group)
		if err != nil {
			log.Error().
				Err(err).
				Int("index", i).
				Str("data", pggroup.Data.String()).
				Msgf("Could not unmarshal group data of group from json")
			return nil, err
		}

		groups[i] = &group
	}

	return groups, nil
}

func (r *PGGroupRepo) Delete(ctx context.Context, groupID int) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM groups WHERE id = $1", groupID)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("Could not delete group %d, affected rows 0", groupID)
	}

	return nil
}