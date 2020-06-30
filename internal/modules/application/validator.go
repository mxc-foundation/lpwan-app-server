package application

import (
	"github.com/mxc-foundation/lpwan-app-server/internal/otp"
)

type validator struct {
	otpValidator *otp.Validator
}

func NewValidator(otpValidator *otp.Validator) *validator {
	return &validator{otpValidator: otpValidator}
}
