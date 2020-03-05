package recovery

type PasswordResetDto struct {
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	Token           string `json:"token"`
}
