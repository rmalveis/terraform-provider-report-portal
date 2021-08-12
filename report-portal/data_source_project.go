package report_portal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rpClient "report-portal-provider/report-portal-client-go"
	"strconv"
	"time"
)

func dataSourceProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProjectsRead,
		Schema: map[string]*schema.Schema{
			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"users_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"launches_quantity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_run": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"entry_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceProjectsRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*rpClient.Client)

	var diags diag.Diagnostics

	projects, err := c.GetAllProjects()
	if err != nil {
		return diag.FromErr(err)
	}

	projectItems := flattenProjectsItemData(projects.Content)

	if err = data.Set("projects", projectItems); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenProjectsItemData(projects []rpClient.Project) []map[string]interface{} {
	if projects != nil {
		content := make([]map[string]interface{}, len(projects), len(projects))
		for i, p := range projects {
			project := make(map[string]interface{}, 0)
			project["id"] = p.Id
			project["project_name"] = p.ProjectName
			project["users_quantity"] = p.UsersQuantity
			project["launches_quantity"] = p.LaunchesQuantity
			project["last_run"] = p.LastRun
			project["creation_date"] = p.CreationDate
			project["entry_type"] = p.EntryType

			content[i] = project
		}

		return content
	}

	return make([]map[string]interface{}, 0)
}
