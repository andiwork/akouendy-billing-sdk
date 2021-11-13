package billing

import (
	utils "github.com/andiwork/gorm-utils"
)

type BillingTransaction struct {
	utils.Model
	OrderID           string
	OrderPaymentToken string
	AppTrxId          string
}

// Migrate return models for gorm
func MigrateModels() (models []interface{}) {
	models = append(models, new(BillingTransaction))
	return
}
