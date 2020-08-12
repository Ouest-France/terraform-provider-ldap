# ldap_group

`ldap_group` is a resource for managing an LDAP group.

## Example Usage

```hcl
resource "ldap_group" "group" {
  ou      = "OU=MyOU,DC=domain,DC=tld"
  name    = "MyGroup"
  members = ["CN=MyUser,OU=MyOU,DC=domain,DC=tld"]
}
```

## Argument Reference

* `ou` - (Required) OU where LDAP group will be created.
* `name` - (Required) LDAP group name.
* `members` - (Optional) LDAP group members.

## Attribute Reference

* `id` - LDAP group DN.