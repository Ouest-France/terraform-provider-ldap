package ldap

import (
	"context"

	"github.com/Ouest-France/goldap"
	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLDAPUser() *schema.Resource {
	return &schema.Resource{
		Description: "`ldap_user` is a data source for retrieving an LDAP user.",
		ReadContext: dataSourceLDAPUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The DN of the LDAP user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ou": &schema.Schema{
				Description: "OU where LDAP user will be search.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": &schema.Schema{
				Description:  "The name of the LDAP user.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "sam_account_name", "user_principal_name"},
			},
			"sam_account_name": &schema.Schema{
				Description: "The sAMAccountName of the LDAP user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"user_principal_name": &schema.Schema{
				Description: "The userPrincipalName of the LDAP user",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"description": &schema.Schema{
				Description: "Description attribute for the LDAP user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func dataSourceLDAPUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceLDAPUserRead(context.WithValue(ctx, CallerTypeKey, DatasourceCaller), d, m)
}

func resourceLDAPUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	user, err := client.ReadUser(d.Get("ou").(string), d.Get("name").(string), d.Get("sam_account_name").(string), d.Get("user_principal_name").(string))

	if err != nil {
		if err.(*ldap.Error).ResultCode == ldap.LDAPResultNoSuchObject {
			// Object doesn't exist

			// If Read is called from a datasource, return an error
			if ctx.Value(CallerTypeKey) == DatasourceCaller {
				return diag.FromErr(err)
			}

			// If not a call from datasource, remove the resource from the state
			// and cleanly return
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId(user["distinguishedName"][0])

	if val, ok := user["name"]; ok {
		if err := d.Set("name", val[0]); err != nil {
			return diag.FromErr(err)
		}
	}

	if val, ok := user["sAMAccountName"]; ok {
		if err := d.Set("sam_account_name", val[0]); err != nil {
			return diag.FromErr(err)
		}
	}

	if val, ok := user["userPrincipalName"]; ok {
		if err := d.Set("user_principal_name", val[0]); err != nil {
			return diag.FromErr(err)
		}
	}

	if val, ok := user["description"]; ok {
		if err := d.Set("description", val[0]); err != nil {
			return diag.FromErr(err)
		}
	}

	return diag.FromErr(err)
}
