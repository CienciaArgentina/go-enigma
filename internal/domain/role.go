package domain

type AssignedRole struct {
	AuthID string `json:"auth_id"`
	Roles  []Role `json:"roles"`
}

// Role Structure of a role to be assumed by a user.
type Role struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Claims      []Claim `json:"claims"`
}

// Claim Defines role permissions.
type Claim struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type AssignRoleRequest struct {
	AuthID int64 `json:"auth_id"`
	RoleID int   `json:"role_id"`
}
