package models

import "testing"

func TestGroupGatherProjectsWorks(t *testing.T) {
	expectedProjects := 5
	group := Group{
		Projects: []*ProjectRef{
			{Id: 1},
			{Id: 5},
		},
		Subgroups: []*Group{
			{
				Projects: []*ProjectRef{
					{Id: 2},
				},
			},
			{
				Subgroups: []*Group{
					{
						Projects: []*ProjectRef{
							{Id: 3},
							{Id: 4},
						},
					},
				},
			},
		},
	}

	projects := group.GatherProjects()

	if len(projects) != expectedProjects {
		t.Fatalf("The number of gathered projects should be %d but was %d", expectedProjects, len(projects))
	}
}

func TestGroupGatherSubgroupsWorks(t *testing.T) {
	expectedSubgroups := 6
	group := Group{
		Subgroups: []*Group{
			{
				Subgroups: []*Group{
					{Id: 2},
				},
			},
			{
				Subgroups: []*Group{
					{
						Subgroups: []*Group{
							{Id: 3},
							{Id: 4},
						},
					},
				},
			},
		},
	}

	subgroups := group.GatherSubgroups()

	if len(subgroups) != expectedSubgroups {
		t.Fatalf("The number of gathered projects should be %d but was %d", expectedSubgroups, len(subgroups))
	}
}

func TestInsertProjectRef(t *testing.T) {
	group := Group{
		Subgroups: []*Group{
			{
				Subgroups: []*Group{
					{Id: 2},
				},
			},
			{
				Subgroups: []*Group{
					{
						Subgroups: []*Group{
							{Id: 3},
							{Id: 4},
						},
					},
				},
			},
		},
	}

	err := group.InsertProjectRef(&Project{Id: 1, AssignedGroupId: 3})
	if err != nil {
		t.Fatal(err)
	}

	if group.Subgroups[1].Subgroups[0].Subgroups[0].Projects[0].Id != 1 {
		t.Error("Project was not correctly assigned")
	}
}


func TestUpdateProjectRef(t *testing.T) {
	group := Group{
		Subgroups: []*Group{
			{
				Subgroups: []*Group{
					{Id: 2},
				},
			},
			{
				Subgroups: []*Group{
					{
						Subgroups: []*Group{
							{Id: 3, Projects: []*ProjectRef{ {Id: 1, Name: "old"} }},
							{Id: 4},
						},
					},
				},
			},
		},
	}

	err := group.UpdateProjectRef(&Project{Id: 1, AssignedGroupId: 3, Name: "new"})
	if err != nil {
		t.Fatal(err)
	}

	if group.Subgroups[1].Subgroups[0].Subgroups[0].Projects[0].Name != "new" {
		t.Error("Project was not correctly updated")
	}
}

func TestDeleteProjectRef(t *testing.T) {
	group := Group{
		Subgroups: []*Group{
			{
				Subgroups: []*Group{
					{Id: 2},
				},
			},
			{
				Subgroups: []*Group{
					{
						Subgroups: []*Group{
							{Id: 3, Projects: []*ProjectRef{ {Id: 1, Name: "old"} }},
							{Id: 4},
						},
					},
				},
			},
		},
	}

	err := group.DeleteProjectRef(&Project{Id: 1, AssignedGroupId: 3})
	if err != nil {
		t.Fatal(err)
	}

	if len(group.Subgroups[1].Subgroups[0].Subgroups[0].Projects) > 0 {
		t.Error("Project was not correctly deleted")
	}
}


