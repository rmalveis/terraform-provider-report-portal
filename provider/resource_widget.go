package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rpClient "github.com/rmalveis/report-portal-client-go/client"
	"strconv"
	"strings"
	"time"
)

func resourceWidget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWidgetCreate,
		ReadContext:   resourceWidgetRead,
		UpdateContext: resourceWidgetUpdate,
		DeleteContext: resourceWidgetDelete,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"share": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"widget_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"widget_type_calculated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"parameters_content_fields": { // TODO: Find a better name: criteria?
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parameters_content_fields_calculated": { // Workaround
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parameters_items_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"options_latest": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"options_timeline": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"options_view_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"options_zoom": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"options_action_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"options_user": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"options_launch_name_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"options_include_methods": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceWidgetDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	// TODO: There is no way to delete the widget without a dashboard
	var diags diag.Diagnostics

	data.SetId("")

	return diags
}

func resourceWidgetRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := i.(*rpClient.Client)
	pn := data.Get("project_name").(string)
	widgetId := data.Id()

	widgetSettings, err := client.ReadFullWidgetDataByProjectName(&pn, &widgetId)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return diags
		}
		return diag.FromErr(err)
	}

	err = data.Set("project_name", pn)
	err = data.Set("description", widgetSettings.Description)
	err = data.Set("name", widgetSettings.Name)
	err = data.Set("owner", widgetSettings.Owner)
	err = data.Set("share", widgetSettings.Share)
	err = data.Set("widget_type_calculated", widgetSettings.WidgetType)
	err = data.Set("filter_ids", getFilterIds(widgetSettings.AppliedFilters))
	err = data.Set("parameters_content_fields_calculated", widgetSettings.ContentParameters.ContentFields)
	err = data.Set("parameters_items_count", widgetSettings.ContentParameters.ItemsCount)
	err = data.Set("options_latest", widgetSettings.ContentParameters.WidgetOptions["latest"])
	err = data.Set("options_view_mode", widgetSettings.ContentParameters.WidgetOptions["viewMode"])
	err = data.Set("options_timeline", widgetSettings.ContentParameters.WidgetOptions["timeline"])
	err = data.Set("options_zoom", widgetSettings.ContentParameters.WidgetOptions["zoom"])
	err = data.Set("options_action_type", widgetSettings.ContentParameters.WidgetOptions["actionTypes"])
	err = data.Set("options_user", widgetSettings.ContentParameters.WidgetOptions["user"])
	err = data.Set("options_launch_name_filter", widgetSettings.ContentParameters.WidgetOptions["launchNameFilter"])
	err = data.Set("options_include_methods", widgetSettings.ContentParameters.WidgetOptions["includeMethods"])

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceWidgetCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(*rpClient.Client)

	widgetSettings, err := getWidgetParameters(data)
	if err != nil {
		return diag.FromErr(err)
	}

	pn := data.Get("project_name").(string)

	savedSettings, err := client.CreateWidgetByProject(&pn, widgetSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(savedSettings.Id))
	resourceWidgetRead(ctx, data, i)

	return diags
}

func resourceWidgetUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(*rpClient.Client)

	widgetParameters, err := getWidgetParameters(data)
	if err != nil {
		return diag.FromErr(err)
	}
	pn := data.Get("project_name").(string)
	wi := data.Id()

	err = client.UpdateWidgetByProject(&pn, &wi, widgetParameters)
	if err != nil {
		return diag.FromErr(err)
	}

	err = data.Set("latest_updated", time.Now().Format(time.RFC3339))
	if err != nil {
		return nil
	}
	resourceWidgetRead(ctx, data, i)

	return diags
}

// Auxiliary functions
func encodeWidgetTypeOption(widgetType string) string {
	return rpClient.WidgetTypes[widgetType]
}

func getCriteriaByWidgetType(widgetType *string, data *schema.ResourceData) ([]string, error) {
	switch *widgetType {
	case "launchesDurationChart":
		return []string{"startTime", "endTime", "name", "number", "status"}, nil
	case "bugTrend":
		return []string{"statistics$defects$product_bug$total",
			"statistics$defects$automation_bug$total",
			"statistics$defects$system_issue$total",
			"statistics$defects$no_defect$total",
			"statistics$defects$to_investigate$total"}, nil
	case "topTestCases":
		r := data.Get("parameters_content_fields").([]interface{})
		if len(r) == 0 || len(r) > 1 {
			return nil, fmt.Errorf("This Widget Type must have one and only one criteria.(parameters_content_fields)")
		}
		return getCriteriaValues(r), nil
	case "flakyTestCases":
		return nil, nil
	default:
		r := data.Get("parameters_content_fields").([]interface{})
		if len(r) == 0 {
			return nil, fmt.Errorf("A criteria must be provided (parameters_content_fields)")
		}
		return getCriteriaValues(r), nil
	}
}

func getCriteriaValues(criteria []interface{}) []string {
	r := make([]string, len(criteria), len(criteria))
	for i, c := range criteria {
		r[i] = rpClient.WidgetCriteria[c.(string)]
	}
	return r
}

func getFilterIds(filters []rpClient.Filter) []int {
	r := make([]int, len(filters), len(filters))
	for i, f := range filters {
		r[i] = f.Id
	}
	return r
}

func getWidgetParameters(data *schema.ResourceData) (*rpClient.WidgetInputPayload, error) {
	widgetType := encodeWidgetTypeOption(data.Get("widget_type").(string))
	contentFields, err := getCriteriaByWidgetType(&widgetType, data)
	if err != nil {
		return nil, err
	}

	var widgetSettings rpClient.WidgetInputPayload
	widgetSettings.WidgetType = widgetType
	widgetSettings.Name = data.Get("name").(string)
	widgetSettings.Description = data.Get("description").(string)
	widgetSettings.ContentParameters.ContentFields = contentFields
	widgetSettings.ContentParameters.WidgetOptions = getWidgetOptions(data)
	widgetSettings.ContentParameters.ItemsCount = data.Get("parameters_items_count").(int)
	widgetSettings.Share = data.Get("share").(bool)
	widgetSettings.FilterIds = data.Get("filter_ids").([]interface{})
	return &widgetSettings, nil
}

func getWidgetOptions(data *schema.ResourceData) map[string]interface{} {
	options := make(map[string]interface{}, 0)
	options["latest"] = data.Get("options_latest").(bool)
	options["timeline"] = data.Get("options_timeline").(string)
	options["viewMode"] = data.Get("options_view_mode").(string)
	options["zoom"] = data.Get("options_zoom").(bool)
	options["actionType"] = data.Get("options_action_type").(string)
	options["user"] = data.Get("options_user").([]interface{})
	options["launchNameFilter"] = data.Get("options_launch_name_filter").(string)
	options["includeMethods"] = data.Get("options_include_methods").(bool)

	return options
}
