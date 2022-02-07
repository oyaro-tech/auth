package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

var (
	db *pgx.Conn
)

func init() {
	var err error

	db, err = pgx.Connect(
		context.Background(),
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_PORT"),
			os.Getenv("POSTGRES_DB"),
		),
	)
	if err != nil {
		log.Fatalf("Can't connect to database: %s\n", err.Error())
	}

	os.Setenv("ACCESS_SECRET", "U5kBnGsmHW1Bchegg7bi8fFvfdqxSAuk")
}

func createToken(userid uint64) (*TokenDetails, error) {
	td := &TokenDetails{}

	var err error

	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = userid
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func verifyToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := extractToken(c)

	if tokenString == "" {
		return nil, fmt.Errorf("no access_token found in cookie")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func extractToken(c *gin.Context) string {
	bearToken, err := c.Cookie("access_token")

	if err != nil {
		return ""
	}

	return bearToken
}

func tokenValid(c *gin.Context) error {
	token, err := verifyToken(c)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

func ExtractTokenMetadata(c *gin.Context) (*AccessDetails, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}

		return &AccessDetails{
			UserId: userId,
		}, nil
	}
	return nil, err
}

func validEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func checkEmailExists(email string) (bool, error) {
	row, err := db.Query(
		context.Background(),
		"SELECT email FROM users WHERE email=$1 LIMIT 1",
		email,
	)
	if err != nil {
		log.Printf("checkEmailExists query: %v\n", err.Error())
		return false, err
	}

	var queryEmail string

	for row.Next() {
		err = row.Scan(&queryEmail)
		if err != nil {
			log.Printf("Row scan: %v\n", err.Error())
			return false, err
		}
	}

	if queryEmail != "" {
		return true, nil
	}

	return false, nil
}

func validUsername(username string) bool {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]{3,64}$`)
	return usernameRegex.MatchString(username)
}

func checkUsernameExists(username string) (bool, error) {
	row, err := db.Query(
		context.Background(),
		"SELECT username FROM users WHERE username=$1 LIMIT 1",
		username,
	)
	if err != nil {
		log.Printf("checkEmailExists query: %v\n", err.Error())
		return false, err
	}

	var queryUsername string

	for row.Next() {
		err = row.Scan(&queryUsername)
		if err != nil {
			log.Printf("Row scan: %v\n", err.Error())
			return false, err
		}
	}

	if queryUsername != "" {
		return true, nil
	}

	return false, nil
}

func validPassword(password string) bool {
	passwordRegex := regexp.MustCompile(`^[a-zA-Z0-9\$\!\@\_\%\^\*\&\(\)]{8,64}$`)
	return passwordRegex.MatchString(password)
}

func createUser(user User) error {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO users (username, password, email) VALUES ($1, $2, $3)",
		user.Username, HashAndSalt([]byte(user.Password)), user.Email,
	)
	if err != nil {
		log.Printf("createUser: exec query: %s\n", err.Error())
		return err
	}

	return nil
}
