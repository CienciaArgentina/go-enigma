package domain

// UserProfile defines model for UserProfile.
type UserProfile struct {
	AuthID                 int64                      `json:"auth_id" db:"id"`
	UserName               string                   `json:"user_name" db:"username"`
	Email                  string                   `json:"email" db:"email"`
}
