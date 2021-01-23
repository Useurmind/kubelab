package models

import "fmt"

// IsNew checks whether the group is new by testing if id is not set.
func (g *Group) IsNew() bool {
	return g.Id == 0
}

// InsertProjectRef tries to insert a reference to the given project
// on the matching group/subgroup. If no matching group/subgroup is found
// an error is returned.
func (g *Group) InsertProjectRef(project *Project) error {
	projectRef := &ProjectRef{
		Id: project.Id,
		Slug: project.Slug,
		Name: project.Name,
	}
	if g.Id == project.AssignedGroupId {
		g.Projects = append(g.Projects, projectRef)
		return nil
	}

	subgroups := g.GatherSubgroups()

	for _, subgroup := range subgroups {
		if subgroup.Id == project.AssignedGroupId {
			subgroup.Projects = append(subgroup.Projects, projectRef)
			return nil
		}
	}

	return fmt.Errorf("Could not find subgroup with matching ID %d in group %s", project.AssignedGroupId, g.Name)
}

// UpdateProjectRef tries to update a reference to the given project.
// If no matching project is found an error is returned.
func (g *Group) UpdateProjectRef(project *Project) error {
	projects := g.GatherProjects()

	for _, projRef := range projects {
		if project.Id == projRef.Id {
			projRef.Name = project.Name
			projRef.Slug = project.Slug
			return nil
		}
	}

	return fmt.Errorf("Could not find project ref for project %s (%d) in group %s (%d)", project.Name, project.Id, g.Name, g.Id)
}

// DeleteProjectRef tries to delete a reference to the given project.
// If no matching project is found an error is returned.
func (g *Group) DeleteProjectRef(project *Project) error {
	assGroup := g.GetGroupByID(project.AssignedGroupId)
	if assGroup == nil {
		return fmt.Errorf("Could not find assigned group for project %s (%d) in group %s (%d)", project.Name, project.Id, g.Name, g.Id)
	}

	for i, projRef := range assGroup.Projects {
		if projRef.Id == project.Id {
			assGroup.Projects[i] = assGroup.Projects[len(assGroup.Projects) - 1]
			assGroup.Projects = assGroup.Projects[:len(assGroup.Projects) - 1]
			return nil
		}
	}

	return fmt.Errorf("Could not delete ref for project %s (%d) in group %s (%d)", project.Name, project.Id, g.Name, g.Id)
}

// GetGroupByID returns the group in the group tree with the matching id or nil if not found.
func (g *Group) GetGroupByID(groupID int64) *Group {
	if groupID == g.Id {
		return g
	}

	subgroups := g.GatherSubgroups()
	for _, subgroup := range subgroups {
		if subgroup.Id == groupID {
			return subgroup
		}
	}

	return nil
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