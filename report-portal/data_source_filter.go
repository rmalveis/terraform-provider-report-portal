package report_portal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rpClient "report-portal-provider/report-portal-client-go"
	"strconv"
	"time"
)

func dataSourceFilters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFiltersRead,
		Schema: map[string]*schema.Schema{
			"project_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"share": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"conditions": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"filtering_field": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"condition": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"orders": {
							Computed: true,
							Type:     schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sorting_column": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"is_asc": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceFiltersRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*rpClient.Client)

	var diags diag.Diagnostics

	projectName := data.Get("project_name").(string)

	filtersMap, diagnostics, done := getAllProjectFilter(c, projectName)
	if done {
		return diagnostics
	}

	if err := data.Set("filters", filtersMap); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func getAllProjectFilter(c *rpClient.Client, projectName string) ([]map[string]interface{}, diag.Diagnostics, bool) {
	filterSlice := make([]map[string]interface{}, 0, 10)
	currentPage := 1
	defaultSize := 100
	sort := "asc"
	for {
		pagination := rpClient.PaginationQuery{
			Page: &currentPage,
			Size: &defaultSize,
			Sort: &sort,
		}

		filters, err := c.GetFiltersByProject(projectName, nil, &pagination)
		if err != nil {
			return nil, diag.FromErr(err), true
		}

		for _, filter := range filters.Content {
			filterMap := mapFilter(filter)
			filterSlice = append(filterSlice, filterMap)
		}

		if filters.Page.TotalPages >= currentPage {
			break
		}
		currentPage++
	}
	return filterSlice, nil, false
}

func mapFilter(filter rpClient.Filter) map[string]interface{} {
	filterMap := make(map[string]interface{})
	filterMap["id"] = filter.Id
	filterMap["owner"] = filter.Owner
	filterMap["share"] = filter.Share
	filterMap["name"] = filter.Name
	filterMap["type"] = filter.Type

	conditionList := make([]map[string]interface{}, 0, 10)
	for _, condition := range filter.Conditions {
		conditionMap := make(map[string]interface{})
		conditionMap["filtering_field"] = condition.FilteringField
		conditionMap["condition"] = condition.Condition
		conditionMap["value"] = condition.Value
		conditionList = append(conditionList, conditionMap)
	}
	filterMap["conditions"] = conditionList

	orderList := make([]map[string]interface{}, 0, 10)
	for _, order := range filter.Orders {
		orderMap := make(map[string]interface{})
		orderMap["sorting_column"] = order.SortingColumn
		orderMap["is_asc"] = order.IsAsc
		orderList = append(orderList, orderMap)
	}
	filterMap["orders"] = orderList
	return filterMap
}
