package ldap

import (
	"github.com/Ouest-France/goldap"
	"github.com/hashicorp/terraform/helper/schema"
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
				Description: "FortiADC username",
			},
			"bind_password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "FortiADC password",
			},
			"tls": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable TLS Verify",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ldap_group": resourceLDAPGroup(),
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
	}

	err := client.Connect()
	if err != nil {
		return nil, err
	}

	return client, nil
}
