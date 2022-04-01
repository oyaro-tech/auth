package auth

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var u User
	err := c.ShouldBindJSON(&u)
	switch {
	case err == io.EOF:
		c.JSON(http.StatusBadRequest, "please send a request body")
		log.Println("empty request body")
		return
	case err != nil:
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		log.Println(err.Error())
		return
	}

	user := getPasswordHashByLogin(u.Username)

	//compare the user from the request, with the one we defined:
	if !comparePasswords(u.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, "please provide valid login details")
		log.Println("wrong credentials")
		return
	}

	ts, err := createToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		log.Println(err.Error())
		return
	}

	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		log.Println(err.Error())
		return
	}

	c.SetCookie("access_token", ts.AccessToken, 3600*12, "/", "", true, true)
}

func Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusUnauthorized, "unauthorized")
}

func Register(c *gin.Context) {
	var user User

	err := c.BindJSON(&user)
	if err != nil {
		log.Printf("register: %s\n", err.Error())
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if !validUsername(user.Username) {
		log.Printf("register: invalid username: %s\n", user.Username)
		c.JSON(http.StatusUnprocessableEntity, "invalid username")
		return
	}

	status, err := checkUsernameExists(user.Username)
	if err != nil {
		log.Printf("register: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if status {
		c.JSON(
			http.StatusUnprocessableEntity,
			fmt.Sprintf("username %s exists!", user.Username),
		)
		return
	}

	if !validEmail(user.Email) {
		log.Printf("register: invalid email: %s\n", user.Email)
		c.JSON(http.StatusUnprocessableEntity, "invalid email")
		return
	}

	status, err = checkEmailExists(user.Email)
	if err != nil {
		log.Printf("register: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if status {
		c.JSON(
			http.StatusUnprocessableEntity,
			fmt.Sprintf("email %s exists!", user.Email),
		)
		return
	}

	if !validPassword(user.Password) {
		log.Printf("register: invalid password: %s\n", user.Password)
		c.JSON(http.StatusUnprocessableEntity, "invalid password!")
		return
	}

	err = createUser(user)
	if err != nil {
		log.Printf("register: createUser: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

func TokenAuthMiddleware(c *gin.Context) {
	err := tokenValid(c)
	if err != nil {
		c.SetCookie("access_token", "", -1, "/", "", true, true)
		c.JSON(http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}
	c.Next()
}
