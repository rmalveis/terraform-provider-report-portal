package terraform_provider_report_portal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rpClient "github.com/rmalveis/report-portal-client-go/v1/client"
	"strconv"
)

func resourceDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDashboardCreate,
		ReadContext:   resourceDashboardRead,
		UpdateContext: resourceDashboardUpdate,
		DeleteContext: resourceDashboardDelete,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"share": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"widgets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"widget_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"widget_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"share": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"position_x": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"position_y": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"height": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"width": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceDashboardCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	description := data.Get("description").(string)
	name := data.Get("name").(string)
	share := data.Get("share").(bool)
	projectName := data.Get("project_name").(string)

	c := i.(*rpClient.Client)

	dashboardId, err := c.CreateDashboard(rpClient.CreateDashboardRequest{
		ProjectName: projectName,
		Description: description,
		Name:        name,
		Share:       share,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(strconv.Itoa(*dashboardId))

	widgetSlice := data.Get("widgets").([]interface{})
	for _, item := range widgetSlice {
		widget := mapToWidget(item, dashboardId)
		err = c.AddWidgetIntoDashboard(projectName, dashboardId, widget)
		if err != nil {
			diag.FromErr(err)
		}

		return resourceDashboardRead(ctx, data, i)
	}

	return resourceDashboardRead(ctx, data, i)
}

func resourceDashboardRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	projectName := data.Get("project_name").(string)
	dashboardId, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := i.(*rpClient.Client)
	dashboard, err := client.GetDashboardById(projectName, &dashboardId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("description", dashboard.Description)
	data.Set("name", dashboard.Name)
	data.Set("share", dashboard.Share)
	data.Set("project_name", projectName)

	widgetSlice := make([]map[string]interface{}, 0, 10)
	for _, widget := range dashboard.Widgets {
		widgetMap := widgetToMap(widget)
		widgetSlice = append(widgetSlice, widgetMap)
	}
	data.Set("widgets", widgetSlice)

	return diags
}

func resourceDashboardUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	//projectName := data.Get("project_name").(string)
	//dashboardId, err := strconv.Atoi(data.Id())
	//if err != nil {
	//	return diag.FromErr(err)
	//}

	//c := i.(*rpClient.Client)
	//dashboard := &rpClient.UpdateDashboardRequest{
	//	CreateDashboardRequest: rpClient.CreateDashboardRequest{
	//		ProjectName: projectName,
	//		Description: "",
	//		Name:        "",
	//		Share:       false,
	//	},
	//	DashboardId:            0,
	//	UpdateWidgets:          nil,
	//}
	//c.UpdateDashboard()

	return resourceDashboardRead(ctx, data, i)
}

func resourceDashboardDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(*rpClient.Client)

	id, err := strconv.Atoi(data.Id())
	projectName := data.Get("project_name").(string)

	err = client.DeleteDashboardById(&projectName, &id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func mapToWidget(item interface{}, dashboardId *int) *rpClient.Widget {
	widgetMap := item.(map[string]interface{})
	widget := &rpClient.Widget{
		Share:      widgetMap["share"].(bool),
		WidgetId:   *dashboardId,
		WidgetName: widgetMap["name"].(string),
		WidgetPosition: struct {
			PositionX int `json:"positionX"`
			PositionY int `json:"positionY"`
		}{
			PositionX: widgetMap["position_x"].(int),
			PositionY: widgetMap["position_y"].(int),
		},
		WidgetSize: struct {
			Height int `json:"height"`
			Width  int `json:"width"`
		}{
			Height: widgetMap["height"].(int),
			Width:  widgetMap["width"].(int),
		},
		WidgetType: widgetMap["type"].(string),
	}
	return widget
}

func widgetToMap(widget rpClient.Widget) map[string]interface{} {
	widgetMap := make(map[string]interface{})
	widgetMap["widget_id"] = widget.WidgetId
	widgetMap["widget_name"] = widget.WidgetName
	widgetMap["share"] = widget.Share
	widgetMap["position_x"] = widget.WidgetPosition.PositionX
	widgetMap["position_y"] = widget.WidgetPosition.PositionY
	widgetMap["height"] = widget.WidgetSize.Height
	widgetMap["width"] = widget.WidgetSize.Width
	return widgetMap
}
