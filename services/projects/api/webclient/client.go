package webclient

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// ErrorNotFound describes the case when the entity could not be found via a REST call.
var ErrorNotFound = errors.New("404: Not found")

// Client provides convenience methods to call the different
// REST methods for the projects service.
type Client struct {
	BaseUrl string
}

// Groups get the group client to interact with groups.
func (c *Client) Groups() *GroupClient {
	return NewGroupClient(c.BaseUrl)
}

// Projects get the project client to interact with projects.
func (c *Client) Projects() *ProjectClient {
	return NewProjectClient(c.BaseUrl)
}

// ExpectResponse checks the response status code and returns an error if it does not match
func ExpectResponse(resp *http.Response, status int, body io.ReadCloser) error {
	if resp.StatusCode == 404 {
		return ErrorNotFound
	}

	if resp.StatusCode != status {
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}

		return fmt.Errorf("Expected response code %d, but got %d: %s", status, resp.StatusCode, string(bodyBytes))
	}

	return nil
}