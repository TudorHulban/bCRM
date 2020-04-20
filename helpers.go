package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(pPassword, pSalt, pHash string) bool {
	errCompare := bcrypt.CompareHashAndPassword([]byte(pHash), []byte(pPassword+pSalt))
	return errCompare == nil
}

func createJWT(ttlSeconds int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Jon Snow"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Second * time.Duration(ttlSeconds)).Unix()

	return token.SignedString([]byte("secret"))
}
