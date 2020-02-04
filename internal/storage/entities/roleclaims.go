package entities

type RoleClaims struct {
	RoleClaimId int    `db:"role_claim_id"`
	RoleId      int    `db:"role_id"`
	ClaimType   string `db:"claim_type"`
	ClaimValue  string `db:"claim_value"`
}
