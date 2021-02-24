package ldap

import (
	"github.com/Ouest-France/goldap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "LDAP host",
			},
			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "LDAP port",
			},
			"bind_user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "LDAP username",
			},
			"bind_password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "LDAP password",
			},
			"tls": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable the TLS encryption for LDAP (LDAPS). Default, is `false`.",
			},
			"tls_ca_certificate": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The TLS CA certificate to trust for the LDAPS connection.",
			},
			"tls_insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Don't verify the server TLS certificate. Default is `false`.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ldap_group": resourceLDAPGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"ldap_group": dataSourceLDAPGroup(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	client := &goldap.Client{
		Host:         d.Get("host").(string),
		Port:         d.Get("port").(int),
		BindUser:     d.Get("bind_user").(string),
		BindPassword: d.Get("bind_password").(string),
		TLS:          d.Get("tls").(bool),
		TLSCACert:    d.Get("tls_ca_certificate").(string),
		TLSInsecure:  d.Get("tls_insecure").(bool),
	}

	err := client.Connect()
	if err != nil {
		return nil, err
	}

	return client, nil
}
