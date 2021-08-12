package terraform_provider_report_portal

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	rpClient "github.com/rmalveis/report-portal-client-go/v1/client"
	"strconv"
	"strings"
	"time"
)

func resourceAuthLdapSettings() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAuthLdapSettingsCreate,
		ReadContext:   resourceAuthLdapSettingsRead,
		UpdateContext: resourceAuthLdapSettingsUpdate,
		DeleteContext: resourceAuthLdapSettingsDelete,
		Schema: map[string]*schema.Schema{
			"ldap_attrs_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"ldap_attrs_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ldap_attrs_base_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ldap_attrs_sync_email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ldap_attrs_sync_fullname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ldap_attrs_sync_photo": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_dn_pattern": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_search_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_search_base": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_search_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password_encoder_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(rpClient.PasswordEncryptionTypes, true),
			},
			"password_attr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"manager_dn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"manager_password": {
				Type:      schema.TypeString,
				Sensitive: true,
				Optional:  true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAuthLdapSettingsDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := i.(*rpClient.Client)
	integrationId, err := strconv.Atoi(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = client.DeleteIntegration(&integrationId)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")

	return diags
}

func resourceAuthLdapSettingsRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := i.(*rpClient.Client)
	ldapSettings, err := client.ReadLdapAuthSettings()
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return diags
		}
		return diag.FromErr(err)
	}

	err = data.Set("ldap_attrs_enabled", ldapSettings.LdapAttributes.Enabled)
	err = data.Set("ldap_attrs_url", ldapSettings.LdapAttributes.Url)
	err = data.Set("ldap_attrs_base_dn", ldapSettings.LdapAttributes.BaseDn)
	err = data.Set("ldap_attrs_sync_email", ldapSettings.LdapAttributes.SynchronizationAttributes.Email)
	err = data.Set("ldap_attrs_sync_fullname", ldapSettings.LdapAttributes.SynchronizationAttributes.FullName)
	err = data.Set("ldap_attrs_sync_photo", ldapSettings.LdapAttributes.SynchronizationAttributes.Photo)
	err = data.Set("user_dn_pattern", ldapSettings.UserDnPattern)
	err = data.Set("user_search_filter", ldapSettings.UserSearchFilter)
	err = data.Set("group_search_base", ldapSettings.GroupSearchBase)
	err = data.Set("group_search_filter", ldapSettings.GroupSearchFilter)
	err = data.Set("password_encoder_type", ldapSettings.PasswordEncoderType)
	err = data.Set("password_attr", ldapSettings.PasswordAttribute)
	err = data.Set("manager_password", ldapSettings.ManagerPassword)
	err = data.Set("password_attr", ldapSettings.ManagerDn)

	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAuthLdapSettingsCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(*rpClient.Client)

	var ldapSettings rpClient.LdapIntegrationParameters
	getSettingsFromData(&ldapSettings, data)

	savedSettings, err := client.CreateAuthLdapSettings(&ldapSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(*savedSettings.Id))
	resourceAuthLdapSettingsRead(ctx, data, i)

	return diags
}

func getSettingsFromData(ldapSettings *rpClient.LdapIntegrationParameters, data *schema.ResourceData) {
	ldapSettings.Enabled = data.Get("ldap_attrs_enabled").(bool)
	ldapSettings.Url = data.Get("ldap_attrs_url").(string)
	ldapSettings.BaseDn = data.Get("ldap_attrs_base_dn").(string)
	ldapSettings.Email = data.Get("ldap_attrs_sync_email").(string)
	ldapSettings.FullName = data.Get("ldap_attrs_sync_fullname").(string)
	ldapSettings.Photo = data.Get("ldap_attrs_sync_photo").(string)
	ldapSettings.UserDnPattern = data.Get("user_dn_pattern").(string)
	ldapSettings.UserSearchFilter = data.Get("user_search_filter").(string)
	ldapSettings.GroupSearchBase = data.Get("group_search_base").(string)
	ldapSettings.GroupSearchFilter = data.Get("group_search_filter").(string)
	ldapSettings.PasswordEncoderType = data.Get("password_encoder_type").(string)
	ldapSettings.PasswordAttribute = data.Get("password_attr").(string)
	ldapSettings.ManagerDn = data.Get("manager_dn").(string)
	ldapSettings.ManagerPassword = data.Get("manager_password").(string)
}

func resourceAuthLdapSettingsUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := i.(*rpClient.Client)

	var ldapSettings rpClient.LdapIntegrationParameters
	getSettingsFromData(&ldapSettings, data)

	savedSettings, err := client.UpdateAuthLdapSettings(&ldapSettings)
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(strconv.Itoa(*savedSettings.Id))

	err = data.Set("latest_updated", time.Now().Format(time.RFC3339))
	if err != nil {
		return nil
	}
	resourceAuthLdapSettingsRead(ctx, data, i)

	return diags
}
