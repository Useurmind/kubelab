package models

// HasAnyProjects checks if the group or any subgroups has any projects in it.
func (g *Group) HasAnyProjects() bool {
	hasProjects := len(g.Projects) > 0

	for _, subgroup := range g.Subgroups {
		hasProjects = hasProjects || subgroup.HasAnyProjects()
	}

	return hasProjects
}
