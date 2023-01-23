# ldap_ou

`ldap_ou` is a data source for managing an LDAP OU.

## Example Usage

```hcl
data "ldap_ou" "ou" {
  ou          = "OU=MyCompany,DC=domain,DC=tld"
  name        = "MyOU"
  scope = 2
}
```

## Argument Reference

* `name` - (Required) LDAP OU name.
* `ou` - (Required) OU where LDAP OU will be search.
* `scope` - (Optional) LDAP search scope (0: BaseObject, 1: SingleLevel, 2: WholeSubtree) Defaults to `0`.

## Attribute Reference

* `id` - The DN of the LDAP OU.
* `description` - Description attribute for the LDAP OU
* `managed_by` - ManagedBy attribute.