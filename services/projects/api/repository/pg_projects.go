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

type pgProject struct {
	ID int64 `db:"id"`
	GroupID int64 `db:"group_id"`
	Slug string `db:"slug"`

	Data types.JSONText `db:"data"`
}

func newPGProject(project *models.Project) (pgproject *pgProject, err error) {
	pgproject = &pgProject{
		ID: project.Id,
		GroupID: project.GroupId,
		Slug: project.Slug,
	}

	bytes, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}

	pgproject.Data = types.JSONText(bytes)

	return pgproject, nil
}

func (pgproject *pgProject) createModel() (*models.Project, error) {
	project := models.Project{
		Id: pgproject.ID,
		GroupId: pgproject.GroupID,
		Slug: pgproject.Slug,
	}
	err := pgproject.Data.Unmarshal(&project)
	if err != nil {
		log.Error().
			Err(err).
			Int64("projectID", pgproject.ID).
			Str("data", pgproject.Data.String()).
			Msg("Could not unmarshal project data from json")
		return nil, err
	}

	return &project, nil
}

// PGProjectRepo is an implementation of the ProjectRepo interface to store projects in a postgres database.
type PGProjectRepo struct {
	tx *sqlx.Tx
}

func (r *PGProjectRepo) CreateOrUpdate(ctx context.Context, project *models.Project) (*models.Project, error) {
	pgproject, err := newPGProject(project)
	if err != nil {
		return nil, err
	}

	if project.IsNew() {
		// insert
		res, err := r.tx.NamedExecContext(ctx, "INSERT INTO projects (name, group_id, slug, data) VALUES (:name, :groupid, :slug, :data)", pgproject)
		if err != nil {
			return nil, err
		}

		id, _ := res.LastInsertId()
		project.Id = id
	} else {
		// update
		res, err := r.tx.NamedExecContext(ctx, "UPDATE projects SET name=:name, group_id=:groupid slug=:slug, data=:data WHERE id=:id", pgproject)
		if err != nil {
			return nil, err
		}

		if rows, _ := res.RowsAffected(); rows == 0 {
			return nil, fmt.Errorf("Could not update project %d, affected rows 0", project.Id)
		}
	}

	return project, nil
}

func (r *PGProjectRepo) Get(ctx context.Context, projectID int64) (*models.Project, error) {
	pgproject := pgProject{}
	err := r.tx.GetContext(ctx, &pgproject, "SELECT * FROM projects WHERE id = $1 LIMIT 1", projectID)
	if err != nil {
		return nil, err
	}

	return pgproject.createModel()
}

func (r *PGProjectRepo) GetBySlugs(ctx context.Context, groupSlug string, projectSlug string) (*models.Project, error) {
	pgproject := pgProject{}
	err := r.tx.GetContext(ctx, &pgproject, "SELECT p.* FROM projects p INNER JOIN groups g ON g.slug = $2 AND g.id = p.group_id WHERE p.slug = $1 LIMIT 1",
		projectSlug, groupSlug)
	if err != nil {
		return nil, err
	}

	return pgproject.createModel()
}

func (r *PGProjectRepo) GetByGroupID(ctx context.Context, groupID int64) ([]*models.Project, error) {
	pgprojects := make([]*pgProject, 0)
	err := r.tx.GetContext(ctx, &pgprojects, "SELECT id, slug, data -> 'name' as name FROM projects WHERE group_id = $1", groupID)
	if err != nil {
		return nil, err
	}

	projects := make([]*models.Project, 0, len(pgprojects))
	for _, pgProj := range pgprojects {
		proj, err := pgProj.createModel()
		if err != nil {
			return nil, err
		}
		projects = append(projects, proj)
	}

	return projects, nil
}

func (r *PGProjectRepo) CountByGroupID(ctx context.Context, groupID int64) (int64, error) {
	var count int64
	err := r.tx.GetContext(ctx, &count, "SELECT COUNT(id) FROM projects WHERE group_id = $1", groupID)
	if err != nil {
		return -1, err
	}

	return count, nil
}


func (r *PGProjectRepo) Delete(ctx context.Context, projectID int64) error {
	res, err := r.tx.ExecContext(ctx, "DELETE FROM projects WHERE id = $1", projectID)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("Could not delete project %d, affected rows 0", projectID)
	}

	return nil
}