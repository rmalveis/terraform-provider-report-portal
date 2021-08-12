package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rpClient "github.com/rmalveis/report-portal-client-go/client"
	"strconv"
	"strings"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		DeleteContext: resourceProjectDelete,
		Schema: map[string]*schema.Schema{
			"creation_date": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"entry_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_run": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"launches_per_user": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"full_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"launches_per_week": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"launches_quantity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"unique_tickets": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"users_quantity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceProjectDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := i.(*rpClient.Client)
	projectId, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteProject(&projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags
}

func resourceProjectCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(*rpClient.Client)

	pn := data.Get("name").(string)

	project, err := client.CreateProject(&pn)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(project.Id))

	resourceProjectRead(ctx, data, i)

	return diags
}

func resourceProjectRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(*rpClient.Client)

	pn := data.Get("name").(string)

	project, err := client.GetProjectByName(&pn)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return diags
		}
		return diag.FromErr(err)
	}

	data.Set("creation_date", project.CreationDate)
	data.Set("entry_type", project.EntryType)
	data.Set("last_run", project.LastRun)

	lpu := make([]map[string]interface{}, len(project.LaunchesPerUser))
	for _, v := range project.LaunchesPerUser {
		l := make(map[string]interface{})
		l["count"] = v.Count
		l["full_name"] = v.FullName
		lpu = append(lpu, l)
	}
	data.Set("launches_per_user", lpu)

	data.Set("launches_per_week", project.LaunchesPerWeek)
	data.Set("launches_quantity", project.LaunchesQuantity)
	data.Set("organization", project.Organization)
	data.Set("unique_tickets", project.UniqueTickets)
	data.Set("users_quantity", project.UsersQuantity)

	return diags
}
