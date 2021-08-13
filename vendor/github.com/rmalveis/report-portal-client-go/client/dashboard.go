package client

//TODO: Need to refactor these methods. Standardize the ap exported objects.
import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type CreateDashboardRequest struct {
	ProjectName string `json:"-"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Share       bool   `json:"share"`
}

type CreateDashboardResponse struct {
	Id int `json:"id"`
}

type GetDashboardByIdResponse struct {
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Owner       string   `json:"owner"`
	Share       bool     `json:"share"`
	Widgets     []Widget `json:"widgets"`
}

type UpdateDashboardRequest struct {
	DashboardId   int      `json:"-"`
	UpdateWidgets []Widget `json:"updateWidgets"`
	CreateDashboardRequest
}

type Widget struct {
	Share          bool   `json:"share"`
	WidgetId       int    `json:"widgetId"`
	WidgetName     string `json:"widgetName"`
	WidgetPosition struct {
		PositionX int `json:"positionX"`
		PositionY int `json:"positionY"`
	} `json:"widgetPosition"`
	WidgetSize struct {
		Height int `json:"height"`
		Width  int `json:"width"`
	} `json:"widgetSize"`
	WidgetType string `json:"widgetType"`
}

type AddWidgetRequest struct {
	AddWidget Widget `json:"addWidget"`
}

func (c *Client) CreateDashboard(d CreateDashboardRequest) (*int, error) {
	reqBody := CreateDashboardRequest{
		Description: d.Description,
		Name:        d.Name,
		Share:       d.Share,
	}

	reqBodyAsJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/v1/%s/dashboard", c.HostUrl, url.PathEscape(d.ProjectName)),
		bytes.NewReader(reqBodyAsJson))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var cdr CreateDashboardResponse
	err = json.Unmarshal(resBody, &cdr)
	if err != nil {
		return nil, err
	}

	return &cdr.Id, nil
}

func (c *Client) UpdateDashboard(updateDashboardRequest *UpdateDashboardRequest) error {
	reqBody, err := json.Marshal(*updateDashboardRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/api/v1/%s/dashboard/%d", c.HostUrl, url.PathEscape(updateDashboardRequest.ProjectName), updateDashboardRequest.DashboardId),
		bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	_, err = c.doRequest(req)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetDashboardById(projectName string, dashboardId *int) (*GetDashboardByIdResponse, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/api/v1/%s/dashboard/%d", c.HostUrl, url.PathEscape(projectName), *dashboardId),
		nil)
	if err != nil {
		return nil, err
	}

	resBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response GetDashboardByIdResponse
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) DeleteDashboardById(projectName *string, dashboardId *int) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%s/api/v1/%s/dashboard/%d", c.HostUrl, url.PathEscape(*projectName), *dashboardId),
		nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AddWidgetIntoDashboard(projectName string, dashboardId *int, widget *Widget) error {
	reqBody, err := json.Marshal(AddWidgetRequest{
		AddWidget: *widget,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/api/v1/%s/dashboard/%d/add", c.HostUrl, url.PathEscape(projectName), *dashboardId),
		bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	_, err = c.doRequest(req)

	if err != nil {
		return err
	}

	return nil
}
