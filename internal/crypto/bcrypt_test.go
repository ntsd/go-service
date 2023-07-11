package crypto

import (
	"testing"
)

func TestBcryptHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		salt     string
		wantErr  bool
	}{
		{
			name:     "success",
			password: "password",
			salt:     "salt",
		},
		{
			name:     "success, no salt",
			password: "password",
			salt:     "",
		},
		{
			name:     "success. no password",
			password: "",
			salt:     "salt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := BcryptHash(tt.password, tt.salt)
			if (err != nil) != tt.wantErr {
				t.Errorf("BcryptHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			BcryptVerify(hash, tt.password, tt.salt)
		})
	}
}
