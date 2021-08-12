package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type WidgetInputPayload struct {
	ContentParameters WidgetContentParameters `json:"contentParameters"`
	Description       string                  `json:"description"`
	FilterIds         []interface{}           `json:"filterIds"`
	Name              string                  `json:"name"`
	Share             bool                    `json:"share"`
	WidgetType        string                  `json:"widgetType"`
}

type WidgetContentParameters struct {
	ContentFields []string               `json:"contentFields"`
	ItemsCount    int                    `json:"itemsCount"`
	WidgetOptions map[string]interface{} `json:"widgetOptions"`

	// TODO: Later
	//WidgetOptions struct {
	//	Latest     bool     `json:"latest"`
	//	Zoom       bool     `json:"zoom"`
	//	Timeline   string   `json:"timeline"`
	//	ViewMode   string   `json:"viewMode"`
	//	ActionType string   `json:"actionType"`
	//	User       []string `json:"user"`
	//	LaunchNameFilter string `json:"launchNameFilter"`
	//	IncludeMethods bool `json:"includeMethods"`
	//
	//} `json:"widgetOptions"`
}

type FullWidgetModel struct {
	AppliedFilters []Filter `json:"appliedFilters"`
	// TODO: Content may be never be used in TF, but could be useful for a generic RP client
	//Content           struct{}                `json:"content"`
	ContentParameters WidgetContentParameters `json:"contentParameters"`
	Description       string                  `json:"description"`
	Id                int                     `json:"id"`
	Name              string                  `json:"name"`
	Owner             string                  `json:"owner"`
	Share             bool                    `json:"share"`
	WidgetType        string                  `json:"widgetType"`
}

type WidgetCreationResponseModel struct {
	Id int `json:"id"`
}

func (c *Client) ReadFullWidgetDataByProjectName(projectName *string, widgetId *string) (*FullWidgetModel, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/widget/%s", c.HostUrl, *projectName, *widgetId), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp FullWidgetModel

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateWidgetByProject(projectName *string, widgetParameters *WidgetInputPayload) (*WidgetCreationResponseModel, error) {
	var err error
	data, err := json.Marshal(widgetParameters)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/%s/widget", c.HostUrl, *projectName), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp WidgetCreationResponseModel
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) UpdateWidgetByProject(projectName *string, widgetId *string, widgetParameters *WidgetInputPayload) error {
	var err error

	data, err := json.Marshal(widgetParameters)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/%s/widget/%s", c.HostUrl, *projectName, *widgetId), bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return err
	}

	var resp WidgetCreationResponseModel
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return err
	}

	return nil
}
