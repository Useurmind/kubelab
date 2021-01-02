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
