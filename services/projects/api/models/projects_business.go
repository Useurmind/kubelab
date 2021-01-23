package models

// IsNew checks whether the project is new by testing if id is not set.
func (p *Project) IsNew() bool {
	return p.Id == 0
}

// ToProjectRef converts the project into a project reference.
func (p *Project) ToProjectRef() *ProjectRef {
	return &ProjectRef{
		Id:   p.Id,
		Name: p.Name,
		Slug: p.Slug,
	}
}
