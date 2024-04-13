package entity

type Role int32

type UserClaims struct {
	Id   string `json:"id"`
	Role Role   `json:"role"`
}
