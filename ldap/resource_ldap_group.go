package ldap

import (
	"context"
	"fmt"
	"strings"

	"github.com/Ouest-France/goldap"
	"github.com/go-ldap/ldap/v3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLDAPGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "`ldap_group` is a resource for managing an LDAP group.",
		CreateContext: resourceLDAPGroupCreate,
		ReadContext:   resourceLDAPGroupRead,
		UpdateContext: resourceLDAPGroupUpdate,
		DeleteContext: resourceLDAPGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The DN of the LDAP group.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ou": {
				Description: "OU where LDAP group will be created.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "LDAP group name.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "Description attribute for the LDAP group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"members": {
				Description: " LDAP group members.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"group_type": {
				Description: "Type of the group",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceLDAPGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	dn := fmt.Sprintf("CN=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	members := []string{}
	memberSet := d.Get("members").(*schema.Set)
	for _, member := range memberSet.List() {
		members = append(members, member.(string))
	}

	err := client.CreateGroup(dn, d.Get("name").(string), d.Get("description").(string), members)
	if err != nil {
		return diag.FromErr(err)
	}
	groupType := d.Get("group_type").(string)
	if groupType != "" {
		err := client.UpdateGroupType(dn, groupType)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(dn)

	return resourceLDAPGroupRead(ctx, d, m)
}

func resourceLDAPGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	dn := d.Id()

	attributes, err := client.ReadGroup(dn)
	if err != nil {
		if err.(*ldap.Error).ResultCode == 32 {
			// Object doesn't exist

			// If Read is called from a datasource, return an error
			if ctx.Value(CallerTypeKey) == DatasourceCaller {
				return diag.Errorf("LDAP group not found: %s", dn)
			}

			// If not a call from datasource, remove the resource from the state
			// and cleanly return
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if err := d.Set("name", attributes["name"][0]); err != nil {
		return diag.FromErr(err)
	}

	// Remove the `CN=<group-name>,` from the DN to get the OU
	ou := strings.ReplaceAll(dn, fmt.Sprintf("CN=%s,", attributes["name"][0]), "")
	if err := d.Set("ou", ou); err != nil {
		return diag.FromErr(err)
	}

	desc := ""
	if val, ok := attributes["description"]; ok {
		desc = val[0]
	}
	if err := d.Set("description", desc); err != nil {
		return diag.FromErr(err)
	}
	groupType := ""
	if val, ok := attributes["groupType"]; ok {
		groupType = val[0]
	}
	if err := d.Set("group_type", groupType); err != nil {
		return diag.FromErr(err)
	}

	members := []string{}
	for name, values := range attributes {
		if name == "member" && len(values) >= 1 {
			members = append(members, values...)
		}
	}
	err = d.Set("members", members)

	return diag.FromErr(err)
}

func resourceLDAPGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)
	dn := fmt.Sprintf("CN=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	if d.HasChange("members") {
		members := []string{}
		memberSet := d.Get("members").(*schema.Set)
		for _, member := range memberSet.List() {
			members = append(members, member.(string))
		}

		if err := client.UpdateGroupMembers(dn, members); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("description") {
		if err := client.UpdateGroupDescription(dn, d.Get("description").(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("group_type") {
		if err := client.UpdateGroupType(dn, d.Get("group_type").(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLDAPGroupRead(ctx, d, m)
}

func resourceLDAPGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	dn := fmt.Sprintf("CN=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	err := client.DeleteGroup(dn)

	return diag.FromErr(err)
}
