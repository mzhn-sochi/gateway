package dto

type CreateTicket struct {
	UserId   string `json:"userId"`
	ShopName string `json:"shopName"`
	ShopAddr string `json:"shopAddr"`
	ImageUrl string `json:"imageUrl"`
}
