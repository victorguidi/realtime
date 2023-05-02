package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func encryptData() {}

func hashPassword(password string) (string, error) {
	var hashedPassword []byte
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func unhashPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("Password does not match")
		return false, err
	}
	return true, nil
}
