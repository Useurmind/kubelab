package repository

import (
	"context"

	"github.com/useurmind/kubelab/services/projects/api/models"
)

// DBSystem is represents the chosen DB technology to use.
// There is only one DBSystem per running service instance.
type DBSystem interface {
	// NewContext returns a new DBContext for the current http request.
	NewContext() DBContext
}

// DBContext is the root entity for accessing the database inside one http request.
// It can be used to retrieve the different repos.
// One new DBContext is created per request.
type DBContext interface {
	// Migrate the database to the current schema version.
	Migrate() error

	// CreateTransaction returns a transaction that can be used to make multiple repo calls
	// atomic.
	CreateTransaction(ctx context.Context) (DBTransaction, error)

	// GetGroupRepo returns a group repository.
	GetGroupRepo(ctx context.Context) (GroupRepo, error)

	// GetProjectRepo returns a project repository.
	GetProjectRepo(ctx context.Context) (ProjectRepo, error)

	// Commit all changes on the shared connection.
	Commit() error

	// Rollback all changes on the shared connection.
	Rollback() error

	// Close all connections to the database.
	Close() error
}

// GroupRepo is an interface for a group repository.
type GroupRepo interface {
	// CreateOrUpdate creates or updates a group depending on whether the primary key is already set.
	// It returns the group as saved in the database.
	CreateOrUpdate(ctx context.Context, group *models.Group) (*models.Group, error)

	// Get retrieves the group with the given id from the database.
	Get(ctx context.Context, groupID int64) (*models.Group, error)

	// List retrieves number groups from the database starting with the given index.
	List(ctx context.Context, startIndex int64, count int64) ([]*models.Group, error)

	// Delete the group with the given id.
	Delete(ctx context.Context, groupID int64) error
}

// ProjectRepo is an interface for a project repository.
type ProjectRepo interface {
	// CreateOrUpdate creates or updates a project depending on whether the primary key is already set.
	// It returns the project as saved in the database.
	CreateOrUpdate(ctx context.Context, project *models.Project) (*models.Project, error)

	// Get retrieves the project with the given id from the database.
	Get(ctx context.Context, projectID int64) (*models.Project, error)

	// Delete the project with the given id.
	Delete(ctx context.Context, projectID int64) error
}
