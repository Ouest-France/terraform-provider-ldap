# LDAP Provider

The LDAP provider is used to interact with any ActiveDirectory server.

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
* `tls` - (Optional) Enable the TLS encryption for LDAP (LDAPS). Default, is `false`.
* `tls_ca_certificate` - (Optional) The TLS CA certificate to trust for the LDAPS connection.
* `tls_insecure` - (Optional) Don't verify the server TLS certificate. Default is `false`.
