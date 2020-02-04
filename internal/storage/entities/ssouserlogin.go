package entities

type SSOUserLogin struct {
	SSOUserLoginId      int    `db:"sso_user_login_id"`
	UserId              int    `db:"user_id"`
	ProviderKey         string `db:"provider_key"`
	ProviderDisplayName string `db:"provider_display_name"`
}
