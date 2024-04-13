package dto

import "github.com/mzhn-sochi/gateway/internal/entity"

type Ticket struct {
	*entity.Ticket
	*entity.User
}
