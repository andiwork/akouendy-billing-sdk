package billing

import "time"

const (
	SUCCESS PaymentStatus = "SUCCESS"
)

type OrderRequest struct {
	CustomerEmail    string
	CustomerFullName string
	CustomerID       string `json:"CustomerId"`
	BillingProvider  string
	CountryAlpha3    string
	PriceID          string `json:"PriceId"`
	AppID            string `json:"AppId"`
	Webhook          string
	ReturnUrl        string
	OwnedBy          string
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

type PaymentWebhook struct {
	Hash          string `json:"Hash"`
	Status        string `json:"Status"`
	TransactionID string `json:"TransactionID"`
}

type PaymentStatus string

type OrderSubsRequest struct {
	CustomerID string `json:"CustomerId"`
	PriceID    string `json:"PriceId"`
}

type OrderSubsReponse struct {
	CustomerID  string `json:"CustomerId"`
	PriceID     string `json:"PriceId"`
	IsSubscribe bool
	ExpireDate  time.Time
}
