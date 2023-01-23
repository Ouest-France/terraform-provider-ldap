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

func resourceLDAPOU() *schema.Resource {
	return &schema.Resource{
		Description:   "`ldap_ou` is a resource for managing an LDAP OU.",
		CreateContext: resourceLDAPOUCreate,
		ReadContext:   resourceLDAPOURead,
		UpdateContext: resourceLDAPOUUpdate,
		DeleteContext: resourceLDAPOUDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The DN of the LDAP OU.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ou": {
				Description: "OU where LDAP OU will be created.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "LDAP OU name.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "Description attribute for the LDAP OU.",
				Type:        schema.TypeString,
				Optional:    true,
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

func resourceLDAPOUCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	dn := fmt.Sprintf("OU=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	err := client.CreateOrganizationalUnit(dn, d.Get("description").(string), d.Get("managed_by").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(dn)

	return resourceLDAPOURead(ctx, d, m)
}

func resourceLDAPOURead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)

	dn := d.Id()

	attributes, err := client.ReadOrganizationalUnit(dn)
	if err != nil {
		if err.(*ldap.Error).ResultCode == 32 {
			// Object doesn't exist

			// If Read is called from a datasource, return an error
			if ctx.Value(CallerTypeKey) == DatasourceCaller {
				return diag.Errorf("LDAP OU not found: %s", dn)
			}

			// If not a call from datasource, remove the resource from the state
			// and cleanly return
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	nameAttr, ok := attributes["ou"]
	if !ok || len(nameAttr) != 1 {
		return diag.Errorf("LDAP attribute \"name\" doesn't exist or is empty for OU: %s", dn)
	}

	if err := d.Set("name", nameAttr[0]); err != nil {
		return diag.FromErr(err)
	}

	// Remove the `OU=<ou-name>,` from the DN to get the OU
	ou := strings.ReplaceAll(dn, fmt.Sprintf("OU=%s,", attributes["name"][0]), "")
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

	managedBy := ""
	if val, ok := attributes["managedBy"]; ok {
		managedBy = val[0]
	}
	if err := d.Set("managed_by", managedBy); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}

func resourceLDAPOUUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)
	dn := fmt.Sprintf("OU=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	if d.HasChange("description") {
		if err := client.UpdateOrganizationalUnitDescription(dn, d.Get("description").(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("managed_by") {
		if err := client.UpdateOrganizationalUnitManagedBy(dn, d.Get("managed_by").(string)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceLDAPOURead(ctx, d, m)
}

func resourceLDAPOUDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*goldap.Client)
	dn := fmt.Sprintf("OU=%s,%s", d.Get("name").(string), d.Get("ou").(string))

	err := client.DeleteOrganizationalUnit(dn)

	return diag.FromErr(err)
}
