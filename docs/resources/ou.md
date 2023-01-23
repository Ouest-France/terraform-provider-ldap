# ldap_ou

`ldap_ou` is a resource for managing an LDAP OU.

## Example Usage

```hcl
resource "ldap_ou" "ou" {
  name        = "MyOU"
  ou          = "OU=MyCompany,DC=domain,DC=tld"
  description = "My OU description"
}
```

## Argument Reference

* `ou` - (Required) OU where LDAP OU will be created.
* `name` - (Required) LDAP OU name.
* `description` - (Optional) Description attribute for the LDAP OU. Defaults to empty.
* `managed_by` - (Optional) ManagedBy attribute. Defaults to ``.

## Attribute Reference

* `id` - The DN of the LDAP OU.

## Import

LDAP OU can be imported using the full LDAP DN (id), e.g.

```
$ terraform import ldap_ou.example OU=Myou,OU=MyCompany,DC=domain,DC=tld
```