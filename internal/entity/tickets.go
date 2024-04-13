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
