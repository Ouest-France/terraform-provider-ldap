package ldap

import (
	"context"
	"fmt"

	"github.com/Ouest-France/goldap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLDAPGroup() *schema.Resource {
	return &schema.Resource{
		Description: "`ldap_group` is a data source for managing an LDAP group.",
		ReadContext: dataSourceLDAPGroupRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The DN of the LDAP group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ou": {
				Description: "OU where LDAP group will be search.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "LDAP group name.",
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
				Description: "Description attribute for the LDAP",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"members": {
				Description: "LDAP group members.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_type": {
				Description: "Type of the group",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceLDAPGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// Get scope
	scope := d.Get("scope").(int)
	if scope < 0 || scope > 2 {
		return diag.FromErr(fmt.Errorf("scope must be between 0 and 2, got %d", scope))
	}

	if scope == 0 {
		// If scope is 0, we keep the old code to ensure backward compatibility
		dn := fmt.Sprintf("CN=%s,%s", d.Get("name").(string), d.Get("ou").(string))

		d.SetId(dn)
	} else {
		// If scope is 1 or 2, we search the group DN given the group name and the OU
		client := m.(*goldap.Client)

		// Search group
		dn, err := client.SearchGroupByName(d.Get("name").(string), d.Get("ou").(string), scope)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId(dn)
	}

	// Add context key to signal the Read is called from a datasource
	return resourceLDAPGroupRead(context.WithValue(ctx, CallerTypeKey, DatasourceCaller), d, m)
}
