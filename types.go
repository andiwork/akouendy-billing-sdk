package billing

import (
	"time"
)

const (
	SUCCESS PaymentStatus = "SUCCESS"
	SANDBOX Environment   = "sandbox"
	PROD    Environment   = "prod"
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

type Environment string

type PaymentRequest struct {
	AppId         string
	TransactionId string
	TotalAmount   int
	Hash          string
	Description   string
	ReturnUrl     string
	Webhook       string
	Env           Environment
	Email         string
	FullName      string
	SuccessMsg    string
	FailedMsg     string
	ClientId      string
}

type PaymentResponse struct {
	Code       string
	Text       string
	Status     PaymentStatus
	Token      string
	PaymentUrl string
}

type PaymentStatusResponse struct {
	Status           PaymentStatus
	TotalPayedAmount int
	Email            string
}
