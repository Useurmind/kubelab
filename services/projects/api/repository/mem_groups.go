package repository

import (
	"context"
	"fmt"

	"github.com/useurmind/kubelab/services/projects/api/models"
)

// MemGroupRepo is an implementation of the GroupRepo interface to store groups in memory.
type MemGroupRepo struct {
	nextGroupId int64
	groups map[int64]*models.Group
}

func NewMemGroupRepo() *MemGroupRepo {
	return &MemGroupRepo{
		nextGroupId: 1,
		groups: make(map[int64]*models.Group),
	}
}

func (r *MemGroupRepo) getNextId() int64 {
	r.nextGroupId++

	return r.nextGroupId
}

func (r *MemGroupRepo) CreateOrUpdate(ctx context.Context, group *models.Group) (*models.Group, error) {
	// give all new subgroups a unique id
	newSubgroups := group.GatherNewSubgroups()
	for _, newSubgroup := range newSubgroups {
		nextID := r.getNextId()
		newSubgroup.Id = nextID
	}

	if group.IsNew() {
		// insert
		nextID := r.getNextId()
		group.Id = nextID
		r.groups[nextID] = group
	} else {
		// update
		r.groups[group.Id] = group
	}

	return group, nil
}

func (r *MemGroupRepo) Get(ctx context.Context, groupID int64) (*models.Group, error) {
	group, ok := r.groups[groupID]
	if !ok {
		return nil, fmt.Errorf("Could not find group in memory store")
	}

	return group, nil
}

func (r *MemGroupRepo) List(ctx context.Context, startIndex int64, count int64) ([]*models.Group, error) {
	groups := make([]*models.Group, 0)
	var index int64 = 0

	for _, group := range r.groups {
		if index >= startIndex && index < startIndex + count {
			groups = append(groups, group)
		}
	}

	return groups, nil
}

func (r *MemGroupRepo) Delete(ctx context.Context, groupID int64) error {
	delete(r.groups, groupID)

	return nil
}

func (r *MemGroupRepo) Close() error {
	return nil
}