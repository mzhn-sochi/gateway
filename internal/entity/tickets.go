package entity

type Filters struct {
	//FromDate string `json:"from_date"`
	//ToDate   string `json:"to_date"`
	//FromTime string `json:"from_time"`
	//ToTime   string `json:"to_time"`
	//From     string `json:"from"`
	//To       string `json:"to"`
	//Sort     string `json:"sort"`
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

type TicketFilters struct {
	Filters
	Status *string `json:"status"`
	UserId *string `json:"user_id"`
}

type Ticket struct {
	Id          string `json:"id"`
	UserId      string `json:"userId"`
	Status      Role   `json:"status"`
	ImageUrl    string `json:"imageUrl"`
	ShopAddress string `json:"shopAddress"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   *int64 `json:"updatedAt"`
}
