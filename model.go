package auth

type User struct {
	ID       uint64 `json:"id"`
	Email    string `form:"email" json:"email"`
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type TokenDetails struct {
	AccessToken string
}

type AccessDetails struct {
	UserId uint64
}
