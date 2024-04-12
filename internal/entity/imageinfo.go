package entity

type ImageInfo struct {
	Product     string            `json:"product"`
	Description string            `json:"description"`
	Price       float32           `json:"price"`
	Measure     *Measure          `json:"measure"`
	Attributes  map[string]string `json:"attributes"`
}

type Measure struct {
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"`
}
