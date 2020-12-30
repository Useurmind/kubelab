package models

// A group that contains projects or can be used as hierarchical element to structure the projects.
type Group struct {
	// The primary key of the group.
	Id int64 `json:"id,omitempty"`
	// The name of the group.
	Name string `json:"name,omitempty"`
	// The groups under this group in the group hierarchy.
	Subgroups []Group `json:"subgroups,omitempty"`
	// Fat references to the projects in this group.
	Projects []ProjectRef `json:"projects,omitempty"`
}

// A pointer to a project.
type ProjectRef struct {
	// The primary key of the referenced project.
	Id int64 `json:"id,omitempty"`
	// The slug of the project.
	Slug string `json:"slug,omitempty"`
	// The name of the project.
	Name string `json:"name,omitempty"`
}

// 
type Project struct {
	// The primary key of this project.
	Id int64 `json:"id,omitempty"`
	// The root group to which this project belongs.
	GroupId int64 `json:"groupId,omitempty"`
	// The name of the (sub)group to which this project belongs.
	AssignedGroupName string `json:"assignedGroupName,omitempty"`
	// Pretty name for this project.
	Name string `json:"name,omitempty"`
	// The short name of this project. Must only contain numbers, letters, dash and underline.
	Slug string `json:"slug,omitempty"`
}