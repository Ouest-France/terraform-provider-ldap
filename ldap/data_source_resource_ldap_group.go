package ldap

import (
	"context"
	"fmt"

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
			"ou": &schema.Schema{
				Description: "OU where LDAP group will be search.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": &schema.Schema{
				Description: "LDAP group name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": &schema.Schema{
				Description: "Description attribute for the LDAP",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"members": &schema.Schema{
				Description: "LDAP group members.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_type": &schema.Schema{
				Description: "Type of the group",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceLDAPGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dn := fmt.Sprintf("CN=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	d.SetId(dn)

	// Add context key to signal the Read is called from a datasource
	return resourceLDAPGroupRead(context.WithValue(ctx, CallerTypeKey, DatasourceCaller), d, m)
}
