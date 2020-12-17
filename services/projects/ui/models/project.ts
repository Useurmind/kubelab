

export interface IGroup {
	id: number;
	name: string;
	subgroupNames: string[];
	projects: IProject[];
}

export interface IProject {
	id: number;
	name: string;
	slug: string;
}