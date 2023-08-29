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
* `scope` - (Optional) LDAP search scope (0: BaseObject, 1: SingleLevel, 2: WholeSubtree) Defaults to `0`.

## Attribute Reference

* `description` - Description attribute for the LDAP
* `group_type` - Type of the group
* `id` - The DN of the LDAP group.
* `members` - LDAP DN of group members
* `members_names` - LDAP name of group members
* `managed_by` - ManagedBy attribute.
* `display_name` - The displayName of the group.