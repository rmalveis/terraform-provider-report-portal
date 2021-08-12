package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type FilterQuery struct {
	Id        *int    `url:"filter.eq.id,omitempty"`
	Name      *string `url:"filter.eq.name,omitempty"`
	Owner     *string `url:"filter.eq.owner,omitempty"`
	ProjectId *int    `url:"filter.eq.projectId,omitempty"`
	Shared    *bool   `url:"filter.eq.shared,omitempty"`
}

type PaginationQuery struct {
	Page *int    `url:"page,omitempty"`
	Size *int    `url:"size,omitempty"`
	Sort *string `url:"sort,omitempty"`
}

type GetFiltersByProjectResponse struct {
	Content []Filter           `json:"content"`
	Page    PaginationResponse `json:"page"`
}

type PaginationResponse struct {
	Number        int `json:"number"`
	Size          int `json:"size"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
}

type CreateFilterByProjectResponse struct {
	Id int `json:"id"`
}

type Filter struct {
	Owner       string      `json:"owner,omitempty"`
	Share       bool        `json:"share"`
	Id          int         `json:"id,omitempty"`
	Name        string      `json:"name"`
	Conditions  []Condition `json:"conditions"`
	Orders      []Order     `json:"orders"`
	Type        string      `json:"type"`
	Description string      `json:"description,omitempty"`
}

type Condition struct {
	FilteringField string `json:"filteringField"`
	Condition      string `json:"condition"`
	Value          string `json:"value"`
}

type Order struct {
	SortingColumn string `json:"sortingColumn"`
	IsAsc         bool   `json:"isAsc"`
}

func (c *Client) GetFiltersByProject(projectName string, filter *FilterQuery, pagination *PaginationQuery) (*GetFiltersByProjectResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/filter", c.HostUrl, url.PathEscape(projectName)), nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = parseQuery(filter) + parseQuery(pagination)
	req.Header.Add("Content-Type", "application/json")

	respBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var response GetFiltersByProjectResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) CreateFilterByProject(projectName string, filter *Filter) (*CreateFilterByProjectResponse, error) {
	body, err := json.Marshal(filter)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/%s/filter", c.HostUrl, url.PathEscape(projectName)), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	respBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp CreateFilterByProjectResponse
	err = json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) GetFilterByProjectAndId(projectName string, filterId int) (*Filter, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/filter/%d", c.HostUrl, url.PathEscape(projectName), filterId), nil)
	if err != nil {
		return nil, err
	}

	respBody, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var filter Filter
	err = json.Unmarshal(respBody, &filter)
	if err != nil {
		return nil, err
	}

	return &filter, nil
}

func (c *Client) UpdateFilterByProjectAndId(projectName string, filter Filter) error {
	filterId := filter.Id
	filter.Id = 0

	reqBody, err := json.Marshal(filter)

	log.Printf("Request body: %s", string(reqBody))

	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/%s/filter/%d", c.HostUrl, url.PathEscape(projectName), filterId), bytes.NewReader(reqBody))
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

func (c *Client) DeleteFilterByProjectAndId(projectName string, filterId int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/%s/filter/%d", c.HostUrl, url.PathEscape(projectName), filterId), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
