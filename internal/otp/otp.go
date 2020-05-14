package otp

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image/png"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// TOTPInfo contains user's TOTP configuration
type TOTPInfo struct {
	Enabled       bool
	Secret        string
	LastTimeSlot  int64
	RecoveryCodes map[int64]string
}

// Store provides methods to access TOTP information for users stored in the database
type Store interface {
	// GetTOTPInfo returns TOTP configuration info for the user
	GetTOTPInfo(ctx context.Context, username string) (TOTPInfo, error)
	// Enable enables TOTP for the user
	Enable(ctx context.Context, username string) error
	// Delete removes TOTP configuration for the user
	Delete(ctx context.Context, username string) error
	// StoreNewSecret stores new secret for user in the database
	StoreNewSecret(ctx context.Context, username, secret string) error
	// DeleteRecoveryCode removes recovery code from the database
	DeleteRecoveryCode(ctx context.Context, username string, codeID int64) error
	// AddRecoveryCodes adds new recovery codes to the database
	AddRecoveryCodes(ctx context.Context, username string, codes []string) error
	// UpdateLastTimeSlot updates last time slot value in the database
	UpdateLastTimeSlot(ctx context.Context, username string, previousValue, newValue int64) error
}

const (
	errOTPNotValid  = "The OTP is not valid"
	errOTPLockedOut = "Too many unsuccessful attemps, try again in 10 minutes"
)

// Validator provides methods to generate TOTP configuration for the user and validate OTPs
type Validator struct {
	issuer string
	block  cipher.Block
	store  Store
}

// NewValidator creates a new TOTP validator using given issuer, master key and store.
func NewValidator(issuer, key string, store Store) (*Validator, error) {
	k, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("couldn't decode the key: %v", err)
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize cipher: %v", err)
	}
	return &Validator{
		issuer: issuer,
		block:  block,
		store:  store,
	}, nil
}

// Configuration contains TOTP configuration generated for user
type Configuration struct {
	URL           string
	Secret        string
	RecoveryCodes []string
	QRCode        string
}

var totpOptions = totp.ValidateOpts{
	Period:    30,
	Digits:    otp.DigitsSix,
	Algorithm: otp.AlgorithmSHA1,
}

// NewConfiguration generates a new TOTP configuration for the user
func (v *Validator) NewConfiguration(ctx context.Context, username string) (*Configuration, error) {
	ti, err := v.store.GetTOTPInfo(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve TOTPInfo for %s: %v", username, err)
	}
	if ti.Enabled {
		return nil, fmt.Errorf("TOTP is already enabled")
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      v.issuer,
		AccountName: username,
		Period:      totpOptions.Period,
		Digits:      totpOptions.Digits,
		SecretSize:  20,
		Algorithm:   totpOptions.Algorithm,
	})
	if err != nil {
		return nil, fmt.Errorf("couldn't generate a new TOTP key: %v", err)
	}

	qrcode, err := encodeImage(key)
	if err != nil {
		return nil, fmt.Errorf("couldn't create QR code: %v", err)
	}

	conf := &Configuration{
		URL:    key.URL(),
		Secret: key.Secret(),
		QRCode: qrcode,
	}

	eSecret, err := v.encrypt(key.Secret())
	if err != nil {
		return nil, fmt.Errorf("couldn't encrypt TOTP secret: %v", err)
	}
	if err := v.store.StoreNewSecret(ctx, username, eSecret); err != nil {
		return nil, fmt.Errorf("couldn't store TOTP secret: %v", err)
	}

	rCodes, err := v.GetRecoveryCodes(ctx, username, true)
	if err != nil {
		return nil, fmt.Errorf("couldn't get recovery codes: %v", err)
	}
	conf.RecoveryCodes = rCodes

	return conf, nil
}

// GetRecoveryCodes returns the list of recovery codes for the user. If
// regenerate is set to true then all the old recovery codes are deleted and
// replaced with the new ones, otherwise the old codes are left in place and
// the new ones only generated to make it to ten.
func (v *Validator) GetRecoveryCodes(ctx context.Context, username string, regenerate bool) ([]string, error) {
	ti, err := v.store.GetTOTPInfo(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("couldn't retrieve TOTPInfo for %s: %v", username, err)
	}
	if regenerate {
		for cid := range ti.RecoveryCodes {
			if err := v.store.DeleteRecoveryCode(ctx, username, cid); err != nil {
				return nil, fmt.Errorf("couldn't remove recovery code: %v", err)
			}
		}
		ti.RecoveryCodes = nil
	}
	var codes []string
	for _, val := range ti.RecoveryCodes {
		dCode, err := v.decrypt(val)
		if err != nil {
			return nil, fmt.Errorf("couldn't decrypt recovery code")
		}
		codes = append(codes, dCode)
	}
	var codesToAdd []string
	for len(codes) < 10 {
		code, err := generateRecoveryCode()
		if err != nil {
			return nil, fmt.Errorf("couldn't generate a recovery code: %v", err)
		}
		codes = append(codes, code)
		eCode, err := v.encrypt(code)
		if err != nil {
			return nil, fmt.Errorf("couldn't encrypt a recovery code: %v", err)
		}
		codesToAdd = append(codesToAdd, eCode)
	}

	if err := v.store.AddRecoveryCodes(ctx, username, codesToAdd); err != nil {
		return nil, fmt.Errorf("couldn't store recovery codes: %v", err)
	}
	return codes, nil
}

// encodes totp keys as Png image and then returns it encoded in base64
func encodeImage(key *otp.Key) (string, error) {
	qrcode, err := key.Image(300, 300)
	if err != nil {
		return "", err
	}
	bRaw := bytes.NewBuffer(nil)
	if err := png.Encode(bRaw, qrcode); err != nil {
		return "", err
	}
	b64 := make([]byte, base64.StdEncoding.EncodedLen(bRaw.Len()))
	base64.StdEncoding.Encode(b64, bRaw.Bytes())
	return string(b64), nil
}

func (v *Validator) encrypt(secret string) (string, error) {
	ciphertext := make([]byte, v.block.BlockSize()+len(secret))
	iv := ciphertext[:v.block.BlockSize()]
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(v.block, iv)
	stream.XORKeyStream(ciphertext[v.block.BlockSize():], []byte(secret))
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(ciphertext)))
	base64.StdEncoding.Encode(b64, ciphertext)
	return string(b64), nil
}

func (v *Validator) decrypt(encrypted string) (string, error) {
	ciphertext := make([]byte, base64.StdEncoding.DecodedLen(len(encrypted)))
	if len(ciphertext) <= v.block.BlockSize() {
		return "", fmt.Errorf("ciphertext contains no data")
	}
	if _, err := base64.StdEncoding.Decode(ciphertext, []byte(encrypted)); err != nil {
		return "", fmt.Errorf("couldn't decode base64: %v", err)
	}
	iv := ciphertext[:v.block.BlockSize()]
	ciphertext = ciphertext[v.block.BlockSize():]
	stream := cipher.NewCFBDecrypter(v.block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

// IsEnabled returns true if user has enabled TOTP
func (v *Validator) IsEnabled(ctx context.Context, username string) (bool, error) {
	ti, err := v.store.GetTOTPInfo(ctx, username)
	if err != nil {
		return false, err
	}
	return ti.Enabled, nil
}

// Enable enables TOTP authentication for user
func (v *Validator) Enable(ctx context.Context, username, otp string) error {
	if err := v.Validate(ctx, username, otp); err != nil {
		return err
	}
	if err := v.store.Enable(ctx, username); err != nil {
		return fmt.Errorf("couldn't enable TOTP: %v", err)
	}
	return nil
}

// Disable disables TOTP for the user and removes configuration from the
// database
func (v *Validator) Disable(ctx context.Context, username string) error {
	return v.store.Delete(ctx, username)
}

// Validate returns nil if otp is correct one-time password or recovery code,
// otherwise it returns the error
func (v *Validator) Validate(ctx context.Context, username, otp string) error {
	ti, err := v.store.GetTOTPInfo(ctx, username)
	if err != nil {
		return fmt.Errorf("couldn't retrieve TOTP information for %s: %v", username, err)
	}

	// make sure that the last time slot is not greater than the current one
	ts := time.Now().Unix() / int64(totpOptions.Period)
	if ts <= ti.LastTimeSlot {
		// we will increase last time slot to up to ten minutes to prevent bruteforce attacks
		if ti.LastTimeSlot-ts < 20 {
			_ = v.store.UpdateLastTimeSlot(ctx, username, ti.LastTimeSlot, ti.LastTimeSlot+1)
		}
		if ts == ti.LastTimeSlot {
			return fmt.Errorf(errOTPNotValid)
		}
		return fmt.Errorf(errOTPLockedOut)
	}

	if len(otp) == 6 {
		return v.validateTOTP(ctx, username, ti, ts, otp)
	}
	if len(otp) == 11 {
		return v.validateRecoveryCode(ctx, username, ti, ts, otp)
	}
	return fmt.Errorf("invalid OTP format: %s", otp)
}

func (v *Validator) validateTOTP(ctx context.Context, username string, ti TOTPInfo, ts int64, otp string) error {
	secret, err := v.decrypt(ti.Secret)
	if err != nil {
		return fmt.Errorf("couldn't decrypt TOTP secret: %v", err)
	}
	if len(secret) < 32 {
		return fmt.Errorf("invalid secret")
	}
	for i := 0; i < 3 && ts > ti.LastTimeSlot; i++ {
		gotp, err := totp.GenerateCodeCustom(secret, time.Unix(ts*int64(totpOptions.Period), 0), totpOptions)
		if err != nil {
			return fmt.Errorf("couldn't validate the code: %v", err)
		}
		if subtle.ConstantTimeCompare([]byte(gotp), []byte(otp)) == 1 {
			return v.store.UpdateLastTimeSlot(ctx, username, ti.LastTimeSlot, ts)
		}
		ts--
	}
	ts++
	_ = v.store.UpdateLastTimeSlot(ctx, username, ti.LastTimeSlot, ts)
	return fmt.Errorf(errOTPNotValid)
}

func (v *Validator) validateRecoveryCode(ctx context.Context, username string, ti TOTPInfo, ts int64, otp string) error {
	otp = strings.ToLower(otp)
	for id, code := range ti.RecoveryCodes {
		rc, err := v.decrypt(code)
		if err != nil {
			return fmt.Errorf("couldn't decrypt recovery code: %v", err)
		}
		if subtle.ConstantTimeCompare([]byte(otp), []byte(rc)) == 1 {
			if err := v.store.DeleteRecoveryCode(ctx, username, id); err != nil {
				return fmt.Errorf(errOTPNotValid)
			}
			return v.store.UpdateLastTimeSlot(ctx, username, ti.LastTimeSlot, ts)
		}
	}
	for i := 0; i < 2 && ts > ti.LastTimeSlot+1; i++ {
		ts--
	}
	_ = v.store.UpdateLastTimeSlot(ctx, username, ti.LastTimeSlot, ts)
	return fmt.Errorf(errOTPNotValid)
}

func generateRecoveryCode() (string, error) {
	binCode := make([]byte, 5)
	_, err := rand.Read(binCode)
	if err != nil {
		return "", err
	}
	code := fmt.Sprintf("%x", binCode)
	return code[:5] + "-" + code[5:], nil
}
