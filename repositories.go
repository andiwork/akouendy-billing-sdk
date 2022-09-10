package billing

import (
	"gorm.io/gorm"
)

type BillingRepository struct{}

func (r BillingRepository) CreateBillingOrder(m *BillingTransaction, db *gorm.DB) (err error) {
	err = db.Debug().Create(m).Error
	return
}

func (r BillingRepository) GetBillingOrderByPaymentToken(paymentToken string, db *gorm.DB) (model BillingTransaction, err error) {
	err = db.First(&model, "order_payment_token = ?", paymentToken).Error
	return
}

func (r BillingRepository) GetBillingOrderByTransactionId(transactionID string, db *gorm.DB) (model BillingTransaction, err error) {
	err = db.First(&model, "app_trx_id = ?", transactionID).Error
	return
}
func (r BillingRepository) GetBillingOrderByPaymentTokenAndTable(table string, paymentToken string, db *gorm.DB) (model BillingTransaction, err error) {
	err = db.Joins("INNER JOIN ? t ON t.token = billing_transactions.app_trx_id").First(&model, "order_payment_token = ?", paymentToken).Error
	return
}
