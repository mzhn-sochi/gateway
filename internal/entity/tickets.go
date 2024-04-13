package entity

type Filters struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

type TicketFilters struct {
	Filters
	Status *string `json:"status"`
	UserId *string `json:"userId"`
	Phone  *string `json:"phone"`
}

type Ticket struct {
	Id          string  `json:"id"`
	UserId      string  `json:"userId"`
	Status      Role    `json:"status"`
	ImageUrl    string  `json:"imageUrl"`
	ShopAddress string  `json:"shopAddress"`
	CreatedAt   int64   `json:"createdAt"`
	UpdatedAt   *int64  `json:"updatedAt"`
	Reason      *string `json:"reason"`
}
