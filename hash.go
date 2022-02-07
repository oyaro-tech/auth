package auth

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		log.Println("HashAndSalt: ", err.Error())
	}

	return string(hash)
}

func comparePasswords(clientPwd, serverPwd string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(serverPwd), []byte(clientPwd)); err != nil {
		log.Println("ComparePasswords: ", err.Error())
		return false
	}

	return true
}

func getPasswordHashByLogin(username string) *User {
	var user User

	row, err := db.Query(context.Background(), "SELECT id, password FROM users WHERE username=$1 LIMIT 1", username)
	if err != nil {
		log.Printf("GetPasswordHashByLogin query: %v\n", err.Error())
		return nil
	}

	for row.Next() {
		err = row.Scan(&user.ID, &user.Password)
		if err != nil {
			log.Printf("Row scan: %v\n", err.Error())
			return nil
		}
	}

	return &user
}
