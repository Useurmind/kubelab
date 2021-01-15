

// A group that contains projects or can be used as hierarchical element to structure the projects.
export interface IGroup {
	// The primary key of the group.
	id?: number;
	// The name of the group.
	name?: string;
	// The short name of the group. Must only contain numbers, letters, dash and underline.
	slug?: string;
	// The groups under this group in the group hierarchy.
	subgroups?: IGroup[];
	// Fat references to the projects in this group.
	projects?: IProjectRef[];
}

// A pointer to a project.
export interface IProjectRef {
	// The primary key of the referenced project.
	id?: number;
	// The slug of the project.
	slug?: string;
	// The name of the project.
	name?: string;
}

// 
export interface IProject {
	// The primary key of this project.
	id?: number;
	// The root group to which this project belongs.
	groupId?: number;
	// The id of the (sub)group to which this project belongs.
	assignedGroupId?: number;
	// Pretty name for this project.
	name?: string;
	// The short name of this project. Must only contain numbers, letters, dash and underline.
	slug?: string;
}