---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "reportportal_project Resource - terraform-provider-report-portal"
subcategory: ""
description: |-
  
---

# reportportal_project (Resource)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String)

### Optional

- **id** (String) The ID of this resource.

### Read-Only

- **creation_date** (Number)
- **entry_type** (String)
- **last_run** (Number)
- **launches_per_user** (List of Object) (see [below for nested schema](#nestedatt--launches_per_user))
- **launches_per_week** (String)
- **launches_quantity** (Number)
- **organization** (String)
- **unique_tickets** (Number)
- **users_quantity** (Number)

<a id="nestedatt--launches_per_user"></a>
### Nested Schema for `launches_per_user`

Read-Only:

- **count** (Number)
- **full_name** (String)


