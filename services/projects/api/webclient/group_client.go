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

func (c *GroupClient) List(startIndex int64, count int64) ([]*models.Group, error) {
	resp, err := http.Get(fmt.Sprintf("%s?start=%d&count=%d", c.GroupBaseUrl, startIndex, count))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	err = ExpectResponse(resp, 200)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)

	groups := make([]*models.Group, 0)
	err = dec.Decode(&groups)
	if err != nil {
		return nil, err
	}

	return groups, err
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

func (c *GroupClient) Delete(groupID int64) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%d", c.GroupBaseUrl, groupID), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = ExpectResponse(resp, 204)
	if err != nil {
		return err
	}

	return nil
}