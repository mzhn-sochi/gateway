package ticketservice

import "errors"

var (
	ErrTicketNotFound = errors.New("ticket not found")
)
