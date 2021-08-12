package client

import (
	"fmt"
	"net/http"
)

func (c *Client) DeleteIntegration(id *int) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/uat/settings/auth/%d", c.HostUrl, *id), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(request)
	if err != nil {
		return err
	}

	return nil
}
