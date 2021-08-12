package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Project struct {
	Id               int    `json:"id"`
	ProjectName      string `json:"projectName"`
	UsersQuantity    int    `json:"usersQuantity"`
	LaunchesQuantity int    `json:"launchesQuantity"`
	LastRun          int    `json:"lastRun"`
	CreationDate     int    `json:"creationDate"`
	EntryType        string `json:"entryType"`
}

type GetAllProjectsResponse struct {
	Content []Project `json:"content"`
}

type CreateProjectRequest struct {
	EntryType   string `json:"entryType"`
	ProjectName string `json:"projectName"`
}

type CreateProjectResponse struct {
	Id int `json:"id"`
}

type LaunchPerUser struct {
	Count    int    `json:"count"`
	FullName string `json:"fullName"`
}

type GetProjectByNameResponse struct {
	Project
	LaunchesPerUser  []LaunchPerUser `json:"launchesPerUser"`
	LaunchesPerWeek  string          `json:"launchesPerWeek"`
	LaunchesQuantity int             `json:"launchesQuantity"`
	Organization     string          `json:"organization"`
	UniqueTickets    int             `json:"uniqueTickets"`
	UsersQuantity    int             `json:"usersQuantity"`
}

func (c *Client) GetAllProjects() (*GetAllProjectsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/project/list", c.HostUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp GetAllProjectsResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) CreateProject(projectName *string) (*CreateProjectResponse, error) {
	request := CreateProjectRequest{
		EntryType:   "INTERNAL",
		ProjectName: *projectName,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/project", c.HostUrl), bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	respBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var jsonBody CreateProjectResponse
	err = json.Unmarshal(respBody, &jsonBody)
	if err != nil {
		return nil, err
	}

	return &jsonBody, nil
}

func (c *Client) GetProjectByName(projectName *string) (*GetProjectByNameResponse, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/project/list/%s", c.HostUrl, url.PathEscape(*projectName)), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(request)
	if err != nil {
		return nil, err
	}

	var response GetProjectByNameResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) DeleteProject(projectId *int) error {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/project/%d", c.HostUrl, *projectId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(request)
	if err != nil {
		return err
	}

	return nil
}
