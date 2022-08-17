package ldap

import (
	"context"
	"encoding/json"
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
			"ou": {
				Description: "OU where LDAP user will be search.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description:  "The name of the LDAP user.",
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"name", "sam_account_name", "user_principal_name", "filter"},
			},
			"sam_account_name": {
				Description:  "The sAMAccountName of the LDAP user.",
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"name", "sam_account_name", "user_principal_name", "filter"},
			},
			"user_principal_name": {
				Description:  "The userPrincipalName of the LDAP user",
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"name", "sam_account_name", "user_principal_name", "filter"},
			},
			"filter": {
				Description:  "The filter for selecting the LDAP user.",
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"name", "sam_account_name", "user_principal_name", "filter"},
			},
			"description": {
				Description: "Description attribute for the LDAP user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"data_json": {
				Description: "JSON-encoded string that that is read as the attributes of the user.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func dataSourceLDAPUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceLDAPUserRead(context.WithValue(ctx, CallerTypeKey, DatasourceCaller), d, m)
}

func resourceLDAPUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	ou := d.Get("ou").(string)
	filter := d.Get("filter").(string)

	var user map[string][]string
	var err error
	if _, ok := d.GetOk("filter"); ok {
		user, err = client.ReadUserByFilter(ou, "("+filter+")")
	} else {
		user, err = client.ReadUser(ou, d.Get("name").(string), d.Get("sam_account_name").(string), d.Get("user_principal_name").(string))
	}

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

	var id string
	if _, ok := d.GetOk("filter"); ok {
		id = "(" + filter + "," + ou + ")"
	} else {
		id = user["distinguishedName"][0]
	}
	d.SetId(id)

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

	jsonData, err := json.Marshal(user)
	if err != nil {
		return diag.Errorf("error marshaling JSON for %q: %s", d.Get("name").(string), err)
	}

	if err := d.Set("data_json", string(jsonData)); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}
