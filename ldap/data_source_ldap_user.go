package ldap

import (
	"context"
	"encoding/json"
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
			"filter": {
				Description: "The filter for selecting the LDAP user.",
				Type:        schema.TypeString,
				Required:    true,
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
	client := m.(*Client)

	ou := d.Get("ou").(string)
	filter := d.Get("filter").(string)

	user, err := client.ReadUserByFilter(ou, "("+filter+")")

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

	id := "(" + filter + "," + ou + ")"
	d.SetId(id)

	jsonData, err := json.Marshal(user)
	if err != nil {
		return diag.Errorf("error marshaling JSON for %q: %s", id, err)
	}

	if err := d.Set("data_json", string(jsonData)); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}
