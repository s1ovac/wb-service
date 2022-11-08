package model

type Item struct {
	ID          int     `json:""`
	OrderID     string  `json:""`
	ChrtID      int     `json:""`
	TrackNumber string  `json:""`
	Price       float64 `json:""`
	Rid         string  `json:""`
	Name        string  `json:""`
	Sale        int     `json:""`
	Size        string  `json:""`
	TotalPrice  float64 `json:""`
	NmID        int     `json:""`
	Brand       string  `json:""`
	Status      int     `json:""`
}
