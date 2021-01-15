package webclient

import (
	"fmt"
	"net/http"
)

// Client provides convenience methods to call the different
// REST methods for the projects service.
type Client struct {
	BaseUrl string
}

// Groups get the group client to interact with groups.
func (c *Client) Groups() *GroupClient {
	return NewGroupClient(c.BaseUrl)
}

// ExpectResponse checks the response status code and returns an error if it does not match
func ExpectResponse(resp *http.Response, status int) error {
	if resp.StatusCode != status {
		return fmt.Errorf("Expected response code %d, but got %d", status, resp.StatusCode)
	}

	return nil
}