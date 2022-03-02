# ldap_user

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
* `name` - (Required) Name of the LDAP user.
* `sam_account_name` - (Optional) sAMAccountName of the LDAP user.
* `user_principal_name` - (Optional) UPN of the LDAP user.

## Attribute Reference

* `id` - The DN of the LDAP user.
* `description` - Description attribute for the LDAP user.