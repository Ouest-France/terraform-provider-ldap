# ldap_group

`ldap_group` is a resource for managing an LDAP group.

## Example Usage

```hcl
resource "ldap_group" "group" {
  ou          = "OU=MyOU,DC=domain,DC=tld"
  name        = "MyGroup"
  members     = ["CN=MyUser,OU=MyOU,DC=domain,DC=tld"]
  description = "My group description"
}
```

## Argument Reference

* `ou` - (Required) OU where LDAP group will be created.
* `name` - (Required) LDAP group name.
* `members` - (Optional) LDAP group members. Defaults to `[]`.
* `description` - (Optional) Description attribute for the LDAP group. Defaults to empty.
* `group_type` - (Optional, Computed) Type of the group.
* `managed_by` - (Optional) ManagedBy attribute. Defaults to ``.
* `display_name` - (Optional) The displayName of the group. Defaults to ``.

## Attribute Reference

* `members_names` - Names of the members
* `id` - The DN of the LDAP group.

## Import

LDAP group can be imported using the full LDAP DN (id), e.g.

```
$ terraform import ldap_group.example CN=MyGroup,OU=MyOU,DC=domain,DC=tld
```