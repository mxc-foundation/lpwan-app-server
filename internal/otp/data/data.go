package data

// TOTPInfo contains user's TOTP configuration
type TOTPInfo struct {
	Enabled       bool
	Secret        string
	LastTimeSlot  int64
	RecoveryCodes map[int64]string
}
