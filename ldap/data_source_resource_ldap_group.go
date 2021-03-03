package ldap

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLDAPGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLDAPGroupRead,

		Schema: map[string]*schema.Schema{
			"ou": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"members": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
