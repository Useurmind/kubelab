package models

// IsNew checks whether the project is new by testing if id is not set.
func (p *Project) IsNew() bool {
	return p.Id == 0
}