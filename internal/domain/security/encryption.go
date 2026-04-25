package security

import "golang.org/x/crypto/bcrypt"

func HashSenha(senha string) ([]byte, error) {

	return bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
}

func VerificaSenha(senhaHashed, senha string) error {

	return bcrypt.CompareHashAndPassword([]byte(senhaHashed), []byte(senha))
}
