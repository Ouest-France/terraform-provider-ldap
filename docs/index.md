# LDAP Provider

The LDAP provider is used to interact with any LDAP server.

## Example Usage

```hcl
# Provider configuration
terraform {
  required_providers {
    ldap = {
      source  = "Ouest-France/ldap"
    }
  }
}

provider "ldap" {
  host          = "ldap.mycompany.tld"
  port          = 389
  bind_user     = "ldap_user"
  bind_password = "ldap_password"
}

...
```

## Argument Reference

* `host` - (Required) LDAP host address formatted like `ldap.mycompany.com`.
* `port` - (Required) LDAP port.
* `bind_user` - (Required) LDAP bind user.
* `bind_password` - (Required) LDAP bind password.
* `tls` - (Optional) Enable TLS. Defaults to `false`.