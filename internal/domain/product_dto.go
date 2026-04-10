package domain

type ProductColorResponse struct {
	ID            int
	ProductCodeID int
	Color         string
}

type ProductSizeTypeResponse struct {
	ID               int
	ProductSizeID    int
	SizeName         string
	Price            float64
	Stock            int
	Weight           int
	CustomerTypeID   int
	CustomerTypeName string
}

