package ldap

import (
	"context"
	"fmt"

	"github.com/Ouest-France/goldap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLDAPOU() *schema.Resource {
	return &schema.Resource{
		Description: "`ldap_ou` is a data source for getting an LDAP OU.",
		ReadContext: dataSourceLDAPOURead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The DN of the LDAP OU.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ou": {
				Description: "OU where LDAP OU will be search.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "LDAP OU name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"scope": {
				Description: "LDAP search scope",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"description": {
				Description: "Description attribute for the LDAP OU",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"managed_by": {
				Description: "ManagedBy attribute",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceLDAPOURead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Get scope
	scope := d.Get("scope").(int)
	if scope < 0 || scope > 2 {
		return diag.FromErr(fmt.Errorf("scope must be between 0 and 2, got %d", scope))
	}

	// If scope is 1 or 2, we search the OU DN given the OU name and the OU
	client := m.(*goldap.Client)

	// Search OU
	dn, err := client.SearchOUByName(d.Get("name").(string), d.Get("ou").(string), scope)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	// Add context key to signal the Read is called from a datasource
	return resourceLDAPOURead(context.WithValue(ctx, CallerTypeKey, DatasourceCaller), d, m)
}
