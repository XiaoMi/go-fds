package fds

type aclOption struct {
	ACL string `param:"acl" header:"-"`
}

// GrantKey is key of Grant
type GrantKey struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// GrantType is USER or GROUP
type GrantType string

// GrantType const
const (
	GrantTypeUser  GrantType = "USER"
	GrantTypeGroup GrantType = "GROUP"
)

// GrantPermission is permission of Grantee
type GrantPermission string

// GrantPermission const
const (
	GrantPermissionRead        GrantPermission = "READ"
	GrantPermissionWrite       GrantPermission = "WRITE"
	GrantPermissionReadObjects GrantPermission = "READ_OBJECTS"
	GrantPermissionSSOWrite    GrantPermission = "SSO_WRITE"
	GrantPermissionFullControl GrantPermission = "FULL_CONTROL"
)

// Grant grants
type Grant struct {
	Grantee    GrantKey        `json:"grantee"`
	Permission GrantPermission `json:"permission"`
	Type       GrantType       `json:"type"`
}

// AccessControlList is access control list
type AccessControlList struct {
	Grants []Grant `json:"accessControlList"`
	Owner  Owner   `json:"owner"`
}

// AddGrant add a grant into ACL
func (acl *AccessControlList) AddGrant(grant Grant) {
	acl.Grants = append(acl.Grants, grant)
}
