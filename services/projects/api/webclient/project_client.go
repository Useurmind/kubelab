package webclient

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"

	"github.com/useurmind/kubelab/services/projects/api/models"
)

// ProjectClient is the web api client to interact with projects
type ProjectClient struct {
	ProjectBaseUrl string
}

func NewProjectClient(baseUrl string) *ProjectClient {
	return &ProjectClient{
		ProjectBaseUrl: fmt.Sprintf("%s/projects", baseUrl),
	}
}

func (c *ProjectClient) Get(projectID int64) (*models.Project, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d", c.ProjectBaseUrl, projectID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	err = ExpectResponse(resp, 200, resp.Body)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)

	project := models.Project{}
	err = dec.Decode(&project)
	if err != nil {
		return nil, err
	}

	return &project, err
}

// func (c *ProjectClient) List(startIndex int64, count int64) ([]*models.Project, error) {
// 	resp, err := http.Get(fmt.Sprintf("%s?start=%d&count=%d", c.ProjectBaseUrl, startIndex, count))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
	
// 	err = ExpectResponse(resp, 200)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dec := json.NewDecoder(resp.Body)

// 	projects := make([]*models.Project, 0)
// 	err = dec.Decode(&projects)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return projects, err
// }

func (c *ProjectClient) Create(project *models.Project) (*models.Project, error) {
	data, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)

	resp, err := http.Post(c.ProjectBaseUrl, "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = ExpectResponse(resp, 200, resp.Body)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)

	createdproject := models.Project{}
	err = dec.Decode(&createdproject)
	if err != nil {
		return nil, err
	}

	return &createdproject, err
}

func (c *ProjectClient) Update(project *models.Project) (*models.Project, error) {
	data, err := json.Marshal(project)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)

	req, err := http.NewRequest(http.MethodPut, c.ProjectBaseUrl, reader)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = ExpectResponse(resp, 200, resp.Body)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)

	createdproject := models.Project{}
	err = dec.Decode(&createdproject)
	if err != nil {
		return nil, err
	}

	return &createdproject, err
}

func (c *ProjectClient) Delete(projectID int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", c.ProjectBaseUrl, projectID), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = ExpectResponse(resp, 204, resp.Body)
	if err != nil {
		return err
	}

	return nil
}