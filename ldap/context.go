package ldap

// CallerTypeKeyType represents an CallerType context key type
type CallerTypeKeyType string

// DatasourceCallerValueType represents an DatasourceCaller context key type
type DatasourceCallerValueType bool

// CallerTypeKey represents an CallerTypeKey context key
const CallerTypeKey CallerTypeKeyType = "callerType"

// DatasourceCaller represents an DatasourceCallerKeyValue context value
const DatasourceCaller DatasourceCallerValueType = true
