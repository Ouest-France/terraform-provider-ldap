package ldap

import (
	"fmt"
	"github.com/go-ldap/ldap/v3"
)

func (c *Client) ReadUserByFilter(ou string, filter string) (entries map[string][]string, err error) {
	req := ldap.NewSearchRequest(
		ou,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		[]string{"*"},
		[]ldap.Control{},
	)

	sr, err := c.Conn.Search(req)
	if err != nil {
		return nil, err
	}

	if len(sr.Entries) == 0 {
		return nil, ldap.NewError(ldap.LDAPResultNoSuchObject, fmt.Errorf("The filter '%s' doesn't match any user in the OU: %s", filter, ou))
	}

	if len(sr.Entries) > 1 {
		return nil, ldap.NewError(ldap.LDAPResultOther, fmt.Errorf("The filter '%s' match more than one user in the OU: %s", filter, ou))
	}

	entries = map[string][]string{}
	for _, entry := range sr.Entries {
		for _, attr := range entry.Attributes {
			entries[attr.Name] = attr.Values
		}
	}

	return entries, nil
}
