package billing

import (
	utils "github.com/andiwork/gorm-utils"
)

type SubscriptionStatus string

const (
	REDIRECT  SubscriptionStatus = "REDIRECT"
	SUBSCRIBE SubscriptionStatus = "SUBSCRIBE"
)

type BillingTransaction struct {
	utils.Model
	OrderID           string
	PaymentToken string
	AppTrxId          string
	CountryAlpha3     string `gorm:"default:SEN"`
	PaymentID           string
}

type SubscriptionResponse struct {
	Data        interface{}
	RedirectUrl string
	Status      SubscriptionStatus
}

// Migrate return models for gorm
func MigrateModels() (models []interface{}) {
	models = append(models, new(BillingTransaction))
	return
}
