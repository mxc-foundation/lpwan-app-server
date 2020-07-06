package authentication

// User contains information about the user
type User struct {
	ID            int64
	Username      string
	IsGlobalAdmin bool
}

// OrgUser contains information about the role of the user in organisation
type OrgUser struct {
	IsOrgUser      bool
	IsOrgAdmin     bool
	IsDeviceAdmin  bool
	IsGatewayAdmin bool
}
