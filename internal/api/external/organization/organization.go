package organization

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/lib/pq/hstore"

	"github.com/brocaar/chirpstack-api/go/v3/ns"
	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/mxc-foundation/lpwan-app-server/internal/api/external/dp"
	nsapi "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	nsd "github.com/mxc-foundation/lpwan-app-server/internal/api/external/ns"
	auth "github.com/mxc-foundation/lpwan-app-server/internal/authentication"
	appd "github.com/mxc-foundation/lpwan-app-server/internal/modules/application/data"
	spmod "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile"
	spd "github.com/mxc-foundation/lpwan-app-server/internal/modules/service-profile/data"
	"github.com/mxc-foundation/lpwan-app-server/internal/nscli"
)

var organizationNameRegexp = regexp.MustCompile(`^[\w-]+$`)

// Organization represents an organization.
type Organization struct {
	ID              int64     `db:"id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	Name            string    `db:"name"`
	DisplayName     string    `db:"display_name"`
	CanHaveGateways bool      `db:"can_have_gateways"`
	MaxDeviceCount  int       `db:"max_device_count"`
	MaxGatewayCount int       `db:"max_gateway_count"`
}

// Validate validates the data of the Organization.
func (o Organization) Validate() error {
	if !organizationNameRegexp.MatchString(o.Name) {
		return errors.New("ErrOrganizationInvalidName")
	}
	return nil
}

// OrgUser represents an organization user.
type OrgUser struct {
	UserID         int64     `db:"user_id"`
	Email          string    `db:"email"`
	IsAdmin        bool      `db:"is_admin"`
	IsDeviceAdmin  bool      `db:"is_device_admin"`
	IsGatewayAdmin bool      `db:"is_gateway_admin"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// OrgFilters provides filters for filtering organizations.
type OrgFilters struct {
	UserID int64  `db:"user_id"`
	Search string `db:"search"`

	// Limit and Offset are added for convenience so that this struct can
	// be given as the arguments.
	Limit  int `db:"limit"`
	Offset int `db:"offset"`
}

// SQL returns the SQL filters.
func (f OrgFilters) SQL() string {
	var filters []string

	if f.UserID != 0 {
		filters = append(filters, "u.id = :user_id")
	}

	if f.Search != "" {
		filters = append(filters, "o.display_name ilike :search")
	}

	if len(filters) == 0 {
		return ""
	}

	return "where " + strings.Join(filters, " and ")
}

// Validator defines struct type for vadidating user access to APIs provided by this package
type Validator struct {
	Credentials *auth.Credentials
	st          Store
}

// Validate defines methods used on struct Validator
type Validate interface {
	ValidateOrganizationAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error)
	ValidateOrganizationsAccess(ctx context.Context, flag auth.Flag) (bool, error)
	ValidateOrganizationUsersAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error)
	GetUser(ctx context.Context) (auth.User, error)
}

// NewValidator returns new Validate instance for this package
func NewValidator(st Store) Validate {
	return &Validator{
		Credentials: auth.NewCredentials(),
		st:          st,
	}
}

// GetUser returns user and corresponding attributes after authenticating the user
func (v *Validator) GetUser(ctx context.Context) (auth.User, error) {
	return v.Credentials.GetUser(ctx)
}

// ValidateOrganizationAccess validates if the client has access to the
// given organization.
func (v *Validator) ValidateOrganizationAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationAccess")
	}

	switch flag {
	case auth.Read:
		return v.st.CheckReadOrganizationAccess(ctx, u.Email, u.ID, organizationID)
	case auth.Update:
		return v.st.CheckUpdateOrganizationAccess(ctx, u.Email, u.ID, organizationID)
	case auth.Delete:
		return v.st.CheckDeleteOrganizationAccess(ctx, u.Email, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationsAccess validates if the client has access to the
// organizations.
func (v *Validator) ValidateOrganizationsAccess(ctx context.Context, flag auth.Flag) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationsAccess")
	}

	switch flag {
	case auth.Create:
		return v.st.CheckCreateOrganizationAccess(ctx, u.Email, u.ID)
	case auth.List:
		return v.st.CheckListOrganizationAccess(ctx, u.Email, u.ID)
	default:
		panic("unsupported flag")
	}
}

// ValidateOrganizationUsersAccess validates if the client has access to
// the organization users.
func (v *Validator) ValidateOrganizationUsersAccess(ctx context.Context, flag auth.Flag, organizationID int64) (bool, error) {
	u, err := v.Credentials.GetUser(ctx)
	if err != nil {
		return false, errors.Wrap(err, "ValidateOrganizationUsersAccess")
	}

	switch flag {
	case auth.Create:
		return v.st.CheckCreateOrganizationUserAccess(ctx, u.Email, u.ID, organizationID)
	case auth.List:
		return v.st.CheckListOrganizationUserAccess(ctx, u.Email, u.ID, organizationID)
	default:
		panic("unsupported flag")
	}
}

// Store defines db APIs used by this package
type Store interface {
	GetDefaultNetworkServer(ctx context.Context) (nsd.NetworkServer, error)
	CreateApplication(ctx context.Context, item *appd.Application) error

	CheckReadOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error)
	CheckUpdateOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error)
	CheckDeleteOrganizationAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error)
	CheckCreateOrganizationAccess(ctx context.Context, username string, userID int64) (bool, error)
	CheckListOrganizationAccess(ctx context.Context, username string, userID int64) (bool, error)
	CheckCreateOrganizationUserAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error)
	CheckListOrganizationUserAccess(ctx context.Context, username string, userID int64, organizationID int64) (bool, error)
	CheckReadOrganizationUserAccess(ctx context.Context, username string, organizationID int64, userID, operatorUserID int64) (bool, error)
}

// DefaultApplicationName defines name of the default application for given org id
// this value is unique: default_application_ORGID
const DefaultApplicationName = "default_application_"

// DefaultDeviceProfileName defines name of the default device profile for given org id and network server id
// this value is unique: default_device_profile_ORGID
const DefaultDeviceProfileName = "default_device_profile_"

// ActivateOrganization creates all necessary default settings for new organization:
// default service profile, default applicaiton, default device profile
func ActivateOrganization(ctx context.Context, st Store, spStore spmod.Store, dpSt dp.Store,
	organizationID int64, nsCli *nscli.Client) {
	// get default network
	n, err := st.GetDefaultNetworkServer(ctx)
	if err != nil {
		logrus.WithError(err).Error("couldn't get default network server")
		// no need to interrupt new user registration for this error
		return
	}

	// create default service profile
	sp := spd.ServiceProfile{
		NetworkServerID: n.ID,
		OrganizationID:  organizationID,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
		Name:            fmt.Sprintf(nsapi.DefaultServiceProfileName+"%d", organizationID),
		ServiceProfile: ns.ServiceProfile{
			AddGwMetadata:    true,
			NwkGeoLoc:        true,
			DevStatusReqFreq: 0,
			DrMin:            0,
			DrMax:            0,
			ChannelMask:      []byte(""),
		},
	}
	spID, err := spmod.CreateServiceProfile(ctx, spStore, &sp, nsCli)
	if err != nil {
		logrus.WithError(err).Error("couldn't create default service profile")
		return
	}

	// create default application
	err = st.CreateApplication(ctx, &appd.Application{
		Name:                 fmt.Sprintf(DefaultApplicationName+"%d", organizationID),
		Description:          fmt.Sprintf("default application for organization %d", organizationID),
		OrganizationID:       organizationID,
		ServiceProfileID:     *spID,
		PayloadCodec:         "",
		PayloadEncoderScript: "",
		PayloadDecoderScript: "",
	})
	if err != nil {
		logrus.WithError(err).Error("couldn't create default application")
		return
	}

	// create default device profile
	err = dp.CreateDeviceProfile(ctx, dpSt, nsCli, &dp.DeviceProfile{
		NetworkServerID:      n.ID,
		OrganizationID:       organizationID,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
		Name:                 fmt.Sprintf(DefaultDeviceProfileName+"org_%d_ns_%d", organizationID, n.ID),
		PayloadCodec:         "",
		PayloadEncoderScript: "",
		PayloadDecoderScript: "",
		Tags:                 hstore.Hstore{},
		UplinkInterval:       0,
		DeviceProfile:        ns.DeviceProfile{},
	})
	if err != nil {
		logrus.WithError(err).Error("couldn't creat default device profile")
		return
	}
}
