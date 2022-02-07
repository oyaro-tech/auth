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
		c.JSON(http.StatusBadRequest, "Please send a request body")
		log.Println("Empty request body")
		return
	case err != nil:
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		log.Println(err.Error())
		return
	}

	user := getPasswordHashByLogin(u.Username)

	//compare the user from the request, with the one we defined:
	if !comparePasswords(u.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		log.Println("Wrong credentials")
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
	_, err := ExtractTokenMetadata(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	c.SetCookie("access_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, "Successfully logged out")
}

func Register(c *gin.Context) {
	var user User

	err := c.BindJSON(&user)
	if err != nil {
		log.Printf("Singup: %s\n", err.Error())
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	if !validUsername(user.Username) {
		log.Printf("Singup: invalid username: %s\n", user.Username)
		c.JSON(http.StatusUnprocessableEntity, "Invalid username!")
		return
	}

	status, err := checkUsernameExists(user.Username)
	if err != nil {
		log.Printf("Singup: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if status {
		c.JSON(
			http.StatusUnprocessableEntity,
			fmt.Sprintf("Username %s exists!", user.Username),
		)
		return
	}

	if !validEmail(user.Email) {
		log.Printf("Singup: invalid email: %s\n", user.Email)
		c.JSON(http.StatusUnprocessableEntity, "Invalid email!")
		return
	}

	status, err = checkEmailExists(user.Email)
	if err != nil {
		log.Printf("Singup: %s\n", err.Error())
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if status {
		c.JSON(
			http.StatusUnprocessableEntity,
			fmt.Sprintf("Email %s exists!", user.Email),
		)
		return
	}

	if !validPassword(user.Password) {
		log.Printf("Singup: invalid password: %s\n", user.Password)
		c.JSON(http.StatusUnprocessableEntity, "Invalid password!")
		return
	}

	err = createUser(user)
	if err != nil {
		log.Printf("Singup: createUser: %s\n", err.Error())
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
