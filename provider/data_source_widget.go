package terraform_provider_report_portal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rpClient "github.com/rmalveis/report-portal-client-go/v1/client"
	"strconv"
)

func dataSourceWidgetsByProjectAndId() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWidgetsByProjectAndIdRead,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"widget_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"share": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"widget_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"applied_filters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"applied_filter_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"applied_filter_share": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"applied_filter_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"applied_filter_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"parameters_content_fields": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"parameters_items_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"options_latest": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"options_view_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"options_timeline": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceWidgetsByProjectAndIdRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*rpClient.Client)

	var diags diag.Diagnostics
	projectName := data.Get("project_name").(string)
	widgetId := strconv.Itoa(data.Get("id").(int))

	rawWidgetData, err := c.ReadFullWidgetDataByProjectName(&projectName, &widgetId)
	if err != nil {
		return diag.FromErr(err)
	}

	widgetsData := flattenWidgetDetailsData(rawWidgetData)
	widgetDetails := append(make([]map[string]interface{}, 0), widgetsData)

	if err = data.Set("widget_details", widgetDetails); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(widgetId)

	return diags
}

func flattenWidgetDetailsData(rawWidgetsData *rpClient.FullWidgetModel) map[string]interface{} {
	if rawWidgetsData != nil {
		flattenedWidgets := make(map[string]interface{}, 0)

		flattenedWidgets["description"] = rawWidgetsData.Description
		flattenedWidgets["name"] = rawWidgetsData.Name
		flattenedWidgets["owner"] = rawWidgetsData.Owner
		flattenedWidgets["share"] = rawWidgetsData.Share
		flattenedWidgets["widget_type"] = rawWidgetsData.WidgetType

		filters := make([]map[string]interface{}, len(rawWidgetsData.AppliedFilters), len(rawWidgetsData.AppliedFilters))
		for i, f := range rawWidgetsData.AppliedFilters {
			filter := make(map[string]interface{}, 0)
			filter["applied_filter_id"] = f.Id
			filter["applied_filter_share"] = f.Share
			filter["applied_filter_name"] = f.Name
			filter["applied_filter_type"] = f.Type

			filters[i] = filter
		}
		flattenedWidgets["applied_filters"] = filters

		contentFields := []string{}
		for _, cf := range rawWidgetsData.ContentParameters.ContentFields {
			contentFields = append(contentFields, cf)
		}
		flattenedWidgets["parameters_content_fields"] = contentFields

		flattenedWidgets["parameters_items_count"] = rawWidgetsData.ContentParameters.ItemsCount
		flattenedWidgets["options_latest"] = rawWidgetsData.ContentParameters.WidgetOptions["latest"]
		flattenedWidgets["options_view_mode"] = rawWidgetsData.ContentParameters.WidgetOptions["viewMode"]
		flattenedWidgets["options_timeline"] = rawWidgetsData.ContentParameters.WidgetOptions["timeline"]

		return flattenedWidgets
	}
	return nil
}
