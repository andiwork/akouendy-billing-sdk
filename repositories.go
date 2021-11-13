package billing

import (
	"gorm.io/gorm"
)

type BillingRepository struct{}

func (r BillingRepository) CreateBillingOrder(m *BillingTransaction, db *gorm.DB) (err error) {
	err = db.Debug().Create(m).Error
	return
}
