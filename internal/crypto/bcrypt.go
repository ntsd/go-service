package crypto

import "golang.org/x/crypto/bcrypt"

func BcryptHash(password string, salt string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
}

func BcryptVerify(hash []byte, password string, salt string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password+salt))
}
