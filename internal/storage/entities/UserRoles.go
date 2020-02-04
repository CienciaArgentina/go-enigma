package entities

type UserRoles struct {
	UserId int `db:"user_id"`
	RoleId int `db:"role_id"`
}
