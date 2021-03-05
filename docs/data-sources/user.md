# ldap_user Data Source

`ldap_user` is a data source for retrieving an LDAP user.

## Example Usage

```hcl
data "ldap_user" "user" {
  ou          = "OU=MyOU,DC=domain,DC=tld"
  name        = "MyUser"
}
```

## Argument Reference

* `ou` - (Required) OU where LDAP user will be search.
* `name` - (Optional) The name of the LDAP user.
* `sam_account_name` - (Optional) The sAMAccountName of the LDAP user.
* `user_principal_name` - (Optional) The userPrincipalName of the LDAP user.

## Attribute Reference

* `id` - LDAP user DN.
* `description` - Description attribute for the LDAP user.
