package main

import (
	"math/rand"
	"time"

	"github.com/TudorHulban/bCRM/pkg/commons"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func checkPasswordHash(pass, salt, hash string) bool {
	errCompare := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass+salt))
	return errCompare == nil
}

func createJWT(ttlSeconds int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Jon Snow"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Second * time.Duration(ttlSeconds)).Unix()

	return token.SignedString([]byte(commons.TokenSecret))
}

func randomString(length int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	result := make([]byte, length)
	for i := range result {
		result[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(result)
}
