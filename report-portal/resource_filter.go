package report_portal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rpClient "report-portal-provider/report-portal-client-go"
	"strconv"
)

func resourceFilter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFilterCreate,
		ReadContext:   resourceFilterRead,
		UpdateContext: resourceFilterUpdate,
		DeleteContext: resourceFilterDelete,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"share": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"orders": {
				Required: true,
				Type:     schema.TypeList,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_asc": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"sorting_column": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"conditions": {
				Required: true,
				Type:     schema.TypeList,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:     schema.TypeString,
							Required: true,
						},
						"filtering_field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceFilterCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*rpClient.Client)

	projectName := data.Get("project_name").(string)

	name := data.Get("name").(string)
	filterType := data.Get("type").(string)
	description := data.Get("description").(string)
	share := data.Get("share").(bool)

	filter := rpClient.Filter{
		Share:       share,
		Name:        name,
		Conditions:  mapToConditions(data),
		Orders:      mapToOrders(data),
		Type:        filterType,
		Description: description,
	}
	result, err := c.CreateFilterByProject(projectName, &filter)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(result.Id))

	return resourceFilterRead(ctx, data, i)
}

func resourceFilterRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*rpClient.Client)

	projectName := data.Get("project_name").(string)
	filterId, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	filter, err := c.GetFilterByProjectAndId(projectName, filterId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.Set("name", filter.Name)
	data.Set("type", filter.Type)
	data.Set("description", filter.Description)
	data.Set("share", filter.Share)
	data.Set("owner", filter.Owner)
	data.Set("orders", ordersToMap(filter))
	data.Set("conditions", conditionsToMap(filter))

	data.SetId(strconv.Itoa(filter.Id))

	var diags diag.Diagnostics
	return diags
}

func resourceFilterUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*rpClient.Client)
	projectName := data.Get("project_name").(string)
	filterId, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	filter := rpClient.Filter{
		Share:       data.Get("share").(bool),
		Id:          filterId,
		Name:        data.Get("name").(string),
		Conditions:  mapToConditions(data),
		Orders:      mapToOrders(data),
		Type:        data.Get("type").(string),
		Description: data.Get("description").(string),
	}
	err = c.UpdateFilterByProjectAndId(projectName, filter)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceFilterRead(ctx, data, i)
}

func resourceFilterDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*rpClient.Client)

	projectName := data.Get("project_name").(string)

	filterId, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = c.DeleteFilterByProjectAndId(projectName, filterId)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func conditionsToMap(filter *rpClient.Filter) []map[string]interface{} {
	conditionsSlice := make([]map[string]interface{}, 0, 10)
	for _, condition := range filter.Conditions {
		conditionMap := make(map[string]interface{})
		conditionMap["filtering_field"] = condition.FilteringField
		conditionMap["condition"] = condition.Condition
		conditionMap["value"] = condition.Value
		conditionsSlice = append(conditionsSlice, conditionMap)
	}
	return conditionsSlice
}

func ordersToMap(filter *rpClient.Filter) []map[string]interface{} {
	ordersSlice := make([]map[string]interface{}, 0, 10)
	for _, order := range filter.Orders {
		orderMap := make(map[string]interface{})
		orderMap["sorting_column"] = order.SortingColumn
		orderMap["is_asc"] = order.IsAsc
		ordersSlice = append(ordersSlice, orderMap)
	}
	return ordersSlice
}

func mapToOrders(data *schema.ResourceData) []rpClient.Order {
	ordersSlice := data.Get("orders").([]interface{})
	orders := make([]rpClient.Order, 0, 10)
	for _, orderItem := range ordersSlice {
		orderMap := orderItem.(map[string]interface{})
		order := rpClient.Order{
			SortingColumn: orderMap["sorting_column"].(string),
			IsAsc:         orderMap["is_asc"].(bool),
		}
		orders = append(orders, order)
	}
	return orders
}

func mapToConditions(data *schema.ResourceData) []rpClient.Condition {
	conditionsList := data.Get("conditions").([]interface{})
	conditions := make([]rpClient.Condition, 0, 10)
	for _, conditionItem := range conditionsList {
		conditionMap := conditionItem.(map[string]interface{})
		cond := rpClient.Condition{
			FilteringField: conditionMap["filtering_field"].(string),
			Condition:      conditionMap["condition"].(string),
			Value:          conditionMap["value"].(string),
		}
		conditions = append(conditions, cond)
	}
	return conditions
}
