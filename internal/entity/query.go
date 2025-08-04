package entity

type GetAdsQuery struct {
	Page     int
	Limit    int
	SortBy   string // "date" or "price"
	SortDir  string // "desc" or "asc"
	MinPrice float64
	MaxPrice float64
}
