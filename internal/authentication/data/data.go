package data

// User contains information about the user
type User struct {
	ID            int64
	Email         string
	IsGlobalAdmin bool
}

// OrgUser contains information about the role of the user in organisation
type OrgUser struct {
	IsOrgUser      bool
	IsOrgAdmin     bool
	IsDeviceAdmin  bool
	IsGatewayAdmin bool
}

// Flag defines the authorization flag.
type Flag int

// Authorization flags.
const (
	Create Flag = iota
	Read
	Update
	Delete
	List
	UpdateProfile
	UpdatePassword
	FinishRegistration
)
