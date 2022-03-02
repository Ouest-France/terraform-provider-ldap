# ldap_group

`ldap_group` is a data source for managing an LDAP group.

## Example Usage

```hcl
data "ldap_group" "group" {
  ou          = "OU=MyOU,DC=domain,DC=tld"
  name        = "MyGroup"
}
```

## Argument Reference

* `name` - (Required) LDAP group name.
* `ou` - (Required) OU where LDAP group will be search.

## Attribute Reference

* `description` - Description attribute for the LDAP
* `group_type` - Type of the group
* `id` - The DN of the LDAP group.
* `members` - LDAP group members.