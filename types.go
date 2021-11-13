package billing

type OrderRequest struct {
	CustomerEmail    string
	CustomerFullName string
	CustomerID       string `json:"CustomerId"`
	BillingProvider  string
	PriceID          string `json:"PriceId"`
	AppID            string `json:"AppId"`
	Webhook          string
}

type OrderResponse struct {
	OrderID      string `json:"OrderId"`
	PaymentUrl   string
	PriceID      string `json:"PriceId"`
	AppID        string `json:"AppId"`
	PaymentToken string `json:"PaymentToken"`
	Description  string
	Code         string
}
