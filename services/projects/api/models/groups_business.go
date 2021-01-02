package models

// IsNew checks whether the group is new by testing if id is not set.
func (g *Group) IsNew() bool {
	return g.Id == 0
}

// HasAnyProjects checks if the group or any subgroups has any projects in it.
func (g *Group) HasAnyProjects() bool {
	hasProjects := len(g.Projects) > 0

	for _, subgroup := range g.Subgroups {
		hasProjects = hasProjects || subgroup.HasAnyProjects()
	}

	return hasProjects
}

// GatherProjects returns a list of all projects in the group.
func (g *Group) GatherProjects() []*ProjectRef {
	projects := make([]*ProjectRef, 0)

	projects = append(projects, g.Projects...)

	for _, subgroup := range g.Subgroups {
		projects = append(projects, subgroup.GatherProjects()...)
	}

	return projects
}

// GatherSubgroups returns a list of all subgroups in the group.
func (g *Group) GatherSubgroups() []*Group {
	groups := make([]*Group, 0)

	groups = append(groups, g.Subgroups...)

	for _, subgroup := range g.Subgroups {
		groups = append(groups, subgroup.GatherSubgroups()...)
	}

	return groups
}

// GatherNewSubgroups returns a list of all subgroups in the group that are new.
func (g *Group) GatherNewSubgroups() []*Group {
	allgroups := g.GatherSubgroups()
	newgroups := make([]*Group, 0)

	for _, g := range allgroups {
		if g.IsNew() {
			newgroups = append(newgroups, g)
		}
	}

	return newgroups
}

// GetNewAndDeletedProjects returns a list of projects that would need to be created or
// deleted when this new group will be saved.
// Returns potentially created projects, potentially deleted projects.
func (g *Group) GetNewAndDeletedProjects(oldGroup *Group) ([]*ProjectRef, []*ProjectRef) {	
	createdProjects := make([]*ProjectRef, 0)
	deletedProjects := make([]*ProjectRef, 0)

	newProjects := g.GatherProjects()
	oldProjects := oldGroup.GatherProjects()

	newProjectsByID, newProjectsWithoutID := getProjectMapByID(newProjects)
	oldProjectsByID, _ := getProjectMapByID(oldProjects)

	// all projects without ID are new
	createdProjects = append(createdProjects, newProjectsWithoutID...)

	// find projects that have an ID and are new
	for _, newProject := range newProjects {
		if newProject.Id != 0 {
			_, found := oldProjectsByID[newProject.Id]
			if !found {
				createdProjects = append(createdProjects, newProject)
			}
		}
	}

	// find projects that are not there anymore
	for _, oldProject := range oldProjects {
		_, found := newProjectsByID[oldProject.Id]
		if !found {
			deletedProjects = append(deletedProjects, oldProject)
		}
	}

	return createdProjects, deletedProjects
}

// getProjectMapById returns the projects in a map by their id and the projects without id.
func getProjectMapByID(projects []*ProjectRef) (map[int64]*ProjectRef, []*ProjectRef) {
	projectMap := make(map[int64]*ProjectRef)
	projectsWithoutID := make([]*ProjectRef, 0)

	for _, project := range projects {
		if project.Id == 0 {
			projectsWithoutID = append(projectsWithoutID, project)
		} else {
			projectMap[project.Id] = project
		}
	}

	return projectMap, projectsWithoutID
}