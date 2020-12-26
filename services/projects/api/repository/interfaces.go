package repository

import (
	"context"

	"github.com/useurmind/kubelab/services/projects/api/models"
)

// RepoFactory is used to get repositories.
type RepoFactory interface {
	// GetGroupRepo returns a group repository.
	GetGroupRepo(ctx context.Context) (GroupRepo, error)
}

// GroupRepo is an interface for a group repository.
type GroupRepo interface {
	// CreateOrUpdate creates or updates a group depending on whether the primary key is already set.
	// It returns the group as saved in the database.
	CreateOrUpdate(ctx context.Context, group *models.Group) (*models.Group, error)

	// Get retrieves the group with the given id from the database.
	Get(ctx context.Context, groupID int) (*models.Group, error)

	// List retrieves number groups from the database starting with the given index.
	List(ctx context.Context, startIndex int, count int) ([]*models.Group, error)

	// Delete the group with the given id.
	Delete(ctx context.Context, groupID int) error
}