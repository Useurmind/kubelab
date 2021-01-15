package repository

import (
	"context"
)

// MemDBSystem is used to manage MemDBContexts.
type MemDBSystem struct {
	// we use the same for all requests
	// currently its not thread safe, so its only fit for simple testing
	context *MemDBContext
}

// NewMemDBSystem returns a db system for in memory repos.
func NewMemDBSystem() *MemDBSystem {
	return &MemDBSystem{
		context: &MemDBContext{},
	}
}

func (s *MemDBSystem) NewContext() DBContext {
	return s.context
}

// MemDBContext is used to create in memory repositories.
type MemDBContext struct {
	groups *MemGroupRepo
}

func (f *MemDBContext) Migrate() error {
	return nil
}

func (f *MemDBContext) GetGroupRepo(ctx context.Context) (GroupRepo, error) {
	if f.groups == nil {
		f.groups = NewMemGroupRepo()
	}

	return f.groups, nil
}

func (f *MemDBContext) Close() error {
	return nil
}
