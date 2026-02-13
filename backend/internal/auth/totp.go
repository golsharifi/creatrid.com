package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TOTPService struct{}

func NewTOTPService() *TOTPService {
	return &TOTPService{}
}

// GenerateSecret creates a new TOTP secret for the given user
func (s *TOTPService) GenerateSecret(email string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "Creatrid",
		AccountName: email,
	})
}

// ValidateCode checks a TOTP code against a secret
func (s *TOTPService) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateBackupCodes creates 8 random 8-digit backup codes
func (s *TOTPService) GenerateBackupCodes() ([]string, error) {
	codes := make([]string, 8)
	for i := range codes {
		n, err := rand.Int(rand.Reader, big.NewInt(100000000))
		if err != nil {
			return nil, err
		}
		codes[i] = fmt.Sprintf("%08d", n.Int64())
	}
	return codes, nil
}
