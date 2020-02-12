package register

type UserSignUpDto struct {
	Username string `json:username`
	Password string `json:password`
	Email    string `json:email`
}
