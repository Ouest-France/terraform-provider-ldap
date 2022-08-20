package ldap

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_HOST", nil),
				Description: "LDAP host",
			},
			"port": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_PORT", nil),
				Description: "LDAP port",
			},
			"bind_user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_USER", nil),
				Description: "LDAP username",
			},
			"bind_password": {
				Type:        schema.TypeString,
				DefaultFunc: schema.EnvDefaultFunc("LDAP_PASSWORD", nil),
				Required:    true,
				Description: "LDAP password",
			},
			"tls": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable the TLS encryption for LDAP (LDAPS). Default, is `false`.",
			},
			"tls_insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Don't verify the server TLS certificate. Default is `false`.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ldap_user": dataSourceLDAPUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := &Client{
		Host:         d.Get("host").(string),
		Port:         d.Get("port").(int),
		BindUser:     d.Get("bind_user").(string),
		BindPassword: d.Get("bind_password").(string),
		TLS:          d.Get("tls").(bool),
		TLSInsecure:  d.Get("tls_insecure").(bool),
	}

	err := client.Connect()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
