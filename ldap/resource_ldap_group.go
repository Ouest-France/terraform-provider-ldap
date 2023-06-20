package ldap

import (
	"context"
	"fmt"
	"regexp"

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
				Description: " LDAP group members DN",
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"members_names": {
				Description: "LDAP group members names.",
				Type:        schema.TypeSet,
				Computed:    true,
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
			"managed_by": {
				Description: "ManagedBy attribute",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
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

	err := client.CreateGroup(dn, d.Get("name").(string), d.Get("description").(string), d.Get("group_type").(string), d.Get("managed_by").(string), members)
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

	attributes, err := client.ReadGroup(dn, 1500)
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

	nameAttr, ok := attributes["name"]
	if !ok || len(nameAttr) != 1 {
		return diag.Errorf("LDAP attribute \"name\" doesn't exist or is empty for group: %s", dn)
	}

	if err := d.Set("name", nameAttr[0]); err != nil {
		return diag.FromErr(err)
	}

	// Remove the `CN=<group-name>` from the DN to get the OU
	// using the regex `^cn=.*?,(.*)$` in case insensitive mode
	reg := regexp.MustCompile(`(?i)^cn=.*?,(.*)$`)
	match := reg.FindStringSubmatch(dn)
	if len(match) != 2 {
		return diag.Errorf("Failed parsing OU from DN (must match regex `^cn=.*?,(.*)$`): %s", dn)
	}
	if err := d.Set("ou", match[1]); err != nil {
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
	managedBy := ""
	if val, ok := attributes["managedBy"]; ok {
		managedBy = val[0]
	}
	if err := d.Set("managed_by", managedBy); err != nil {
		return diag.FromErr(err)
	}

	members := []string{}
	members_names := []string{}
	for name, values := range attributes {
		if name == "member" && len(values) >= 1 {
			members = append(members, values...)

			for _, member := range values {
				regName := regexp.MustCompile(`CN=(.*?),`)
				matches := regName.FindStringSubmatch(member)
				if len(matches) == 2 {
					members_names = append(members_names, matches[1])
				}
			}
		}
	}
	err = d.Set("members", members)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("members_names", members_names)

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

	if d.HasChange("managed_by") {
		if err := client.UpdateGroupManagedBy(dn, d.Get("managed_by").(string)); err != nil {
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
