package crypt

import "golang.org/x/crypto/bcrypt"

func Encode(txt string) (string, error) {

	txtEncoded, error := bcrypt.GenerateFromPassword([]byte(txt), bcrypt.MinCost)

	if error != nil {
		return "", error
	}
	return string(txtEncoded[:]), nil

}
func Compare(txt string, txtEncoded string) bool {

	error := bcrypt.CompareHashAndPassword([]byte(txtEncoded), []byte(txt))

	return error == nil

}
