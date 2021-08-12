package report_portal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	rpClient "report-portal-provider/report-portal-client-go"
	"strconv"
	"time"
)

func dataSourceAuthLdapSettings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLdapSettingsRead,
		Schema: map[string]*schema.Schema{
			"auth_ldap_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ldap_attrs_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ldap_attrs_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ldap_attrs_base_dn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ldap_attrs_sync_email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ldap_attrs_sync_fullname": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ldap_attrs_sync_photo": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_dn_pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_search_filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_search_base": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_search_filter": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password_encoder_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(rpClient.PasswordEncryptionTypes, true),
						},
						"password_attr": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLdapSettingsRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*rpClient.Client)

	var diags diag.Diagnostics

	rawLdapSettings, err := c.ReadLdapAuthSettings()
	if err != nil {
		return diag.FromErr(err)
	}

	ldapSettings := flattenLdapSettingsItemData(rawLdapSettings)
	state := append(make([]map[string]interface{}, 0), ldapSettings)

	if err = data.Set("auth_ldap_settings", state); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func flattenLdapSettingsItemData(settings *rpClient.LdapSettings) map[string]interface{} {
	if settings != nil {
		flattenedSettings := make(map[string]interface{}, 0)
		flattenedSettings["id"] = settings.Id
		flattenedSettings["ldap_attrs_enabled"] = settings.LdapAttributes.Enabled
		flattenedSettings["ldap_attrs_url"] = settings.LdapAttributes.Url
		flattenedSettings["ldap_attrs_base_dn"] = settings.LdapAttributes.BaseDn
		flattenedSettings["ldap_attrs_sync_email"] = settings.LdapAttributes.SynchronizationAttributes.Email
		flattenedSettings["ldap_attrs_sync_fullname"] = settings.LdapAttributes.SynchronizationAttributes.FullName
		flattenedSettings["ldap_attrs_sync_photo"] = settings.LdapAttributes.SynchronizationAttributes.Photo
		flattenedSettings["user_dn_pattern"] = settings.UserDnPattern
		flattenedSettings["user_search_filter"] = settings.UserSearchFilter
		flattenedSettings["group_search_base"] = settings.GroupSearchBase
		flattenedSettings["group_search_filter"] = settings.GroupSearchFilter
		flattenedSettings["password_encoder_type"] = settings.PasswordEncoderType
		flattenedSettings["password_attr"] = settings.PasswordAttribute

		return flattenedSettings
	}
	return nil
}
