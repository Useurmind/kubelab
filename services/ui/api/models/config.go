package models

// Configuration for the ui. Will be served and configured from the backend.
type UIConfig struct {
	// The base url of the ui backend service
	UiBaseUrl string `json:"uiBaseUrl"`
	// The base url of the projects backend service
	ProjectsBaseUrl string `json:"projectsBaseUrl"`
}