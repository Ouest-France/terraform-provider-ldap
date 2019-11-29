package ldap

import (
	"fmt"

	"github.com/Ouest-France/goldap"
	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLDAPGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceLDAPGroupCreate,
		Read:   resourceLDAPGroupRead,
		Delete: resourceLDAPGroupDelete,

		Schema: map[string]*schema.Schema{
			"ou": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"members": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceLDAPGroupCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*goldap.Client)

	dn := fmt.Sprintf("CN=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	members := []string{}
	memberSet := d.Get("members").(*schema.Set)
	for _, member := range memberSet.List() {
		members = append(members, member.(string))
	}

	err := client.CreateGroup(dn, d.Get("name").(string), members)
	if err != nil {
		return err
	}

	d.SetId(dn)

	return resourceLDAPGroupRead(d, m)
}

func resourceLDAPGroupRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*goldap.Client)

	dn := fmt.Sprintf("CN=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	attributes, err := client.ReadGroup(dn)
	if err != nil {
		if err.(*ldap.Error).ResultCode == 32 {
			// Object doesn't exist
			d.SetId("")
		}
	}

	members := []string{}
	for name, values := range attributes {

		if name == "member" && len(values) >= 1 {
			members = append(members, values...)
		}
	}
	err = d.Set("members", members)

	return err
}

func resourceLDAPGroupDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*goldap.Client)

	dn := fmt.Sprintf("CN=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	err := client.DeleteGroup(dn)

	return err
}
