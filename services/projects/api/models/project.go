package main

type Group struct {
	Id int `json:"id"`
	Name string `json:"name"`
	SubgroupNames []string `json:"subgroupNames"`
	Projects []Project `json:"projects"`
}

type Project struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}