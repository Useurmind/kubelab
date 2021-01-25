

// Configuration for the ui. Will be served and configured from the backend.
export interface IUIConfig {
	// The base url of the ui backend service
    uiBaseUrl?: string;
	// The base url of the projects backend service
    projectsBaseUrl?: string;
}