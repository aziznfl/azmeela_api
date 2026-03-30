package domain

type ProductColorResponse struct {
	ID            int    `json:"id"`
	ProductCodeID int    `json:"product_code_id"`
	Color         string `json:"color"`
}

type ProductSizeTypeResponse struct {
	ID       int     `json:"id"`
	SizeName string  `json:"size_name"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
	Weight   int     `json:"weight"`
}
