package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type LdapIntegrationParameters struct {
	Enabled             bool   `json:"-"`
	PasswordEncoderType string `json:"passwordEncoderType"`
	Url                 string `json:"url"`
	BaseDn              string `json:"baseDn"`
	Email               string `json:"email"`
	UserDnPattern       string `json:"userDnPattern"`
	UserSearchFilter    string `json:"userSearchFilter"`
	GroupSearchBase     string `json:"groupSearchBase"`
	GroupSearchFilter   string `json:"groupSearchFilter"`
	PasswordAttribute   string `json:"passwordAttribute"`
	FullName            string `json:"fullName"`
	Photo               string `json:"photo"`
	ManagerDn           string `json:"managerDn"`
	ManagerPassword     string `json:"managerPassword"`
}

type LdapIntegrationPayload struct {
	Enabled               bool                       `json:"enabled"`
	IntegrationParameters *LdapIntegrationParameters `json:"integrationParameters"`
	Name                  string                     `json:"name"`
}

type LdapSettings struct {
	Id             *int `json:"id"`
	LdapAttributes struct {
		Enabled                   *bool   `json:"enabled"`
		Url                       *string `json:"url"`
		BaseDn                    *string `json:"baseDn"`
		SynchronizationAttributes struct {
			Email    *string `json:"email"`
			FullName *string `json:"fullName"`
			Photo    *string `json:"photo"`
		} `json:"synchronizationAttributes"`
	} `json:"ldapAttributes"`
	UserDnPattern       *string `json:"userDnPattern"`
	UserSearchFilter    *string `json:"userSearchFilter"`
	GroupSearchBase     *string `json:"groupSearchBase"`
	GroupSearchFilter   *string `json:"groupSearchFilter"`
	PasswordEncoderType *string `json:"passwordEncoderType"`
	PasswordAttribute   *string `json:"passwordAttribute"`
	ManagerDn           *string `json:"managerDn"`
	ManagerPassword     *string `json:"managerPassword"`
}

func (l *LdapSettings) String() {
	fmt.Printf(
		"{\n"+
			"\"id\": %d, \n"+
			"\"ldapAttributes\": { \n"+
			"\"enabled\": %t, \n"+
			"\"url\": \"%s\", \n"+
			"\"baseDn\": \"%s\", \n"+
			"\"synchronizationAttributes\": { \n"+
			"\"email\": \"%s\", \n"+
			"\"fullName\": \"%s\", \n"+
			"\"photo\": \"%s\" \n"+
			"}\n"+
			"}\n"+
			"\"userDnPattern\": \"%s\", \n"+
			"\"userSearchFilter\": \"%s\", \n"+
			"\"groupSearchBase\": \"%s\", \n"+
			"\"groupSearchFilter\": \"%s\" \n"+
			"\"passwordEncoderType\": \"%s\",\n"+
			"\"passwordAttribute\": \"%s\"\n"+
			"}",
		*l.Id,
		*l.LdapAttributes.Enabled,
		*l.LdapAttributes.Url,
		*l.LdapAttributes.BaseDn,
		*l.LdapAttributes.SynchronizationAttributes.Email,
		*l.LdapAttributes.SynchronizationAttributes.FullName,
		*l.LdapAttributes.SynchronizationAttributes.Photo,
		*l.UserDnPattern,
		*l.UserSearchFilter,
		*l.GroupSearchBase,
		*l.GroupSearchFilter,
		*l.PasswordEncoderType,
		*l.PasswordAttribute,
	)
}

func (c *Client) CreateAuthLdapSettings(config *LdapIntegrationParameters) (*LdapSettings, error) {
	var payload LdapIntegrationPayload
	var err error
	payload.Enabled = config.Enabled
	payload.IntegrationParameters = config

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(data))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/uat/settings/auth/ldap", c.HostUrl), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))

	var resp LdapSettings
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)

	return &resp, nil
}

func (c *Client) ReadLdapAuthSettings() (*LdapSettings, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/uat/settings/auth/ldap", c.HostUrl), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp LdapSettings
	fmt.Println(string(body))
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) UpdateAuthLdapSettings(config *LdapIntegrationParameters) (*LdapSettings, error) {
	var payload LdapIntegrationPayload
	var err error
	payload.Enabled = config.Enabled
	payload.IntegrationParameters = config

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(data))
	log.Println("[DEBUG] {SETTINGS:UpdateAuthLdapSettings} Request Data: " + string(data))
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/uat/settings/auth/ldap", c.HostUrl), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(body))
	log.Println("[DEBUG] {SETTINGS:UpdateAuthLdapSettings} Response Data: " + string(body))

	var resp LdapSettings
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
