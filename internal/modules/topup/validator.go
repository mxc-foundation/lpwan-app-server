package topup

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

type validator struct {
	otpValidator *otp.Validator
}

func NewValidator(otpValidator *otp.Validator) *validator {
	return &validator{otpValidator: otpValidator}
}

// API key subjects.
const (
	SubjectUser   = "user"
	SubjectAPIKey = "api_key"
)

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
	FinishRegistration
)
