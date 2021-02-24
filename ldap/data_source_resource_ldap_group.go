package ldap

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLDAPGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLDAPGroupRead,

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

func dataSourceLDAPGroupRead(d *schema.ResourceData, m interface{}) error {
	dn := fmt.Sprintf("CN=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	d.SetId(dn)

	return resourceLDAPGroupRead(d, m)
}
