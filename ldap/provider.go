package ldap

import (
	"github.com/Ouest-France/goldap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
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
			"tls_ca_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The TLS CA certificate to trust for the LDAPS connection.",
			},
			"tls_insecure": {
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
			"ldap_user":  dataSourceLDAPUser(),
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

	if logging.IsDebugOrHigher() {
		client.Conn.Debug.Enable(true)
	}

	return client, nil
}
