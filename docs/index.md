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

* `host` - (Required) LDAP host, can also be provided with env var **LDAP_HOST**.

* `port` - (Required) LDAP port, can also be provided with env var **LDAP_PORT**.

* `bind_user` - (Required) LDAP username, can also be provided with env var **LDAP_USER**.

* `bind_password` - (Required) LDAP password, can also be provided with env var **LDAP_PASSWORD**.

* `tls` - (Optional) Enable the TLS encryption for LDAP (LDAPS). Default, is `false`.

* `tls_ca_certificate` - (Optional) The TLS CA certificate to trust for the LDAPS connection. Default is empty.

* `tls_insecure` - (Optional) Don't verify the server TLS certificate. Default is `false`.