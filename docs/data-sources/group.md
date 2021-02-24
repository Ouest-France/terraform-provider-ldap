# ldap_group Data Source

`ldap_group` is a data source for managing an LDAP group.

## Example Usage

```hcl
data "ldap_group" "group" {
  ou          = "OU=MyOU,DC=domain,DC=tld"
  name        = "MyGroup"
}
```

## Argument Reference

* `ou` - (Required) OU where LDAP group will be created.
* `name` - (Required) LDAP group name.

## Attribute Reference

* `id` - LDAP group DN.
* `members` - LDAP group members.
* `description` - Description attribute for the LDAP group.