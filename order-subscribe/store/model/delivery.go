package model

type Delivery struct {
	OrderID string `json:""`
	Name    string `json:""`
	Phone   string `json:""`
	Zip     string `json:""`
	City    string `json:""`
	Address string `json:""`
	Region  string `json:""`
	Email   string `json:""`
}
