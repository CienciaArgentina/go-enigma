package entities

type Roles struct {
	RoleId         int    `db:"role_id"`
	Name           string `db:"name"`
	NormalizedName string `db:"normalized_name"`
}
