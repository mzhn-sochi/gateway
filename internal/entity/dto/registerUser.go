package dto

import "github.com/mzhn-sochi/gateway/internal/entity"

type RegisterUser struct {
	entity.User
	Password string `json:"password"`
}
