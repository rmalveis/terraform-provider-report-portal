---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "reportportal_filters Data Source - terraform-provider-report-portal"
subcategory: ""
description: |-
  
---

# reportportal_filters (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **project_name** (String)

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **filters** (List of Object) (see [below for nested schema](#nestedatt--filters))

<a id="nestedatt--filters"></a>
### Nested Schema for `filters`

Read-Only:

- **conditions** (List of Object) (see [below for nested schema](#nestedobjatt--filters--conditions))
- **id** (Number)
- **name** (String)
- **orders** (List of Object) (see [below for nested schema](#nestedobjatt--filters--orders))
- **owner** (String)
- **share** (Boolean)
- **type** (String)

<a id="nestedobjatt--filters--conditions"></a>
### Nested Schema for `filters.conditions`

Read-Only:

- **condition** (String)
- **filtering_field** (String)
- **value** (String)


<a id="nestedobjatt--filters--orders"></a>
### Nested Schema for `filters.orders`

Read-Only:

- **is_asc** (Boolean)
- **sorting_column** (String)


