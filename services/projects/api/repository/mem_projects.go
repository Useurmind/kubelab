package repository

import (
	"context"

	"github.com/useurmind/kubelab/services/projects/api/models"
)

// MemGroupRepo is an implementation of the GroupRepo interface to store groups in memory.
type MemProjectRepo struct {
	nextId int64
	projects map[int64]*models.Project
}

func NewMemProjectRepo() *MemProjectRepo {
	return &MemProjectRepo{
		nextId: 1,
		projects: make(map[int64]*models.Project),
	}
}

func (r *MemProjectRepo) getNextId() int64 {
	r.nextId++

	return r.nextId
}

func (r *MemProjectRepo) CreateOrUpdate(ctx context.Context, project *models.Project) (*models.Project, error) {
	if project.IsNew() {
		// insert
		nextID := r.getNextId()
		project.Id = nextID
		r.projects[nextID] = project
	} else {
		// update
		r.projects[project.Id] = project
	}

	return project, nil
}

func (r *MemProjectRepo) Get(ctx context.Context, projectID int64) (*models.Project, error) {
	project, ok := r.projects[projectID]
	if !ok {
		return nil, nil
	}

	return project, nil
}

func (r *MemProjectRepo) GetByGroupID(ctx context.Context, groupID int64) ([]*models.Project, error) {
	projects := make([]*models.Project, 0)

	for _, proj := range r.projects {
		if proj.GroupId == groupID {
			projects = append(projects, proj)
		}
	}

	return projects, nil
}

func (r *MemProjectRepo) CountByGroupID(ctx context.Context, groupID int64) (int64, error) {
	projects, _ := r.GetByGroupID(ctx, groupID)
	return int64(len(projects)), nil
}

func (r *MemProjectRepo) Delete(ctx context.Context, projectID int64) error {
	delete(r.projects, projectID)

	return nil
}

func (r *MemProjectRepo) Close() error {
	return nil
}