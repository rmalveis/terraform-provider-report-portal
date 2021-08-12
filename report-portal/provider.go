package report_portal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rpClient "github.com/rmalveis/report-portal-client-go"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
			"host": {

				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"reportportal_project":            resourceProject(),
			"reportportal_dashboard":          resourceDashboard(),
			"reportportal_auth_ldap_settings": resourceAuthLdapSettings(),
			"reportportal_filter":             resourceFilter(),
			"reportportal_widget":             resourceWidget(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"reportportal_projects":                  dataSourceProjects(),
			"reportportal_auth_ldap_settings":        dataSourceAuthLdapSettings(),
			"reportportal_widget_by_project_details": dataSourceWidgetsByProjectAndId(),
			"reportportal_filters":                   dataSourceFilters(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	host := d.Get("host").(string)

	c, err := rpClient.NewClient(&rpClient.ReportPortalClientConfig{
		Username: username,
		Password: password,
		Host:     host,
	}, nil)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return c, nil
}
