package entity

type Role int32

type User struct {
	Phone      string `json:"phone"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
}

type UserCredentials struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UserClaims struct {
	Id   string `json:"id"`
	Role Role   `json:"role"`
}
