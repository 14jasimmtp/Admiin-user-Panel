package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPwd(pwd string) string {
	HashPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("error while hashing password")
	}
	return string(HashPwd)
}
