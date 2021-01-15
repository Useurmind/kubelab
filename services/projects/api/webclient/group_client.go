package webclient

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"

	"github.com/useurmind/kubelab/services/projects/api/models"
)

// GroupClient is the web api client to interact with groups
type GroupClient struct {
	GroupBaseUrl string
}

func NewGroupClient(baseUrl string) *GroupClient {
	return &GroupClient{
		GroupBaseUrl: fmt.Sprintf("%s/groups", baseUrl),
	}
}

func (c *GroupClient) Get(groupId int64) (*models.Group, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%d", c.GroupBaseUrl, groupId))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	err = ExpectResponse(resp, 200)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)

	group := models.Group{}
	err = dec.Decode(&group)
	if err != nil {
		return nil, err
	}

	return &group, err
}

func (c *GroupClient) Create(group *models.Group) (*models.Group, error) {
	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(data)

	resp, err := http.Post(c.GroupBaseUrl, "application/json", reader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = ExpectResponse(resp, 200)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)

	createdGroup := models.Group{}
	err = dec.Decode(&createdGroup)
	if err != nil {
		return nil, err
	}

	return &createdGroup, err
}